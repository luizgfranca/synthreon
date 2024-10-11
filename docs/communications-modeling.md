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
