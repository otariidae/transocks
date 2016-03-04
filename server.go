package transocks

import (
	"io"
	"net"
	"time"

	"github.com/cybozu-go/log"
	"golang.org/x/net/proxy"
)

var (
	defaultDialer = &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
	}
)

// Server provides transparent proxy server functions.
type Server struct {
	config   *Config
	dialer   proxy.Dialer
	listener net.Listener
}

// NewServer creates Server.
// If c is not valid, this returns non-nil error.
func NewServer(c *Config) (*Server, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	dialer := defaultDialer
	if c.Dialer != nil {
		dialer = c.Dialer
	}
	proxy_dialer, err := proxy.FromURL(c.ProxyURL, dialer)
	if err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", c.Listen)
	if err != nil {
		return nil, err
	}

	return &Server{c, proxy_dialer, l}, nil
}

// Serve accepts and handles new connections forever.
func (s *Server) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Critical(err.Error(), nil)
			return err
		}
		tcp_conn, ok := conn.(*net.TCPConn)
		if !ok {
			conn.Close()
			panic("not a TCPConn!")
		}
		go s.handleConnection(tcp_conn)
	}
}

func (s *Server) handleConnection(c *net.TCPConn) {
	defer c.Close()

	var addr string

	switch s.config.Mode {
	case ModeNAT:
		orig_addr, err := GetOriginalDST(c)
		if err != nil {
			log.Error(err.Error(), nil)
			return
		}
		addr = orig_addr.String()
	default:
		addr = c.LocalAddr().String()
	}

	if log.Enabled(log.LvDebug) {
		log.Debug("making proxy connection", map[string]interface{}{
			"_dst": addr,
		})
	}

	pconn, err := s.dialer.Dial("tcp", addr)
	if err != nil {
		log.Error(err.Error(), nil)
		return
	}
	defer pconn.Close()

	ch := make(chan error, 2)
	go copyData(c, pconn, ch)
	go copyData(pconn, c, ch)
	for i := 0; i < 2; i++ {
		e := <-ch
		if e != nil {
			log.Error(err.Error(), nil)
			break
		}
	}

	if log.Enabled(log.LvDebug) {
		log.Debug("closing proxy connection", map[string]interface{}{
			"_dst": addr,
		})
	}
}

func copyData(dst io.Writer, src io.Reader, ch chan<- error) {
	_, err := io.Copy(dst, src)
	if tdst, ok := dst.(*net.TCPConn); ok {
		tdst.CloseWrite()
	}
	if tsrc, ok := src.(*net.TCPConn); ok {
		tsrc.CloseRead()
	}
	ch <- err
}
