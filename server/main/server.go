// Create by Yale 2018/1/9 14:10
package main

import (
	"fmt"
	"gosocket"
	"time"
)

type ServerHandlerImp struct {
}

func (s *ServerHandlerImp) Connect(session *gosocket.Session) {

	fmt.Println("Connect : ")

}
func (s *ServerHandlerImp) ReadData(session *gosocket.Session, bytes *[]byte, n int) bool {
	fmt.Println("ReadData : " + string((*bytes)[0:n]))
	bt := []byte("aaa")
	session.HandleData(&bt)
	return true
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

	sh := &ServerHandlerImp{}
	server := gosocket.NewServer(sh)
	server.Start(&gosocket.Config{
		Network: "tcp", Address: ":7777", NetworkListen: "tcp", ReadTimeout: int(time.Second), WriteTimeout: int(time.Second)})

}
