package sproto_test

import (
	"testing"

	"errors"
	"reflect"

	"github.com/xjdrew/gosproto"
)

type FoobarRequest struct {
	What *string `sproto:"string,0,name=what"`
}

type FoobarResponse struct {
	Ok   *bool   `sproto:"boolean,0,name=ok"`
	What *string `sproto:"string,1,name=what"`
}

type FooResponse struct {
	Ok *bool `sproto:"boolean,0,name=ok"`
}

var protocols []*sproto.Protocol = []*sproto.Protocol{
	&sproto.Protocol{
		Type:       1,
		Name:       "test.foobar",
		MethodName: "Test.Foobar",
		Request:    reflect.TypeOf(&FoobarRequest{}),
		Response:   reflect.TypeOf(&FoobarResponse{}),
	},
	&sproto.Protocol{
		Type:       2,
		Name:       "test.foo",
		MethodName: "Test.Foo",
		Response:   reflect.TypeOf(&FooResponse{}),
	},
	&sproto.Protocol{
		Type:       3,
		Name:       "test.bar",
		MethodName: "Test.Bar",
	},
}

func checkRequest(rpc *sproto.Rpc, name string, session int32, sp interface{}) (interface{}, error) {
	chunk, err := rpc.RequestEncode(name, session, sp)
	if err != nil {
		return nil, err
	}

	mode, name1, session1, sp1, err := rpc.Dispatch(chunk)
	if err != nil {
		return nil, err
	}
	if mode != sproto.RpcRequestMode || name != name1 || session != session1 {
		return nil, errors.New("dipatch failed: unmatch meta info")
	}
	return sp1, nil
}

func TestRpcRequest(t *testing.T) {
	rpc, err := sproto.NewRpc(protocols)
	if err != nil {
		t.Fatalf("new rpc failed with error:%s", err)
	}

	// request & dispatch request
	what := "hello"
	sp, err := checkRequest(rpc, "test.foobar", 1, &FoobarRequest{What: &what})
	if err != nil {
		t.Fatalf("check request failed:%s", err)
	}
	req := sp.(*FoobarRequest)
	if *req.What != what {
		t.Fatalf("check failed: unmatch data")
	}

	// nil request
	sp, err = checkRequest(rpc, "test.foo", 2, nil)
	if err != nil {
		t.Fatalf("check request failed:%s", err)
	}
	if sp != nil {
		t.Fatalf("check failed: unmatch data")
	}

	// nil request
	sp, err = checkRequest(rpc, "test.bar", 0, nil)
	if err != nil {
		t.Fatalf("check request failed:%s", err)
	}
	if sp != nil {
		t.Fatalf("check failed: unmatch data")
	}
}

func TestRpcResponse(t *testing.T) {
	rpc, err := sproto.NewRpc(protocols)
	if err != nil {
		t.Fatalf("new rpc failed with error:%s", err)
	}

	// request & dispatch request
	what := "hello"
	_, err = rpc.RequestEncode("test.foobar", 18, &FoobarRequest{What: &what})
	if err != nil {
		t.Fatalf("request encode failed:%s", err)
	}

	// check response
	chunk, err := rpc.ResponseEncode("test.foobar", 18, &FoobarResponse{Ok: sproto.Bool(true)})
	if err != nil {
		t.Fatalf("response encode failed:%s", err)
	}

	mode, name, session, sp, err := rpc.Dispatch(chunk)
	if err != nil {
		t.Fatalf("dispatch failed:%s", err)
	}
	if mode != sproto.RpcResponseMode || name != "test.foobar" || session != 18 {
		t.Fatalf("dispatch failed:unmatch meta info")
	}
	response := sp.(*FoobarResponse)
	if !*response.Ok {
		t.Fatalf("dispatch failed:unmatch data")
	}
}
