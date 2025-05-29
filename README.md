# Synthreon

Synthreon is a project that aims to provide a platform for an application developer, a team, or an entire company to create internal tools for quick automation and simplify the management of an application's resources or internal processes.

Here is a quick demonstration of what it can do in its current state of development:





https://github.com/user-attachments/assets/68625bdf-5db4-40b6-b137-dd374f4fb8de






## Status 
This project is still in the early prototyping stage and is currently just a proof-of-concept of the idea. You are invited to test and your feedback is highly valuable.


## How to I use it

You can follow the [Getting Started](/docs/getting-started.md) to setup a server and build your first tool using it.

### Running the server

You can use the Docker image, or build the server from source. Here is how you do it.

#### Docker image

To run the docker image execute this on your terminal:

```bash
docker run --name st \
  -e ACCESS_TOKEN_SECRET_KEY={{your_secret_key}} \
  -e ROOT_EMAIL={{your_login_email}} \
  -e ROOT_PASSWORD={{your_login_password}} \
  -p "25256:25256" \
  -v "$PWD"/data:/data \
synthreon/synthreon
```

Replacing `{{your_secret_key}}`, `{{your_login_email}}`, `{{your_login_password}}` with your own information.
(more specific details about all the supported parameters of the docker image can be found at: https://hub.docker.com/r/synthreon/synthreon)

#### Building from source

Requirements
 - go >= 1.23
 - node.js >= 22
 - yarn

To execute the project from source, enter the `control-panel` directory and execute the `dev.sh` script.

```
cd control-panel
./dev.sh
```

This is going to create a database called `test.db` in the current working directory, with a default user `test@test.com` and a default password `password`. If you want to edit this details you can edit their variables directly in the `./dev.sh` script.

## Run the demonstrations
There are some early-stage technology demonstrations you can already run, to do it you need to run the server with one of the methods instructed above, and then run the example.

To run the examples you need to first set up the dependencies. You will need:
 - node.js 22

After installing them enter http://localhost:25256 using your browser and log in with the test user credentials you set up during the first execution.

To start any of the example tool providers open another terminal window, choose the example you want to run inside the `packages/node-sdk/example` folder, and run it with `ts-node`

For instance, here's how you would run the `ping` example:
```bash
cd packages/node-sdk
ts-node example/ping.ts 
```

Select the "Sandbox" project.

Select the "sandbox" tool in the sidebar, and a screen containing the tool being run will appear ready for your interaction.

## Components
 - **control-panel**: will host the tools
 - **control-panel/web**: frontend of the tool server
 - **packages/node-sdk**: sdk for tool development
 - **packages/js-core**: generic javascript code used around the project
