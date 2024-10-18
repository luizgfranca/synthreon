package connectionmgr

import (
	"encoding/json"
	"log"
	"platformlab/controlpanel/model"
	"sync"

	"github.com/gorilla/websocket"
)

type ToolClient struct {
	id                   uint
	connection           *websocket.Conn
	chSend               chan []byte
	done                 chan bool
	activeHandlerThreads sync.WaitGroup
	manager              *ClientMgr
}

type ClientMgr struct {
	clients             []ToolClient
	ProviderMgr         *ProviderMgr
	clientAdditionMutex sync.Mutex
}

func (m *ClientMgr) NewClient(connection *websocket.Conn) uint {
	m.clientAdditionMutex.Lock()
	idx := uint(len(m.clients))
	m.clients[idx] = *NewToolClient(m, idx, connection)
	m.clientAdditionMutex.Unlock()
	log.Print("Registered new client with ID: ", idx)

	return idx
}

func (m *ClientMgr) SendEvent(client uint, e *model.ToolEvent) error {
	m.clientAdditionMutex.Lock()
	if len(m.clients) <= int(client) {
		m.clientAdditionMutex.Unlock()
		return &model.GenericLogicError{Message: "client not found"}
	}
	m.clientAdditionMutex.Unlock()

	m.clients[client].SendEvent(e)
	return nil
}

func NewToolClient(manager *ClientMgr, id uint, connection *websocket.Conn) *ToolClient {
	client := ToolClient{
		id:                   id,
		connection:           connection,
		chSend:               make(chan []byte),
		done:                 make(chan bool),
		activeHandlerThreads: sync.WaitGroup{},
		manager:              manager,
	}

	client.activeHandlerThreads.Add(2)

	go client.messageReceiver()
	go client.messageSender()

	return &client
}

func (c *ToolClient) SendEvent(event *model.ToolEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Print("[toolprovider] error encoding response: ", err.Error())
	}

	c.chSend <- data
}

func (c *ToolClient) Close() {
	c.done <- false
	c.activeHandlerThreads.Wait()

	c.connection.Close()
}

func (c *ToolClient) messageReceiver() {
	var event model.ToolEvent

	for {
		select {
		case <-c.done:
			return
		default:
			// should continue
		}

		msgtype, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Print("[toolclient] websocket message receiving error: ", err.Error())
			break
		}

		if msgtype != websocket.TextMessage {
			log.Print("[toolclient] unexpeted message type, type sould be a textMessage")
			break
		}

		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Print("[toolclient] websocket message payload parsing error: ", err.Error())
			break
		}

		event.Client = &c.id

		if c.manager.ProviderMgr == nil {
			log.Print("[toolclient] no provider found to forward message")
			break
		}
		c.manager.ProviderMgr.SendEvent(&event)
	}

	c.activeHandlerThreads.Done()
}

func (c *ToolClient) messageSender() {
	select {
	case <-c.done:
		return
	case data := <-c.chSend:
		err := c.connection.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Print("[toolclient] websocket message sending error: ", err.Error())
			break
		}
	}

	c.activeHandlerThreads.Done()
}
