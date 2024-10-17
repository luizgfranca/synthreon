package connectionmgr

import (
	"encoding/json"
	"log"
	"platformlab/controlpanel/model"
	"sync"

	"github.com/gorilla/websocket"
)

type Provider struct {
	connection           *websocket.Conn
	chSend               chan []byte
	chRecv               chan string
	done                 chan bool
	activeHandlerThreads sync.WaitGroup
}

func NewProviderConnectionMgr(connection *websocket.Conn) *Provider {
	provider := Provider{
		connection:           connection,
		chSend:               make(chan []byte),
		chRecv:               make(chan string),
		done:                 make(chan bool),
		activeHandlerThreads: sync.WaitGroup{},
	}

	provider.activeHandlerThreads.Add(2)

	go provider.providerMessageReceiver()
	go provider.providerMessageSender()

	return &provider
}

func (p *Provider) SendEvent(event *model.ToolEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Print("[toolprovider] error encoding response: ", err.Error())
	}

	p.chSend <- data
}

func (p *Provider) Close() {
	p.done <- false
	p.activeHandlerThreads.Wait()

	p.connection.Close()
}

func (p *Provider) providerMessageReceiver() {
	defer p.connection.Close()
	var event model.ToolEvent

	for {
		select {
		case <-p.done:
			return
		default:

		}

		msgtype, message, err := p.connection.ReadMessage()
		if err != nil {
			log.Print("[toolprovider] websocket message receiving error: ", err.Error())
			break
		}

		if msgtype != websocket.TextMessage {
			log.Print("[toolprovider] unexpeted message type, type sould be a textMessage")
			break
		}

		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Print("[toolprovider] websocket message payload parsing error: ", err.Error())
			break
		}

		log.Print("[toolclient] EVENT: class ", event.Class)

		if event.Class == model.EventClassAnnouncement {
			ack := model.NewAnnouncementAckEvent(&event)
			data, err := json.Marshal(ack)
			if err != nil {
				log.Print("[toolprovider] error encoding response: ", err.Error())
			}

			err = p.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Print("[toolprovider] websocket message sending error: ", err.Error())
				break
			}
		}
	}
}

func (p *Provider) providerMessageSender() {
	select {
	case <-p.done:
		return
	case data := <-p.chSend:
		err := p.connection.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Print("[toolprovider] websocket message sending error: ", err.Error())
			break
		}
	}
}
