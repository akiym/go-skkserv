package skkserv

import (
	"bufio"
	"net"
	"testing"
)

type TestHandler struct{}

func (s *TestHandler) Request(text string) (words []string, err error) {
	if text != "not-exists" {
		words = []string{"a", "b", "c"}
	}
	return words, nil
}

type Test struct {
	request          string
	delim            byte
	expectedResponse string
}

var tests = []Test{
	{"1eee ", '\n', "1/a/b/c/\n"},
	{"1not-exists ", '\n', "4not-exists\n"},
	{"2", ' ', SkkServVersion + " "},
	{"3", ' ', "127.0.0.1:55100 "},
}

func TestRequest(t *testing.T) {
	server := NewServer(":55100", &TestHandler{})
	go server.Run()

	conn, err := net.Dial("tcp", "localhost:55100")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, test := range tests {
		if _, err := conn.Write([]byte(test.request)); err != nil {
			t.Fatal(err)
		}
		resp, err := bufio.NewReader(conn).ReadString(test.delim)
		if err != nil {
			t.Fatal(err)
		}
		if resp != test.expectedResponse {
			t.Fail()
		}
	}
}
