package nng

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"nanomsg.org/go/mangos/v2"
	"time"
)

type RpcReq struct {
	ID     string                 `json:"id"`
	CMD    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

type RpcResp struct {
	ID   string      `json:"id"`
	CMD  string      `json:"cmd"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"` // msg 和 data互斥，只会存在一个
	Data interface{} `json:"data"`
}

func NewRPCServer(rpcAddr string) *NngRPCServer {
	rpc := &NngRPCServer{cmds: make(map[string]RpcHandler)}
	repSock, err := InitRepSock(rpcAddr)
	if err != nil {
		return nil
	}
	rpc.sock = repSock

	return rpc
}

type RpcHandler func(map[string]interface{}) interface{}

type NngRPCServer struct {
	cmds map[string]RpcHandler
	sock mangos.Socket
}

// 注册命令
func (nr *NngRPCServer) Register(rpcCmd string, handler RpcHandler) {
	nr.cmds[rpcCmd] = handler
}

// 处理命令
func (nr *NngRPCServer) HandleRequest() {
	for {
		var err error
		var msgByte []byte

		if msgByte, err = nr.sock.Recv(); err != nil {
			time.Sleep(time.Nanosecond)
			continue
		}

		resp := make(map[string]interface{})
		req := make(map[string]interface{})
		err = json.Unmarshal(msgByte, &req)
		if err != nil {
			resp["code"] = 201
			resp["msg"] = "请求错误"
			nngSend(nr.sock, resp)
			continue
		}
		rpcID, ok := req["id"].(string)
		if !ok {
			resp["code"] = 202
			resp["msg"] = "id错误"
			nngSend(nr.sock, resp)
			continue
		}
		resp["id"] = rpcID

		rpcCmd, ok := req["cmd"].(string)
		if !ok {
			resp["code"] = 203
			resp["msg"] = "cmd错误"
			nngSend(nr.sock, resp)
			continue
		}
		resp["cmd"] = rpcCmd

		rpcArgs, ok := req["params"].(map[string]interface{})
		if !ok {
			resp["code"] = 204
			resp["msg"] = "params错误"
			nngSend(nr.sock, resp)
			continue
		}

		handler, ok := nr.cmds[rpcCmd]
		if !ok {
			resp["code"] = 205
			resp["msg"] = "不支持的命令"
			nngSend(nr.sock, resp)
			continue
		}

		result := handler(rpcArgs)
		resp["code"] = 200
		resp["result"] = result
		nngSend(nr.sock, resp)
	}
}

func nngSend(sock mangos.Socket, resp map[string]interface{}) {
	respData, err := json.Marshal(resp)
	if err != nil {
		respData = []byte(fmt.Sprintf("{\"error\":\"%v\"}", err))
	}

	sock.Send(respData)
}

func NewRPCClient(rpcAddr string) *NngRPCClient {
	rpc := &NngRPCClient{}
	reqSock, err := InitReqSock(rpcAddr)
	if err != nil {
		return nil
	}
	rpc.sock = reqSock

	return rpc
}

type NngRPCClient struct {
	sock mangos.Socket
}

func (nc *NngRPCClient) Call1(rpcCmd string, params map[string]interface{}) (result map[string]interface{}, err error) {
	req := make(map[string]interface{})
	req["cmd"] = rpcCmd
	req["id"] = uuid.New().String()
	req["params"] = params

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("json marshal数据失败:%+v", req)
	}

	respData, err := SendAndReceive(nc.sock, reqData)
	if err != nil {
		return nil, err
	}

	result = make(map[string]interface{})
	err = json.Unmarshal(respData, &result)
	return
}

func (nc *NngRPCClient) Call2(rpcCmd string, params interface{}, result interface{}) (err error) {
	req := make(map[string]interface{})
	req["cmd"] = rpcCmd
	req["id"] = uuid.New().String()
	req["params"] = params

	reqData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json marshal数据失败:%+v", req)
	}

	respData, err := SendAndReceive(nc.sock, reqData)
	if err != nil {
		return err
	}

	// {"id":"476eff86-ed89-4d39-ae54-a26b235fcd17","cmd":"ticker.subscribe","code":200,"error":"","data":null}
	type response struct {
		ID   string          `json:"id"`
		Cmd  string          `json:"cmd"`
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}

	var res response
	err = json.Unmarshal(respData, &res)
	if err != nil {
		return err
	}

	//code, ok := res.ID
	if res.Code != 200 {
		return fmt.Errorf("%v", res.Msg)
	}

	err = json.Unmarshal(res.Data, &result)
	return
}
