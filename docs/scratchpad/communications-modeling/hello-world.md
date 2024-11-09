

### Hello world case
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
    "type": "command/finish",
    "tool": "tool-y"
    // still to specify how to transmit the destination client
    "result": {
        "status": "success" | "error" | "undetermined"
        "message": "Hello world"
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
            "status": "success" | "error" | "undetermined"
            "message": "Hello World"
        }
    }
}
```
