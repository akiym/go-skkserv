package skkserv

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const SkkServVersion = "0.0.1"

type Handler interface {
	Request(text string) ([]string, error)
}

type SKKServ struct {
	Port    string
	Handler Handler
}

func NewServer(port string, handler Handler) *SKKServ {
	server := &SKKServ{
		Port:    port,
		Handler: handler,
	}
	return server
}

func (s *SKKServ) Run() {
	ln, err := net.Listen("tcp", s.Port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *SKKServ) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		c, err := reader.ReadByte()
		if err != nil {
			return
		}
		switch c {
		case '0':
			return // end of connection
		case '1':
			buf, err := reader.ReadBytes(' ')
			if err != nil {
				return
			}
			s.serverRequest(conn, buf)
		case '2':
			s.serverVersion(conn)
		case '3':
			s.serverHost(conn)
		}
	}
}

// "1eee " eee is keyword in EUC code with ' ' at the end
func (s *SKKServ) serverRequest(conn net.Conn, buf []byte) error {
	var resp string
	text := string(buf[:len(buf)-1])
	words, err := s.Handler.Request(text)
	if err != nil {
		// server error
		resp = "0"
	} else {
		if len(words) > 0 {
			// found
			resp = "1/" + strings.Join(words, "/") + "/"
		} else {
			// not found
			resp = "4" + text
		}
	}
	if _, err := conn.Write([]byte(resp + "\n")); err != nil {
		return err
	}
	return nil
}

// "2" skkserv version number
func (s *SKKServ) serverVersion(conn net.Conn) error {
	if _, err := conn.Write([]byte(SkkServVersion + " ")); err != nil {
		return err
	}
	return nil
}

// "3" hostname and its IP addresses
func (s *SKKServ) serverHost(conn net.Conn) error {
	addr := conn.LocalAddr().String()
	if _, err := conn.Write([]byte(addr + " ")); err != nil {
		return err
	}
	return nil
}
