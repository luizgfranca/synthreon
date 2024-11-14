# PlatformLab

This project aims to be a easy way to provide a platform to create internal tools to create quick automations and provide ways to simplify the management of application's resources.


## Status 
This project is still in early prototyping stage, please come back later for a full proof-of-concept.


## Run the demonstrations
There are 2 early-stage technology demonstrations you can already run:

### Dependencies
To run the control panel and look at the examples you need to first prepare the dependencies
You will need:
 - go >= 1.22
 - nodejs 22
 - yarn (recomended)

### Preparing

Install ts-node globally
```bash
npm install -g ts-node
```

Go to the control panel's web directory
```bash
cd control-panel/web
```

Install javascript's dependencies
```bash
yarn install
```

Go back to the root directory and enter the node's SDK directory
```bash
cd ../..
cd sdk/node
```

Install javascript's dependencies
```bash
yarn install
```

### Building and running

First start the backend server:
```bash
cd control-panel
go run .
```

Open another terminal window and start the frontend development server
```bash
cd control-panel/web
yarn dev
```

Open another terminal window, choose the example you want to run and run it with `ts-node`
For instance, here's how you would run the `ping-pong` example:
```bash
cd sdk/node
ts-node example/ping-pong.ts 
```

Now to access and test the tool enter the browser on `http://localhost:5173/exp`


## Components
 - **control-panel**: will host the tools
 - **control-panel/web**: frontend of the tool server
 - **sdk**: will enable applications to connect to the server and provide the tools
 - **agent**: prototype for internal agent for infrasctructure management (development paused)