

### Ping Pong case
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
    "type": "command/display",
    "project": "project-x",
    "tool": "tool-y",
    "provider_id": UUID,
    "handler_id": UUID,
    "context_id": UUID,
    "execution_id": UUID,
    "display": {
        "type": "prompt",
        "prompt": {
            title: "Input ping to receive pong",
            type: "string"
        }
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
    "terminal_id": UUID,
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
```jsonc
v0.0|{
    "type": "interaction/input",
    "project": "tool-y",
    "tool": "tool-y",
    "execution_id": UUID,
    "context_id": UUID,
    "terminal_id": UUID,
    "input": {
        "fields": {
            "name": "p_name",
            "value": "ping"
        }
    }
}
```

server to provider
```jsonc
v0.0|{
    "type": "interaction/input",
    "project": "tool-y",
    "tool": "tool-y",
    "provider_id": UUID,
    "handler_id": UUID,
    "execution_id": UUID,
    "input": {
        "fields": {
            "name": "p_name",
            "value": "ping"
        }
    }
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
    "execution_id": UUID,
    "result": {
        "status": "success",
        "message": "pong"
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
        "type": "result",
        "information": {
            "status": "success",
            "message": "pong"
        }
    }
}
```
