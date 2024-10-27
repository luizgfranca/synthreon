import { PlatformConnection } from "../src/platform-connection";

const connection = new PlatformConnection({
    endpoint: 'ws://localhost:8080/api/tool/provider/ws'
})

connection.listen();