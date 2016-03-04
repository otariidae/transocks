// +build linux

package transocks

import (
	"net"
	"testing"
)

func TestGetOriginalDST(t *testing.T) {
	t.Skip()

	l, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 1081})
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	orig_addr, err := GetOriginalDST(c.(*net.TCPConn))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(orig_addr.String())
}
