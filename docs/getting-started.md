## Getting started

This is a step-by-step guide that will help you create your first simple Synthreon tool. Here we will explore the basic concepts in order to understand how to work with Synthreon, how to setup your development environment, and finally, how to build the code for your interaction. All of this, using as a practical example a simple task time tracker tool, a simple timer to monitor how much time you are taking to do a set of tasks so you can understand how long you've taken on each one. 

All the code used here can be accessed in the example repository on github: [link](https://github.com/luizgfranca/synthreon-getting-started-example)


#### Setting up the server

First, to run the tools you need a server instance. A server hosts the many automations you will develop, these are called providers, and they will announce themselves to the server, which will allow the user to use any of the ones it knows about. These "providers" may be standalone tools like our example, but they may also be just an application that exposes itself using the SDK.   

The best way to start the server is using the Docker image. To download it and start a new container you can just run the following command line: 

```bash
docker run --name st \
  -e ACCESS_TOKEN_SECRET_KEY={{your_secret_key}} \
  -e ROOT_EMAIL={{your_login_email}} \
  -e ROOT_PASSWORD={{your_login_password}} \
  -p "25256:25256" \
  -v "$PWD"/data:/data \
synthreon/synthreon
```
Remember to replace {{your_secret_key}}, {{your_login_email}}, {{your_login_password}} with your own information.
(more specific details about all the supported parameters of the docker image can be found at: https://hub.docker.com/r/synthreon/synthreon)

After running it, if you access http:localhost:25256 this screen should appear:

![login screen](/docs/img/getting-started-0.png)

You can now log in with the email and password you used on the docker command, and a list of the current projects should appear:

![projects screen](/docs/img/getting-started-1.png)

Currently only the default ones will be there, "Test" and "Sandbox". If you enter your Sandbox you will be shown this screen:

![Sandbox project tools](/docs/img/getting-started-2.png)

This is the screen for interacting with the tools you created. Select the `Sandbox` tool in the `Sandbox` example and the message `Waiting provider...` should appear, this means it is waiting for our tool implementation to start. So lets implement it then.


#### Creating the tool project.

For the tool development you can currently choose between Node.js or Bun. Although Bun is highly recommended. 

The following instructions are for Bun, you can install Bun following the instructions in [https://bun.sh](https://bun.sh). But doing the equivalent using Node.js with Yarn or NPM should also be possible.

First create your new project and enter its directory:
```
mkdir tasktimer
cd tasktimer
bun init -y
```

Install the SDK:
```
bun add @synthreon/sdk
```

Bun should have created the main file `index.ts` for you. You can also use pure Javascript creating a `index.js`, but for this tutorial we will use typescript.

Add this basic hello world as the content, replace the username and password with the ones you set up on the server:

```ts
import { ToolProvider } from '@synthreon/sdk'

const provider = new ToolProvider({
    project: 'sandbox',
    endpoint: 'ws://localhost:25256/ws/tool/provider',
    credentials: {
        username: '{{your_login_email}}',
        password: '{{your_login_password}}',
    },
})


provider.tool('sandbox', async ({ io }) => {
    return 'hello world' 
})

provider.listen()
```

Now add the following script configuration to your package.json:
```
  "scripts": {
    "start": "bun index.ts"
  },
```

And run:
```
bun run start
```

Now, openning your browser tab, the following screen should appear, this means your tool implementation has successfully connected to the server.

![hello world display](/docs/img/getting-started-3.png)

What we did here was the initialization of your provider creating a new instance of `Provider`, setting it up with the details of the server it should connect to, your credentials for authentication, and the project to whitch it should attach its tools.

After this we should declare all of our tools and their implementation using the `tool()` method from our recently created `Provider`. The declaration method receives a unique ID for the tool in the current project, and the tool implementation, which is provided by a function.

Finally, after declaring all tools we can call `listen()` to connect to the server and announce tools. Once connected `provider.listen()` will block until it loses the connection to the server or the user sends a termination signal to the application.


### Tracker implementation

Now lets look into implementing our application logic. The first thing we need is somewhere to save the tasks we are performing at the moment, for this we must first create a global map to keep their information.

```ts
let taskMap: Record<string, Date> = {};
```

This dictionary saves the name of the task as a string identifier and a date, that should be the time in which the task was created. Using this approach, we can always calculate the elapsed time of the tracked task by subtracting the current time from this one.

The first thing we as users will need to do is to be able to create new tasks, so then lets declare a new tool for this functionality, to receive the name of the task that is being  started and save it to our dictionary.

```ts
provider.tool('create', async ({ io }) => {
    const name = await io.prompt({
        title: 'Task name:',
        type: 'string'
    })

    if (taskMap[name]) {
        throw 'task already exists'
    }

    taskMap[name] = new Date();


    return 'Created.'
})
```

See that now we used a new kind of function. `io.prompt()` prompts the user for a value to be typed, using it you don't need to worry about building the UI yourself, you can just give a title explaining to the user the value needed and the field type. It returns a promise with the value typed by the user when they respond. In this case we prompt the user for the name of the task to be created, and subsequently use this value to create a new position on the task dictionary for it.

The next step is to create another one to show us the elapsed time of the currently tracked tasks.

First, let's create a simple utility function to print the elapsed type in a prettier format.

```ts
function timeStr(seconds: number): string {
    const days = Math.floor(seconds / (24 * 60 * 60));
    seconds %= 24 * 60 * 60;

    const hours = Math.floor(seconds / (60 * 60));
    seconds %= 60 * 60;

    const minutes = Math.floor(seconds / 60);
    seconds %= 60;

    const parts: string[] = [];
    if (days > 0) {
        parts.push(`${days}d`);
    }

    if (hours > 0) {
        parts.push(`${hours}h`);
    }

    if (minutes > 0) {
        parts.push(`${minutes}m`);
    }

    if (seconds > 0 || parts.length === 0) {
        parts.push(`${Math.floor(seconds)}s`);
    }

    return parts.join(", ");
}
```

And now our tool that displays a table with all the current tasks and their elapsed time.

```ts
provider.tool('show', async ({ io }) => {
    const data = Object.keys(taskMap).map(k => ({
        name: k,
        timeSpent: timeStr(((new Date()).getTime() - taskMap[k].getTime()) / 1000),
    }))

    await io.table({
        title: 'Currently active tasks',
        content: data,
        columns: {
            name: 'Name',
            timeSpent: 'Time spent',
        }
    })

    return 'ok';
})
````

For this one we use a `io.table` component, it renders a table on the screen with the data in the array passed in the `content` attribute. In this case we have 2 properties in the objects from this array, `name` and `timeSpent`, so for legibility, we also setup in the `columns` attribute the friendly names that should appear as the column names for this properties.

To see how it works stop the `bun run start` command execution, run it again, go to your browser, refresh the browser tab, create a new "task A", wait a moment and then create a new "task B". 

![new tools visible](/docs/img/getting-started-4.png)

![create tool](/docs/img/getting-started-5.png)

Entering the show tool it should look something loke this:

![showing table of tools](/docs/img/getting-started-6.png)


Now, to allow the user to remove the completed tasks from the list, create a new tool that will remove it from the map.

```ts
provider.tool('delete', async ({ io }) => {
    const name = await io.selection({
        description: 'Select the task to be deleted:',
        options: Object.keys(taskMap).map((name) => ({
            key: name,
            text: name
        }))
    );

    if (!taskMap[name]) {
        throw 'task not found'
    }

    delete taskMap[name]

    return 'Deleted.'
})
```

This one uses a `io.selection` component. This component shows a list of cards from which the user can only select one of the options. The options are given by an array with their unique key, and a text description for the user to choose between them. When the user makes the choice on the UI, the selected option's key is returned from the function. In this case it is used to obtain the name of the task the user chose to remove when they entered the tool.

![delete tool selection](/docs/img/getting-started-7.png)

With this simple steps we have a complete set of behaviors for a simple task timer.


### Where to go from here.

As was said in the beginning, this is a very simple example to demonstrate the basic principles of developing tools with Synthreon, but there are many ways it is useful to automate processes, and easily expose interfaces to control your applications, without taking on big projects. If it was your intention to make this automation more complete, your could, for instancce, use a generic ORM or database library here to integrate with an existing or a dedicated database, any concept you can use in JS development, you can also use here, without being tied down by a specific runtime, or low-code development process.

This is where Synthreon can help you the most, in the creation of simple tools to simplify your life, and to improve your team's productivity.
