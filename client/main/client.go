// Create by Yale 2018/1/9 14:13
package main

import (
	"fmt"
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
	fmt.Println("Write hello ")
	connection.Write([]byte("hellow"))

	sig := make(chan bool)
	go receive(connection, sig)
	<-sig
	connection.Close()

}
func receive(con *net.TCPConn, sig chan bool) {
	fmt.Println("receive ")
	buff := make([]byte, 1024)
	for {
		n, error := con.Read(buff)
		if error != nil {
			fmt.Println(error)
			continue
		}
		if n < 1024 {
			fmt.Println("receive form server :" + string(buff[0:n]))
		}
	}
	sig <- true
}
