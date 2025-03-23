provider to server
```jsonc
v0.0|{
    "type": "handshake/request",
    "project": "proj-x",
    "handshake_id": UUID
}
```

server to provider
```jsonc
v0.0|{
    "type": "handshake/ack",
    "project": "proj-x",
    "handshake_id": UUID,
    "provider_id": UUID
}
```

(if the server refuses the provider)
server to provider
```jsonc
v0.0|{
    "type": "handshake/nack",
    "project": "proj-x",
    "handshake_id": UUID,
    "reason": "reason_code"
}
```


provider to server
```jsonc
v0.0|{
    "type": "announcement/handler",
    "project": "proj-x",
    "tool": "tool-y",
    "handshake_id": UUID,
    "provider_id": UUID,
    "tool_properties": {
        "name": "Tool X",
        "description": "x's description"
    }
}
```

server to provider
```jsonc
v0.0|{
    "class": "announcement/ack",
    "project": "proj-x",
    "tool": "tool-y",
    "handshake_id": UUID,
    "provider_id": UUID,
    "handler_id": UUID
}
```

(if the server does not accept the tool provider registration)
server to provider
```jsonc
v0.0|{
    "class": "announcement/nack",
    "project": "proj-x",
    "tool": "tool-y",
    "handshake_id": UUID,
    "provider_id": UUID,
    "reason": "nack_reason_code"
}
```