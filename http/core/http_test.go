package http

import (
	"bufio"
	"net"
	"net/http"
	"testing"
)

func TestTinyHttp(t *testing.T) {
	go func() {
		//server
		lis, err := net.Listen("tcp4", "127.0.0.1:3000")
		if err != nil {
			t.Errorf("listen error:%s", err)
		}
		// for {
		conn, err := lis.Accept()
		if err != nil {
			t.Errorf("accept error:%s", err)
			// break
		}
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			t.Logf("got from conn:%s", scanner.Text())
		}
		/*
			reader := bufio.NewReader(conn)

			bytes, err := reader.ReadBytes('\n')

			if err != nil {
				t.Errorf("read bytes error:%s", err)
			}
			t.Logf("got from conn:%s", string(bytes))
		*/
		// }
	}()

	//client
	_, err := http.Get("http://127.0.0.1:3000/test")
	if err != nil {
		t.Errorf("http get error:%s", err)
	}

}
