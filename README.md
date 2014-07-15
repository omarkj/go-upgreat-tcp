# UpGreat to TCP

"Upgrade" a HTTP connection to a TCP socket

## What is it

Some PaaSs[1] and servers support HTTP Upgrade, and some allow you to upgrade
to whatever you like. This tiny library tries to upgrade to a `upgreat-tcp`,
making it possible to go through the HTTP handshake (often important for routing)
but end up with a TCP socket usable for whatever.

You can probably find some use for this.

## How does it work

This library only handles handshaking, it does not make outbound connections. That
means you should be able to use it with Go's Net and TLS sockets. It also allows
contains a server-side component that will return a raw socket after a successful
handshake.

### On the client

```
conn, err := net.Dial("tcp", "myapp.herokuapp.com:80")
conn, err = ClientHandshake(conn, "myapp.herokuapp.com", "", "", nil)
```

### On the server

```
func serve(w http.ResponseWriter, req *http.Request) {
	conn, err := Attach(w, req)
}
```

## Help!

I don't really write Go. Help me make this code less bad.

[1] At least Heroku