provider to server
```jsonc
v0.0|{
    "type": "handshake/request",
    "project": "proj-x",
    "announcement_id": UUID
}
```

server to provider
```jsonc
v0.0|{
    "type": "handshake/accept",
    "project": "proj-x",
    "announcement_id": UUID,
    "provider_id": UUID
}
```

(if the server refuses the provider)
server to provider
```jsonc
v0.0|{
    "type": "handshake/noaccept",
    "project": "proj-x",
    "announcement_id": UUID,
    "reason": "reason_code"
}
```


provider to server
```jsonc
v0.0|{
    "type": "announcement/handler",
    "project": "proj-x",
    "tool": "tool-y",
    "announcement_id": UUID,
    "provider_id": UUID
}
```

server to provider
```jsonc
v0.0|{
    "class": "announcement/ack",
    "project": "proj-x",
    "tool": "tool-y",
    "announcement_id": UUID,
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
    "announcement_id": UUID,
    "provider_id": UUID,
    "reason": "nack_reason_code"
}
```