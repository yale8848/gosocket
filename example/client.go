// Create by Yale 2018/1/11 17:57
package main

import (
	"fmt"
	"gosocket"
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
	connection.Write(protocol.Encode([]byte("hellowd123456789012458888888465454488889448")))

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
	ticker := time.Tick(20 * time.Second)

	for {
		select {
		case <-ticker:
			con.Write(protocol.Encode(nil))

		}
	}
}

func receive(con *net.TCPConn, protocol *gosocket.Protocol, sig chan bool) {

	buff := make([]byte, 1024)
	for {
		n, error := con.Read(buff)
		if error != nil {
			fmt.Println(error)
			continue
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
	sig <- true
}
