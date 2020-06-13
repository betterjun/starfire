package mq

import "sync"

type MessageQueue interface {
	Send(message interface{})
	Receive() (message interface{})
	GetChan() chan interface{}
}

var mqMap sync.Map

// 创建MessageQueue
func CreateMQ(mqName string, cacheSize int) MessageQueue {
	val, ok := mqMap.Load(mqName)
	if ok {
		return val.(MessageQueue)
	}

	mq := &messageQueue{
		Name: mqName,
		Size: cacheSize,
		ch:   make(chan interface{}, cacheSize),
	}
	mqMap.Store(mqName, mq)
	return mq
}

// 获取MessageQueue
func GetMQ(mqName string) MessageQueue {
	val, ok := mqMap.Load(mqName)
	if !ok {
		return nil
	}
	return val.(MessageQueue)
}

type messageQueue struct {
	Name string
	Size int
	ch   chan interface{}
}

func (mq *messageQueue) Send(message interface{}) {
	mq.ch <- message
}

func (mq *messageQueue) Receive() (message interface{}) {
	message = <-mq.ch
	return message
}

func (mq *messageQueue) GetChan() chan interface{} {
	return mq.ch
}
