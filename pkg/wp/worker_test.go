package wp

import (
	"fmt"
	"testing"
)

// 工作接口
type WorkerA struct {
	name string
}

func (w WorkerA) Do() {
	fmt.Printf("WorkerA [%v] Do() called\n", w.name)
}

// 工作接口
type WorkerNeedResponse struct {
	name string
	ch   chan interface{}
}

func (w WorkerNeedResponse) Do() {
	fmt.Printf("WorkerNeedResponse [%v] Do() called\n", w.name)
	w.ch <- w.name
}

func TestWorkerPool(t *testing.T) {
	wp := InitDefaultWorkerPool(2)

	a := WorkerA{name: "a"}
	wp.Put(a)

	b := WorkerA{name: "b"}
	GetDefaultWorkerPool().Put(b)
	GetDefaultWorkerPool().Put(a)

	c := WorkerNeedResponse{name: "c", ch: make(chan interface{})}
	GetDefaultWorkerPool().Put(c)

	r := <-c.ch
	s, ok := r.(string)
	if !ok {
		t.Errorf("WorkerNeedResponse response failed")
	}
	if s != "c" {
		t.Errorf("WorkerNeedResponse response failed")

	}
}
