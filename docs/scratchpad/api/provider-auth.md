provider to server

```http
/api/tool/provider/ws
Authorization Basic `username:password`b64
```
server should only upgrade the websocket connection if the credentials are corerct
otherwhise HTTP 401