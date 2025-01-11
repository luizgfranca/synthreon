package toolentity

import (
	"log"
	tooleventmodule "platformlab/controlpanel/modules/toolevent"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketToolEntity struct {
	connection           *websocket.Conn
	chSend               chan []byte
	chRecv               chan string
	done                 chan bool
	activeHandlerThreads sync.WaitGroup

	eventReceivedCallback func(event *tooleventmodule.ToolEvent)
	disconnectedCallback  func()
}

func NewWebsocketToolEntity(connection *websocket.Conn) *WebsocketToolEntity {
	e := WebsocketToolEntity{
		connection:           connection,
		chSend:               make(chan []byte),
		chRecv:               make(chan string),
		done:                 make(chan bool),
		activeHandlerThreads: sync.WaitGroup{},
	}
	return &e
}

func (e *WebsocketToolEntity) StartHandler() error {
	e.activeHandlerThreads.Add(2)
	go e.messageSenderThread()
	go e.messageReceiverThread()

	return nil
}

func (e *WebsocketToolEntity) OnEventReceived(handler func(event *tooleventmodule.ToolEvent)) {
	e.eventReceivedCallback = handler
}

func (e *WebsocketToolEntity) SendEvent(event *tooleventmodule.ToolEvent) error {
	data, err := tooleventmodule.WriteV0EventString(event)
	if err != nil {
		return err
	}

	e.chSend <- []byte(*data)

	return nil
}

func (e *WebsocketToolEntity) OnDisconnect(handler func()) {
	e.disconnectedCallback = handler
}

func (e *WebsocketToolEntity) Close() {
	e.done <- false
	e.activeHandlerThreads.Wait()
	e.connection.Close()
}

func (e *WebsocketToolEntity) messageSenderThread() {
	log.Print("[WebsocketToolEnity] starting messageSender loop")

	for {
		log.Print("[WebsocketToolEnity] waiting for new messages")

		select {
		case <-e.done:
			e.activeHandlerThreads.Done()
			return
		case data := <-e.chSend:
			err := e.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Print("[WebsocketToolEnity] websocket message sending error: ", err.Error())
				return
			}
		}
	}
}

func (e *WebsocketToolEntity) messageReceiverThread() {
	log.Print("[WebsocketToolEntity] starting messageReceiver loop")
	if e.eventReceivedCallback == nil {
		log.Fatal("[WebsocketToolEntity] Assumption violation: No event handler added for WebSocketToolEntity.")
	}

	for {
		select {
		case <-e.done:
			return
		default:
			// continue loop if not marked as done
		}

		log.Print("[WebsocketToolEntity] waiting for messages")
		msgtype, message, err := e.connection.ReadMessage()
		log.Print("[WebsocketToolEntity] message received")
		if err != nil {
			log.Print("[WebsocketToolEntity] websocket message receiving error: ", err.Error())

			if _, ok := err.(*websocket.CloseError); ok {
				log.Print("[WebsocketToolEntity] websocket conneciton closed")
				if e.disconnectedCallback != nil {
					// TODO: should think of a way to pass down more information
					// 		 about the reason of the disconnection
					e.disconnectedCallback()
				}
			}
			break
		}

		if msgtype != websocket.TextMessage {
			log.Print("[WebsocketToolEntity] unexpeted message type, type sould be a textMessage")
			break
		}

		messageString := string(message)
		event, err := tooleventmodule.ParseEventString(&messageString)
		if err != nil {
			log.Print("[WebsocketToolEntity] websocket message payload parsing error: ", err.Error())
			break
		}

		log.Print("[WebsocketToolEntity] event received: ", event)
		if e.eventReceivedCallback == nil {
			log.Fatal("[WebsocketToolEntity] Assumption violation: No event handler added for WebSocketToolEntity.")
		}

		e.eventReceivedCallback(event)
	}

	if e.disconnectedCallback != nil {
		e.disconnectedCallback()
	}

	e.activeHandlerThreads.Done()
}
