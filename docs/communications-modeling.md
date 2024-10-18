### USER ENTERS TOOL EVENT

client to server
```json
{
    "class": "interaction",
    "type": "open",
    "project": "proj-x",
    "tool": "tool-y",
}
```

server to client
```json
{
    "class": "operation",
    "type": "display",
    "project": "proj-x",
    "tool": "tool-y",
    "display": {
        "type": "prompt",
        "elements":[
            {
                "type": "input",
                "label": "This is a test",
                "name": "generated"
            }
        ]
    }
}
```

client to server
```json
{
    "class": "interaction",
    "type": "input",
    "project": "proj-x",
    "tool": "tool-y",
    "input": {
        "fields": [
            {
                "name": "generated",
                "value": "user input"
            }
        ]
    }
}
```

server to client
```json
    {
        "class": "operation",
        "type": "display",
        "project": "proj-x",
        "tool": "tool-y",
        "display": {
            "type": "result",
            "result": {
                "success": true,
                "message": "Hello user input",
            }
        }
    }
```



### ALTERATIVE RESPONSE OPTION

server to client
```json
{
    "class": "operation",
    "type": "display",
    "project": "proj-x",
    "tool": "tool-y",
    "display": {
        "type": "view",
        "elements":[
            {
                "type": "text",
                "message": "hello world"
            }
        ]
    }
}
```


### PROVIDER

self announcement
provider to server
```json
{
    "class": "announcement",
    "type": "provider",
    "project": "proj-x",
    "tool": "tool-y",
    "provider": 1
}
```

server to provider
```json
{
    "class": "announcement",
    "type": "ack",
    "project": "proj-x",
    "tool": "tool-y",
    "provider": 1
```

in the future i should have an announcementId in the announcement to be acknowleged in the ack

provider to server
```json
    {
        "class": "operation",
        "type": "display",
        "project": "proj-x",
        "tool": "tool-y",
        "display": {
            "type": "result",
            "result": {
                "success": true,
                "message": "Hello user input",
            }
        },
        "client": 1
    }
```

**Every time a client sends a message to a provider, the backend should resolve the provider for that tool, and also complement the event with the client's id for the provider to be able to know to which client it should respond to**
(there's a security problem with this implementation, in which a rogue provider would send a message to another client, i should think in a solution for this)