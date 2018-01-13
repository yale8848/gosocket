# gosocket

Simple Golang server socket long connection skeleton

## Useage

```go
go get github.com/yale8848/gosocket
```

- server

```go

package main

import (
	"fmt"
	"github.com/yale8848/gosocket"
)

type ServerHandlerImp struct {
}

func (s *ServerHandlerImp) Connect(session *gosocket.Session) {

	fmt.Println("Connect : ")

}
func (s *ServerHandlerImp) HandleData(session *gosocket.Session, protocol *gosocket.Protocol) {

	if protocol.IsHeartBeat() {
		fmt.Println("ReadData : IsHeartBeat")
		d := protocol.Encode(nil)
		session.WriteData(d)
	} else {
		fmt.Println("ReadData :" + protocol.String())
		d := protocol.Encode(protocol.GetData())
		session.WriteData(d)
	}

}
func (s *ServerHandlerImp) Close(session *gosocket.Session) {

	fmt.Println("Close : ")
}
func (s *ServerHandlerImp) AcceptError(err error) {

	fmt.Println("AcceptError : " + err.Error())
}

func (s *ServerHandlerImp) ReadTimeout(err error) {
	fmt.Println("ReadTimeout : " + err.Error())
}
func main() {
	server := gosocket.NewServer(&ServerHandlerImp{}, &gosocket.Protocol{})
	server.Start(&gosocket.Config{
		Network: "tcp", Address: ":7777", NetworkListen: "tcp", ReadTimeout: 20})

}


```

- client

```go
package main

import (
	"fmt"
	"github.com/yale8848/gosocket"
	"net"
	"time"
)

func main() {
	hawkServer, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.DialTCP("tcp4", nil, hawkServer)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(10 * time.Second)
	protocol := &gosocket.Protocol{
		Version: 1,
		Reserve: 0,
	}
	connection.Write(protocol.Encode([]byte("Hello World")))

	sig := make(chan bool)
	go heartBeat(connection)
	go receive(connection, protocol, sig)
	<-sig
	connection.Close()

}

func heartBeat(con *net.TCPConn) {
	protocol := &gosocket.Protocol{
		Version: 1,
		Reserve: gosocket.HEART_BEAT,
	}
	ticker := time.Tick(10 * time.Second)

	count := 0
	for {
		select {
		case <-ticker:
			con.Write(protocol.Encode(nil))
			count++
			if count == 3 {
				return
			}

		}
	}
}

func receive(con *net.TCPConn, protocol *gosocket.Protocol, sig chan bool) {

	buff := make([]byte, 1024)
	for {
		n, error := con.Read(buff)
		if error != nil {
			fmt.Println(error)
			break
		}
		finish := protocol.Decode(buff[0:n])
		if finish {
			if protocol.IsHeartBeat() {
				fmt.Println("receive from server: IsHeartBeat")
			} else {
				fmt.Println("receive from server: " + string(protocol.GetData()))
			}
		}

	}
	fmt.Printf("read finish")
	sig <- true
}


```

