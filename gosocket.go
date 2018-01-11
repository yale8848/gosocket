// Create by Yale 2018/1/9 9:25
package gosocket

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Session struct {
	Connect       *ServerConnect
	writeChannel  chan []byte
	handleChannel chan []byte
}

type ServerHandler interface {
	Connect(*Session)
	HandleData(*Session, *Protocol)

	Close(*Session)
	AcceptError(error)
	ReadTimeout(error)
}

type ServerConnect struct {
	net.Conn
}
type Server struct {
	Handler  ServerHandler
	protocol *Protocol
	config   *Config
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

func (session *Session) WriteData(bytes []byte) {
	session.writeChannel <- bytes
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func NewServer(handler ServerHandler, protocol *Protocol) *Server {
	return &Server{Handler: handler, protocol: protocol}
}
func (server *Server) defaultConfig(config *Config) {
	server.config = config
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 20
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 20
	}

}
func (server *Server) Start(config *Config) {

	server.defaultConfig(config)

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
		session := &Session{Connect: &ServerConnect{Conn: conn}, writeChannel: make(chan []byte, config.WriteChannelSize),
			handleChannel: make(chan []byte, config.HandleChannelSize)}
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

		session.Connect.Conn.SetReadDeadline(time.Now().Add(time.Duration(server.config.ReadTimeout) * time.Second))

		hb := make([]byte, n)
		copy(hb, buff)
		session.handleChannel <- hb

	}
}
func (server *Server) handleRoutine(session *Session) {
	for {
		select {
		case bytes := <-session.handleChannel:
			if bytes == nil {
				fmt.Println("handleRoutine stop")
				return
			}

			ptcl := server.protocol
			finish := ptcl.Decode(bytes)

			if finish && ptcl.success {

				server.Handler.HandleData(session, server.protocol)
			}
		}
	}
}
func (server *Server) writeRoutine(session *Session) {

	for {
		select {
		case bytes := <-session.writeChannel:
			if bytes == nil {
				fmt.Println("writeRoutine stop")
				return
			}
			session.Connect.Conn.Write(bytes)
		}
	}

}
