// Package rtmp implements RTMP protocol.
package rtmp

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/protocol"
)

// Handler represents handler for `flv` structure.
type Handler struct {
	rawURL string

	url        *url.URL
	serverHost string   // server host or ip address, without port
	serverPort uint     // server port, 1935 by default
	conn       net.Conn // conection after dial
}

// NewHandler creates RTMP Handler.
func NewHandler(url string) (*Handler, error) {
	h := Handler{rawURL: url}
	if err := h.parseURL(); err != nil {
		return nil, err
	}

	return &h, nil
}

// Connect connects rtmp server.
func (h *Handler) Connect(timeout time.Duration) error {
	startTime := time.Now()
	serverAddress := net.JoinHostPort(h.serverHost, strconv.FormatInt(int64(h.serverPort), 10))
	glog.Infof("connect %s", serverAddress)
	conn, err := net.DialTimeout("tcp", serverAddress, timeout)
	if err != nil {
		return err
	}
	h.conn = conn

	// timeout for the connecting
	if err = h.conn.SetWriteDeadline(startTime.Add(timeout)); err != nil {
		return err
	}

	// handshark
	if err = h.handshark(); err != nil {
		return err
	}

	glog.Infof("connected %s, takes %f seconds", serverAddress, time.Since(startTime).Seconds())
	return nil
}

// Close closes the handler.
func (h *Handler) Close() {
	if h.conn != nil {
		if err := h.conn.Close(); err != nil {
			glog.Warning(err)
		}
		h.conn = nil
	}
}

// parseURL
func (h *Handler) parseURL() error {
	if len(h.rawURL) == 0 {
		return fmt.Errorf("url is empty")
	}

	var err error
	if h.url, err = url.Parse(h.rawURL); err != nil {
		return err
	}

	if h.url == nil {
		return fmt.Errorf("empty URL after parsed %s", h.rawURL)
	}
	if h.url.Scheme != protocol.SchemaRTMP {
		return fmt.Errorf("invalid schema %s", h.url.Scheme)
	}
	hostPort := strings.Split(h.url.Host, ":")
	if len(hostPort) == 0 {
		return fmt.Errorf("invalid url %s", h.rawURL)
	}
	h.serverHost = hostPort[0]
	if len(hostPort) >= 2 {
		if port, err := strconv.ParseUint(hostPort[1], 10, 16); err != nil {
			return err
		} else {
			h.serverPort = uint(port)
		}
	} else {
		h.serverPort = protocol.PortRTMP // default port for rtmp
	}

	return nil
}
