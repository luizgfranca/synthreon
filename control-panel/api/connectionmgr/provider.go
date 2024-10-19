package connectionmgr

import (
	"encoding/json"
	"log"
	"platformlab/controlpanel/model"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type ProviderMgr struct {
	connection           *websocket.Conn
	chSend               chan []byte
	chRecv               chan string
	done                 chan bool
	activeHandlerThreads sync.WaitGroup
	ClientMgr            *ClientMgr
}

func NewProviderConnectionMgr(connection *websocket.Conn) *ProviderMgr {
	provider := ProviderMgr{
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

func (p *ProviderMgr) SendEvent(event *model.ToolEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Print("[toolprovider] error encoding response: ", err.Error())
	}

	log.Print("[toolprovider] queueing event to be sent")
	p.chSend <- data
}

func (p *ProviderMgr) Close() {
	p.done <- false
	p.activeHandlerThreads.Wait()

	p.connection.Close()
}

func (p *ProviderMgr) providerMessageReceiver() {
	var event model.ToolEvent

	log.Print("[toolprovider] starting providerMessageReceiver loop")
	for {
		select {
		case <-p.done:
			return
		default:

		}

		log.Print("[toolprovider] waiting for messages")
		msgtype, message, err := p.connection.ReadMessage()
		log.Print("[toolprovider] message received")
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

		log.Print("[toolprovider] event: ", event)
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

			log.Print("[toolprovider] new provider registration acknowleged")
		} else {
			if event.Client == "" {
				log.Print("[toolprovider] no specified client to forward event")
				break
			}

			idx, err := strconv.ParseUint(event.Client, 10, 64)
			if err != nil {
				log.Print("[toolprovider] invalid client specified")
				break
			}

			log.Print("[toolprovider] forwarding event to client: ", event.Client)

			if p.ClientMgr == nil {
				panic("no clientManager found for provider")
			}
			p.ClientMgr.SendEvent(uint(idx), &event)
		}
	}

	p.activeHandlerThreads.Done()
}

func (p *ProviderMgr) providerMessageSender() {
	log.Print("[toolprovider] starting providerMessageSender loop")

	for {
		log.Print("[toolprovider] waiting for new messages")

		select {
		case <-p.done:
			p.activeHandlerThreads.Done()
			return
		case data := <-p.chSend:
			err := p.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Print("[toolprovider] websocket message sending error: ", err.Error())
				return
			}
		}
	}
}
