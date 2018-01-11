// Create by Yale 2018/1/11 18:14
package main

import (
	"fmt"
	"gosocket"
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
