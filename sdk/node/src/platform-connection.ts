import WebSocket, { RawData } from "ws";
import { ToolEvent } from "./tool-event";
import { EventEmitter } from "node:events";

type PlatformToolConnectionOptions = {
  endpoint: string;
  toolFunction: () => string;
};

type ToolExecution = {
  project: string;
  tool: string;
  client: string;
  messageBus: EventEmitter;
  promise: Promise<string>;
};

type ConnectionStatus =
  | "connected"
  | "waiting_ack"
  | "connecting"
  | "disconnected";

export class PlatformConnection {
  #endpoint: string;
  #websocket?: WebSocket;
  #toolFunction: () => string;
  #executions: ToolExecution[];
  #status: ConnectionStatus;

  constructor(options: PlatformToolConnectionOptions) {
    this.#executions = [];
    this.#toolFunction = options.toolFunction;

    this.#endpoint = options.endpoint;
    this.#status = "disconnected";
  }

  listen() {
    try {
      this.#websocket = new WebSocket(this.#endpoint);
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
      console.debug("event received: ", event);
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
    const execution: ToolExecution = {
      messageBus: new EventEmitter(),
      promise: new Promise((resolve, reject) => {
        try {
          const resultMessage = this.#toolFunction();
          resolve(resultMessage);
        } catch (e) {
          reject(e);
        }
      }),
      project,
      tool,
      client,
    };

    execution.promise
      .then((resultMessage) => {
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
              message: errorMessage,
            },
          },
        });
      });

    this.#executions.push(execution);
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
    // TODO: need to add sanity validations here
    switch (this.#status) {
      case "connected":
        if (event.class == "interaction" && event.type == "open") {
          this.#onToolOpen(event.project, event.tool, event.client ?? "");
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
}
