package transocks

import (
	"io"
	"net"
	"os"
	"testing"
	"time"
)

func TestHTTPDialer(t *testing.T) {
	t.Skip()

	// This test only works if Squid allowing CONNECT to port 80 is
	// running on the local machine on port 3128.

	d := &httpDialer{
		addr:    "127.0.0.1:3128",
		forward: &net.Dialer{Timeout: 5 * time.Second},
	}

	conn, err := d.Dial("tcp", "www.yahoo.com:80")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.Write([]byte("GET / HTTP/1.1\r\nHost: www.yahoo.com:80\r\nConnection: close\r\n\r\n"))
	io.Copy(os.Stdout, conn)
}
