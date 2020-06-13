package nng

import (
	"nanomsg.org/go/mangos/v2"
)

func NewPublisher(addr string) *Publisher {
	pubSock, err := InitPubSock(addr)
	if err != nil {
		return nil
	}

	p := &Publisher{}
	p.pubSock = pubSock
	return p
}

type Publisher struct {
	pubSock mangos.Socket
}

// 发布数据
func (p *Publisher) SendPub(msg []byte) error {
	if err := p.pubSock.Send(msg); err != nil {
		return err
	}
	return nil
}
