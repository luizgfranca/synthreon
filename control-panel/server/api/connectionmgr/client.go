package connectionmgr

import (
	"encoding/json"
	"log"
	"platformlab/controlpanel/model"
	genericmodule "platformlab/controlpanel/modules/commonmodule"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ToolClient struct {
	id                   string
	connection           *websocket.Conn
	chSend               chan []byte
	done                 chan bool
	activeHandlerThreads sync.WaitGroup
	manager              *ClientMgr
}

type ClientMgr struct {
	clients             []*ToolClient
	ProviderMgr         *ProviderMgr
	clientAdditionMutex sync.Mutex
}

func (m *ClientMgr) NewClient(connection *websocket.Conn) string {
	m.clientAdditionMutex.Lock()
	idx := strconv.Itoa(len(m.clients))
	m.clients = append(m.clients, NewToolClient(m, idx, connection))
	m.clientAdditionMutex.Unlock()
	log.Print("[ToolClient] Registered new client with ID: ", idx)

	return idx
}

func (m *ClientMgr) SendEvent(client uint, e *model.ToolEvent) error {
	m.clientAdditionMutex.Lock()
	log.Print("[ToolClient] clents: ", m.clients)
	if len(m.clients) <= int(client) {
		m.clientAdditionMutex.Unlock()
		return &genericmodule.GenericLogicError{Message: "client not found"}
	}
	m.clientAdditionMutex.Unlock()

	m.clients[client].SendEvent(e)
	return nil
}

func NewToolClient(manager *ClientMgr, id string, connection *websocket.Conn) *ToolClient {
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

	log.Print("[toolprovider] enqueuing message to be sent to client")
	c.chSend <- data
}

func (c *ToolClient) Close() {
	c.done <- false
	c.activeHandlerThreads.Wait()

	c.connection.Close()
}

func (c *ToolClient) messageReceiver() {
	var event model.ToolEvent

	log.Print("[toolclient] starting providerMessageReceiver loop")
	for {
		select {
		case <-c.done:
			return
		default:
			// should continue
		}

		log.Print("[toolclient] waiting for messages for client: ", c.id)
		msgtype, message, err := c.connection.ReadMessage()
		log.Print("[toolclient] message received from client: ", c.id)
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

		event.Client = c.id

		for c.manager.ProviderMgr == nil {
			log.Print("[toolclient] no provider found to forward message")
			time.Sleep(1 * time.Second)
		}

		log.Print("[toolclient] forwarding event to provider ", event)

		if c.manager.ProviderMgr == nil {
			panic("[toolclient] provider null when going to forward message")
		}
		c.manager.ProviderMgr.SendEvent(&event)
	}

	c.activeHandlerThreads.Done()
}

func (c *ToolClient) messageSender() {
	log.Print("[toolclient] starting messageSender loop for client: ", c.id)

	for {
		select {
		case <-c.done:
			c.activeHandlerThreads.Done()
			return
		case data := <-c.chSend:
			log.Print("[toolclient] sending message to client: ", c.id)
			err := c.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Print("[toolclient] websocket message sending error: ", err.Error())
				return
			}
		}
	}
}
