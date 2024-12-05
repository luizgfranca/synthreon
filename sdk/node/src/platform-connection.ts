import WebSocket, { RawData } from "ws";
import { PromptType, ToolEvent, DisplayDefinition, EventInput } from "./tool-event";
import { EventEmitter } from "node:events";

type UserCredentials = {
  username: string,
  password: string
}

type PlatformToolConnectionOptions = {
  endpoint: string;
  toolFunction: (elements: ToolComponents) => Promise<string>;
  credentials: UserCredentials
};

type ToolExecution = {
  project: string;
  tool: string;
  client: string;
  messageBus: EventEmitter;
  components: ToolComponents;
  promise: Promise<string>;
};

type ConnectionStatus =
  | "connected"
  | "waiting_ack"
  | "connecting"
  | "disconnected";

export type ToolComponents = {
  io: {
    prompt: (description: string, type: PromptType) => Promise<string>;
    textBox: (content: string) => Promise<void>
  };
};

type ClientDefinition = {
  project: string;
  tool: string;
  client?: string;
};

export class PlatformConnection {
  #endpoint: string;
  #websocket?: WebSocket;
  #toolFunction: (components: ToolComponents) => Promise<string>;
  #executions: ToolExecution[];
  #status: ConnectionStatus;
  #credentials: UserCredentials

  constructor(options: PlatformToolConnectionOptions) {
    this.#executions = [];
    this.#toolFunction = options.toolFunction;

    this.#endpoint = options.endpoint;
    this.#status = "disconnected";
    this.#credentials = options.credentials
  }

  listen() {
    try {
      this.#websocket = new WebSocket(this.#endpoint, {
        headers: {
          Authorization: `Basic ${Buffer.from(
            `${this.#credentials.username}:${this.#credentials.password}`
          ).toString("base64")}`,
        },
      });
      this.#status = "connecting";
    } catch (e) {
      throw new Error(`unable to connect to websocket: ${e}`);
    }

    this.#websocket.on("open", () => {
      console.debug("successfully connected to platform");
      this.#makeAnnouncement();
    });

    this.#websocket.on("error", (error) => {
      console.debug(`websocket connection error: ${error.message}`);
    });

    this.#websocket.on("message", (data: RawData) => {
      const event = JSON.parse(data.toString()) as ToolEvent;
      this.#onEventReceived(event);
    });
  }

  #sendEvent(event: ToolEvent) {
    console.debug("sending event: ", event);

    if (!this.#websocket) {
      throw new Error("Internal error: expected websocket to be defined.");
    }
    this.#websocket.send(JSON.stringify(event));
  }

  #onToolOpen(project: string, tool: string, client: string) {
    const messageBus = new EventEmitter();
    const components = this.#getToolComponents(messageBus);
    this.#setupOutboundHandlers(messageBus, { project, tool, client });

    const execution: ToolExecution = {
      promise: new Promise((resolve, reject) => {
        try {
          this.#toolFunction(components)
            .then((resultMessage) => resolve(resultMessage))
            .catch((e) => reject(e));
        } catch (e) {
          reject(e);
        }
      }),
      messageBus,
      components,
      project,
      tool,
      client,
    };

    this.#setupExecutionResultHandlers(execution);

    this.#executions.push(execution);
  }

  #onInput(origin: ClientDefinition, input: EventInput) {
    console.log("handling new input");

    // TODO: this is a simplified and VERY BAD initial approach
    //       every provider execution should have an id specifically
    //       addressable in the response to avoid uncertainty
    const destination = this.#executions.find(
      (execution) =>
        origin.project === execution.project &&
        origin.tool === execution.tool &&
        origin.client === execution.client
    );

    if (!destination) {
      throw new Error("no destination found to handle input");
    }

    if (!input) {
      throw new Error("input event received without input data");
    }

    destination.messageBus.emit("input", input);
  }

  #onToolInteraction(event: ToolEvent) {
    const { project, tool, client} = event
    
    switch (event.type) {
      case "open":
        return this.#onToolOpen(event.project, event.tool, event.client ?? "");
      case "input":
        if (!event.input) {
          throw new Error('expected to receive input parameter on "input" event')
        }

        return this.#onInput({ project, tool, client}, event.input)
    }
  }

  #makeAnnouncement() {
    if (!this.#websocket) {
      throw new Error(
        "Internal error: expected websocket to be defined to make announcement"
      );
    }

    console.debug("making provider annnouncement");
    this.#sendEvent({
      class: "announcement",
      type: "provider",
      project: "proj-x",
      tool: "tool-y",
      provider: 1,
    });

    this.#status = "waiting_ack";
    console.debug("announcement sent");
  }

  #onEventReceived(event: ToolEvent) {
    console.log("event received: ", event);

    // TODO: need to add sanity validations here
    switch (this.#status) {
      case "connected":
        if (event.class == "interaction") {
          this.#onToolInteraction(event);
        }
        break;
      case "waiting_ack":
        if (event.class == "announcement" && event.type == "ack") {
          this.#status = "connected";
          console.debug("announcement acknowleged");
        }
        break;
    }
  }

  #getPromptFunction(
    messageBus: EventEmitter
  ): (title: string, type: PromptType) => Promise<string> {
    return (title: string, type: PromptType) => {
      console.debug("getPromptFunction", { title, type });
      return new Promise((resolve, reject) => {
        console.debug("executing prompt dispatch", { title, type });
        const display: DisplayDefinition = {
          type: "prompt",
          prompt: {
            title,
            type,
          },
        };

        console.debug("sending display definition to messageBus", display);
        messageBus.emit("display", display);

        messageBus.on("display", (definition) => {
          console.debug("[probe] definition received", definition);
        });

        messageBus.on("input", (input: EventInput) => {
          console.log("input received by prompt handler: ", input);

          if (input.fields.length == 0) {
            return reject("Internal Error: did not receive field result");
          }

          resolve(input.fields[0].value);
        });
      });
    };
  }

  #getTextBoxFunction(
    messageBus: EventEmitter
  ): (content: string) => Promise<void> {
    return (content: string) => {
      return new Promise((resolve, reject) => {
        console.debug("executing textBox dispatch", { content });
        const display: DisplayDefinition = {
          type: "textbox",
          textBox: {
            content
          }
        };

        console.debug("sending display definition to messageBus", display);
        messageBus.emit("display", display);

        messageBus.on("display", (definition) => {
          console.debug("[probe] definition received", definition);
        });

        messageBus.on("input", (input: EventInput) => {
          console.log("input received by prompt handler: ", input);
          resolve();
        });
      })
    }
  }

  #getToolComponents(messageBus: EventEmitter): ToolComponents {
    return {
      io: {
        prompt: this.#getPromptFunction(messageBus),
        textBox: this.#getTextBoxFunction(messageBus)
      },
    };
  }

  #setupExecutionResultHandlers(execution: ToolExecution) {
    console.debug("setupExecutionBaseHandlers");
    execution.promise
      .then((resultMessage) => {
        console.log("executing success result trap");

        this.#sendEvent({
          class: "operation",
          type: "display",
          project: execution.project,
          tool: execution.tool,
          client: execution.client,
          display: {
            type: "result",
            result: {
              success: true,
              message: resultMessage,
            },
          },
        });
      })
      .catch((errorMessage) => {
        console.log("executing error result trap");
        this.#sendEvent({
          class: "operation",
          type: "display",
          project: execution.project,
          tool: execution.tool,
          client: execution.client,
          display: {
            type: "result",
            result: {
              success: false,
              message: errorMessage,
            },
          },
        });
      });
  }

  #setupOutboundHandlers(messageBus: EventEmitter, client: ClientDefinition) {
    console.debug("setupExecutionOutboundHandlers");

    messageBus.on("display", (definition: DisplayDefinition) => {
      console.debug("display definition received on messageBus", definition);

      const event: ToolEvent = {
        class: "operation",
        type: "display",
        project: client.project,
        tool: client.tool,
        client: client.client,
        display: definition
      };

      if (!this.#websocket) {
        throw new Error(
          "Internal error: expected websocket to be defined to send event"
        );
      }

      this.#sendEvent(event);
    });
  }
}
