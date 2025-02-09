# PlatformLab

PlatformLab (name subject to change) is a project that aims to provide a platform for an application developer, a team, or an entire company to create internal tools for quick automation and simplify the management of an application's resources.

Here is a quick demonstration of what it can do in its current state of development:


https://github.com/user-attachments/assets/a6ccb9be-b8a5-43b0-a4d2-37e543baf9ff


## Status 
This project is still in the early prototyping stage and is currently just a proof-of-concept of the idea, please return later for a MVP.


## Run the demonstrations
There are 3 early-stage technology demonstrations you can already run, to do it you need to run the server and one of the examples.

### Server (Docker)

Run the following command to build and run a docker container locally:

```bash
tools/run-docker.sh
```

This will instantiate a docker container running the Control Panel's server


### Examples

To run the examples you need first to set up the dependencies. You will need:
 - node.js 22
 - yarn (recommended)


After installing them enter http://localhost:5173 using your browser and log in with the test user credentials, they can be found in control-panel/demo.env, but the default ones are test@test.com/password

To start any of the example tool providers open another terminal window, choose the example you want to run inside the `packages/node-sdk/example` folder, and run it with `ts-node`

For instance, here's how you would run the `ping-pong` example:
```bash
cd packages/node-sdk
ts-node example/ping-pong.ts 
```

Select the "Sandbox" project.

Select the "sandbox" tool in the sidebar, and a screen containing the bein run will appear for the user to interact with it.

## Components
 - **control-panel**: will host the tools
 - **control-panel/web**: frontend of the tool server
 - **packages/node-sdk**: sdk for tool development
 - **packages/js-core**: generic javascript code used around the project
 - **agent**: prototype for the internal agent for infrastructure management (development paused)
