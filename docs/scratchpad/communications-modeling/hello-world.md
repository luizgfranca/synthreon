

### Hello world case
client to server
```jsonc
v0.0|{
    "type": "interaction/open",
    "project": "proj-x",
    "tool": "tool-y",
    "session_id": "shdfkajhdflaksjdfhalksdjh",
}
```

server to provider
```jsonc
v0.0|{
    "type": "interaction/open",
    "project": "tool-y",
    "tool": "tool-y",
    "provider_id": UUID,
    "handler_id": UUID,
    "context_id": UUID
}
```

provider to server
```jsonc
v0.0|{
    "type": "command/finish",
    "project": "tool-y",
    "tool": "tool-y",
    "provider_id": UUID,
    "handler_id": UUID,
    "context_id": UUID,
    "execution_id": UUID,
    "result": {
        "status": "success" | "error" | "undetermined",
        "message": "Hello world"
    }
}
```

server to client
```jsonc
v0.0|{
    "type": "command/display",
    "project": "tool-y",
    "tool": "tool-y",
    "context_id": UUID,
    "display": {
        "type": "information",
        "information": {
            "type": "success" | "error" | "undetermined"
            "message": "Hello World"
        }
    }
}
```
