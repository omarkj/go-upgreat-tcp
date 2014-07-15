// This tiny library tries to upgrade to a `upgreat-tcp`,
// making it possible to go through the HTTP handshake (often important for routing)
// but end up with a TCP socket usable for whatever.
package upgreat
import (
	"bytes"
	"errors"
	"strings"
	"net"
	"net/http"
)

// ClientHandshake does the client-side Upgrade handshake returning 
// a Connection
func ClientHandshake(conn net.Conn, host string, verb string,
	path string, headers map[string]string) (net.Conn, error) {
	request := bytes.NewBuffer(nil)
	addVerb(verb, request)
	addPath(path, request)
	addHost(host, request)
	addHeaders(headers, request)
	request.WriteString("\r\n")
	reqstr := request.Bytes()
	conn.Write(reqstr)
	data := make([]byte, 79)
	_, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Attach does the server-side Upgrade handshake returning a Connection
func Attach(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	if strings.ToUpper(r.Header.Get("Upgrade")) != "UPGREAT-TCP" {
		return nil, errors.New("No or invalid upgrade header")
	}
	if strings.ToUpper(r.Header.Get("Connection")) != "UPGRADE" {
		return nil, errors.New("No upgrade header")
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("Hijacking not supported")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, err
	}
	bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	bufrw.WriteString("upgrade: upgreat-tcp\r\n")
	bufrw.WriteString("connection: upgrade\r\n")
	bufrw.WriteString("\r\n")
	bufrw.Flush()
	return conn, nil
}

func addVerb(verb string, request *bytes.Buffer) {
	if verb != "" {
		request.WriteString(verb)
		request.WriteString(" ")
	} else {
		request.WriteString("GET ")
	}
}

func addPath(path string, request *bytes.Buffer) {
	if path != "" {
		request.WriteString(path)
		request.WriteString(" ")
	} else {
		request.WriteString("/ ")
	}
	request.WriteString("HTTP/1.1\r\n")
}

func addHost(host string, request *bytes.Buffer) {
	request.WriteString("host: ")
	request.WriteString(host)
	request.WriteString("\r\n")
}

func addHeaders(headers map[string]string, request *bytes.Buffer) {
	for key, val := range headers {
		request.WriteString(key)
		request.WriteString(": ")
		request.WriteString(val)
		request.WriteString("\r\n")
	}
	request.WriteString("upgrade: upgreat-tcp\r\n")
	request.WriteString("connection: upgrade\r\n")
}
