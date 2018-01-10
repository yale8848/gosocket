// Create by Yale 2018/1/9 9:25
package gosocket

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Session struct {
	Connect       *ServerConnect
	ID            int64
	Ext           interface{}
	writeChannel  chan *[]byte
	handleChannel chan *[]byte
}

type ServerHandler interface {
	Connect(*Session)
	ReadData(*Session, *[]byte, int) bool

	Close(*Session)
	AcceptError(error)
	ReadTimeout(error)
}

type ServerConnect struct {
	net.Conn
}
type Server struct {
	Handler ServerHandler
}

type Config struct {
	Network           string
	Address           string
	NetworkListen     string
	ReadTimeout       int
	WriteTimeout      int
	WriteChannelSize  int
	HandleChannelSize int
}

func (session *Session) WriteData(bytes *[]byte) {
	session.writeChannel <- bytes
}
func (session *Session) HandleData(bytes *[]byte) {
	session.handleChannel <- bytes
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func NewServer(handler ServerHandler) *Server {
	return &Server{Handler: handler}
}
func defaultConfig(config *Config) {

	if config.ReadTimeout == 0 {
		config.ReadTimeout = 20000
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 20000
	}
}
func (server *Server) Start(config *Config) {

	defaultConfig(config)

	tcpAdd, error := net.ResolveTCPAddr(config.Network, config.Address)
	checkError(error)

	listen, error := net.ListenTCP(config.NetworkListen, tcpAdd)
	checkError(error)

	fmt.Println("start at [ " + config.Address + " ]")
	for {
		conn, error := listen.Accept()
		if error != nil {
			fmt.Println(config.Address + ":" + error.Error())
			continue
		}
		session := &Session{Connect: &ServerConnect{Conn: conn}, writeChannel: make(chan *[]byte, config.WriteChannelSize),
			handleChannel: make(chan *[]byte, config.HandleChannelSize)}
		server.Handler.Connect(session)

		go server.readRoutine(session)
		go server.handleRoutine(session)
		go server.writeRoutine(session)
	}
}
func (server *Server) readRoutine(session *Session) {

	buff := make([]byte, 1024)
	for {
		n, error := session.Connect.Conn.Read(buff)
		if error != nil {
			fmt.Println(error)
			if error == io.EOF {
				continue
			}
			if e, ok := error.(net.Error); ok && e.Timeout() {
				server.Handler.ReadTimeout(error)
			}
			server.Handler.Close(session)
			session.writeChannel <- nil
			session.handleChannel <- nil
			return
		}
		finish := server.Handler.ReadData(session, &buff, n)
		if finish {

		}

	}
}
func (server *Server) handleRoutine(session *Session) {
	for {
		select {
		case bytes := <-session.handleChannel:
			if bytes == nil {
				fmt.Println("handleConnect stop")
				return
			}
			session.WriteData(bytes)
		}
	}
}
func (server *Server) writeRoutine(session *Session) {

	for {
		select {
		case bytes := <-session.writeChannel:
			if bytes == nil {
				fmt.Println("writeConnect stop")
				return
			}
			session.Connect.Conn.Write(*bytes)
		}
	}

}
