package toolentity

import (
	"log"
	tooleventmodule "synthreon/modules/toolevent"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebsocketToolEntity struct {
	id                   string
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
		id:                   uuid.NewString(),
		connection:           connection,
		chSend:               make(chan []byte),
		chRecv:               make(chan string),
		done:                 make(chan bool),
		activeHandlerThreads: sync.WaitGroup{},
	}
	e.log("new connection")
	return &e
}

func (e *WebsocketToolEntity) StartHandler() error {
	e.log("starting handlers")

	e.activeHandlerThreads.Add(2)
	go e.messageSenderThread()
	go e.messageReceiverThread()

	return nil
}

func (e *WebsocketToolEntity) OnEventReceived(handler func(event *tooleventmodule.ToolEvent)) {
	e.eventReceivedCallback = handler
}

func (e *WebsocketToolEntity) SendEvent(event *tooleventmodule.ToolEvent) error {
	e.log("preparing event to be sent", event)
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
	e.log("draining connections")
	e.done <- true
	// e.done <- true
	e.connection.Close()
	e.activeHandlerThreads.Wait()
	e.log("closing websocket")
}

func (e *WebsocketToolEntity) log(v ...any) {
	x := append([]any{"[WebsocketToolEntity-" + e.id + "]"}, v...)

	log.Println(x...)
}

func (e *WebsocketToolEntity) messageSenderThread() {
	e.log("starting messageSender loop")

	for {
		e.log("waiting for new messages to send")

		select {
		case <-e.done:
			e.log("stopping sender")
			e.activeHandlerThreads.Done()
			return
		case data := <-e.chSend:
			e.log("sending message: ", string(data))
			err := e.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				e.log("websocket message sending error: ", err.Error())
				return
			}
		}
	}
}

func (e *WebsocketToolEntity) messageReceiverThread() {
	e.log("starting messageReceiver loop")
	if e.eventReceivedCallback == nil {
		log.Fatal("[WebsocketToolEntity] Assumption violation: No event handler added for WebSocketToolEntity.")
	}

	for {
		select {
		case <-e.done:
			e.log("stopping receiver")
			e.activeHandlerThreads.Done()
			return
		default:
			// continue loop if not marked as done
		}

		e.log("waiting for messages")
		// FIXME: e.connection.SetReadDeadline(time.Now().Add(time.Second))
		msgtype, message, err := e.connection.ReadMessage()
		e.log("message received")
		if err != nil {
			e.log("websocket message receiving error: ", err.Error())
			// TODO: should handle in a smarter way different websocket message errors
			// TODO: should think of a way to pass down more information
			// 		 about the reason of the disconnection
			if e.disconnectedCallback != nil {
				e.disconnectedCallback()
			}

			// return instead of break here to avoid duplicate disconnection
			// notifications
			e.activeHandlerThreads.Done()
			return
		}

		if msgtype != websocket.TextMessage {
			e.log("unexpeted message type, type sould be a textMessage")
			break
		}

		messageString := string(message)

		// TODO: when the event is invalid, should i just drop it like this?
		event, err := tooleventmodule.ParseEventString(&messageString)
		if err != nil {
			e.log("websocket message payload parsing error: ", err.Error())
			break
		}

		e.log("event received: ", event)
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
