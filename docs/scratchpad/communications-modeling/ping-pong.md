

### Ping Pong case
client to server
```
v0.0|{
    "type": "interaction/open",
    "project": "proj-x",
    "tool": "tool-y"
}
```

server to provider
```
v0.0|{
    "type": "interaction/open",
    "tool": "tool-y"
    // still to specify how to transmit the origin client
}
```

provider to server
```
v0.0|{
    "type": "command/display",
    "tool": "tool-y"
    // still to specify how to transmit the destination client
    //
    // need to find a way to rereference the execution inside the provider
    // to where the responding input should go
    // (
    //    this should be something the server knows but the client doesn't
    //    so the client cannot try to send commands to another execution
    // )
    "display": {
        "type": "prompt",
        "prompt": {
            title: 'Input ping to receive pong',
            type: "string"
        }
    }
}
```

server to client
```
v0.0|{
    "type": "command/display",
    "tool": "tool-y"
    "execution_id": UUID,
    // it should return a reference to the provider that is serving this
    // still to specify how to transmit the destination client

    "display": {
        "type": "prompt",
        "prompt": {
            "name": "p_name",
            "title": "Input ping to receive pong",
            "type": "string"
        }
    }
}
```

client to server
```
v0.0|{
    "type": "interaction/input",
    "tool": "tool-y"
    "execution_id": UUID,
    // it should return a reference to the provider that is serving this
    // still to specify how to transmit the destination client

    "input": {
        "fields": {
            "name": "p_name",
            "value": "ping"
        }
    }
}
```

server to provider
```
v0.0|{
    "type": "interaction/input",
    "tool": "tool-y"
    // it should return a reference to the provider that is serving this
    // still to specify how to transmit the destination client

    "input": {
        "fields": {
            "name": "p_name",
            "value": "ping"
        }
    }
}
```

provider to server
```
v0.0|{
    "type": "command/finish",
    "tool": "tool-y"
    // still to specify how to transmit the destination client
    "result": {
        "status": "success"
        "message": "pong"
    }
}
```

server to client
```
v0.0|{
    "type": "command/display",
    "tool": "tool-y"
    "execution_id": UUID,
    // it should return a reference to the provider that is serving this
    // still to specify how to transmit the destination client

    "display": {
        "type": "result",
        "result": {
            "status": "success",
            "message": "pong"
        }
    }
}
```
