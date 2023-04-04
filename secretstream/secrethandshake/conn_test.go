package secrethandshake

import (
	"encoding/base64"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestAuth(t *testing.T) {
	log.SetOutput(os.Stdout)
	/* ERROR
	rwServer := new(io.ReadWriter)
	rwClient := new(io.ReadWriter)*/

	//appkey := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	appkey := make([]byte, 32)
	io.ReadFull(random(0xFF), appkey)
	t.Log("appkey= " + base64.StdEncoding.EncodeToString(appkey))
	keySrv, err := GenEdKeyPair(random(0x11))
	if err != nil {
		t.Fatal(err)
	}
	serverState, err := NewServerState(appkey, *keySrv)
	if err != nil {
		t.Error("NewServerState err:", err)
	}
	t.Log("server public key=", keySrv.Public)
	t.Log("server id: @" + base64.StdEncoding.EncodeToString(keySrv.Public) + ".ed25519")

	keyCli, err := GenEdKeyPair(random(0x22))
	if err != nil {
		t.Fatal(err)
	}
	clientState, err := NewClientState(appkey, *keyCli, keySrv.Public) //client知道来源
	if err != nil {
		t.Error("NewServerState err:", err)
	}
	t.Log("client public key=", keyCli.Public)
	t.Log("client id: @" + base64.StdEncoding.EncodeToString(keyCli.Public) + ".ed25519")
	//构建 server和client的rw, io.Pipe()的r和w是会对当前申请的io锁住w和r，符合场景
	rServer, wClient := io.Pipe()
	rClient, wServer := io.Pipe()
	rwServer := rw{rServer, wServer}
	rwClient := rw{rClient, wClient}

	ch := make(chan error, 2)
	go func() {
		err := ServerShake(serverState, rwServer)
		ch <- err
		//w EOF
		wServer.Close()
	}()
	go func() {
		err := ClientShack(clientState, rwClient)
		ch <- err
		//w EOF
		wClient.Close()
	}()

	if err = <-ch; err != nil {
		t.Errorf("ch-1 read: %v", err)
	}
	if err = <-ch; err != nil {
		t.Errorf("ch-2 read: %v", err)
	}
	t.Log("server secret=", serverState.secret)
	t.Log("client secret=", clientState.secret)
	if reflect.DeepEqual(serverState.secret, clientState.secret) == false {
		t.Error("2 peer secret not equal")
	}
}

type rw struct { //implents io.ReadWriter
	io.Reader
	io.Writer
}

type random byte

// a io.Reader
func (r random) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(r)
	}
	n = len(p)
	return
}
