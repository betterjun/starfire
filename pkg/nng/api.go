package nng

import (
	"log"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pair"
	"nanomsg.org/go/mangos/v2/protocol/pub"
	"nanomsg.org/go/mangos/v2/protocol/pull"
	"nanomsg.org/go/mangos/v2/protocol/push"
	"nanomsg.org/go/mangos/v2/protocol/rep"
	"nanomsg.org/go/mangos/v2/protocol/req"
	"nanomsg.org/go/mangos/v2/protocol/sub"
	_ "nanomsg.org/go/mangos/v2/transport/all"
)

// 初始化命令套接字
func InitReqSock(url string) (sock mangos.Socket, err error) {
	if sock, err = req.NewSocket(); err != nil {
		return
	}

	if err = sock.Dial(url); err != nil {
		return nil, err
	}
	return
}

// 初始化命令服务端套接字
func InitRepSock(url string) (sock mangos.Socket, err error) {
	if sock, err = rep.NewSocket(); err != nil {
		return
	}
	if err = sock.Listen(url); err != nil {
		return nil, err
	}
	return
}

// 初始化订阅套接字
func InitSubSock(url string) (sock mangos.Socket, err error) {
	if sock, err = sub.NewSocket(); err != nil {
		return nil, err
	}

	if err = sock.Dial(url); err != nil {
		return nil, err
	}

	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		return nil, err
	}
	return
}

// 初始化订阅服务端套接字
func InitPubSock(url string) (sock mangos.Socket, err error) {
	if sock, err = pub.NewSocket(); err != nil {
		return nil, err
	}
	if err = sock.Listen(url); err != nil {
		return nil, err
	}
	return
}

// 初始化pair的服务端套接字
func InitPairServerSock(url string) (sock mangos.Socket, err error) {
	if sock, err = pair.NewSocket(); err != nil {
		return nil, err
	}
	if err = sock.Listen(url); err != nil {
		return nil, err
	}
	return
}

// 初始化pair的客户端套接字
func InitPairClientSock(url string) (sock mangos.Socket, err error) {
	if sock, err = pair.NewSocket(); err != nil {
		return nil, err
	}
	if err = sock.Dial(url); err != nil {
		return nil, err
	}
	return
}

// 初始化pipeline的拉取端套接字
func InitPipelinePullSock(url string) (sock mangos.Socket, err error) {
	if sock, err = pull.NewSocket(); err != nil {
		return nil, err
	}
	if err = sock.Listen(url); err != nil {
		return nil, err
	}
	return
}

// 初始化pipeline的推送端套接字
func InitPipelinePushSock(url string) (sock mangos.Socket, err error) {
	if sock, err = push.NewSocket(); err != nil {
		return nil, err
	}
	if err = sock.Dial(url); err != nil {
		return nil, err
	}
	return
}

// 发送和读取数据
func SendAndReceive(sock mangos.Socket, reqData []byte) (respData []byte, err error) {
	if err = sock.Send(reqData); err != nil {
		return
	}
	log.Println("数据已成功发送", string(reqData))

	respData, err = sock.Recv()
	log.Println("数据已接收", string(respData))

	return
}
