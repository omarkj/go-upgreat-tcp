package upgreat
import (
	"testing"
	"time"
	"bytes"
	"net"
	"net/http"
)

func TestClientHandshake(t *testing.T) {
	server := initTest()
	conn, err := net.Dial("tcp", server)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		panic(err)
	}
	conn, err = ClientHandshake(conn, "localhost", "", 
		"", nil)
	if err != nil {
		panic(err)
	}
	msg := []byte("hi")
	conn.Write(msg)
	data := make([]byte, 2)
	_, err = conn.Read(data)
	if bytes.Compare(msg, data) != 0 {
		panic("Messages do not compare")
	}
	
}

func initTest() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go http.Serve(listener, http.HandlerFunc(serve))
	return listener.Addr().String()
}

func serve(w http.ResponseWriter, req *http.Request) {
	conn, err := Attach(w, req)
	if err != nil {
		panic(err)
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	data := make([]byte, 512)
	_, err = conn.Read(data)
	if err != nil {
		panic(err)
	}
	conn.Write(data)
	conn.Close()
}
