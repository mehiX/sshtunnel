package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	Local       *Endpoint
	Server      *Endpoint
	Destination *Endpoint
	Config      *ssh.ClientConfig
	Log         *log.Logger
}

func NewSSHTunnel(tunnelSrvr string, auth ssh.AuthMethod, destSrvr string) *SSHTunnel {

	local := NewEndpoint("localhost:0")
	srvr := NewEndpoint(tunnelSrvr)
	if srvr.Port == 0 {
		srvr.Port = 22
	}

	cfg := &ssh.ClientConfig{
		User: srvr.User,
		Auth: []ssh.AuthMethod{auth},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return &SSHTunnel{
		Local:       local,
		Server:      srvr,
		Destination: NewEndpoint(destSrvr),
		Config:      cfg,
	}

}

func (t *SSHTunnel) Logf(fmt string, args ...interface{}) {
	if t.Log != nil {
		t.Log.Printf(fmt, args...)
	}
}

func (t *SSHTunnel) Start() error {

	// start the local listener
	localListener, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}

	defer localListener.Close()

	localPort := localListener.Addr().(*net.TCPAddr).Port

	fmt.Printf("Listening on port %d\n", localPort)

	for {
		c, err := localListener.Accept()
		if err != nil {
			return err
		}

		go t.fwd(c)
	}

}

func (t *SSHTunnel) fwd(c net.Conn) {

	// connect and authenticate to the SSH server
	client, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		t.Logf("server dial error: %s\n", err)
		return
	}

	defer client.Close()

	t.Logf("connected to %s (1 of 2)\n", t.Server.String())

	// connect to destination
	destConn, err := client.Dial("tcp", t.Destination.String())
	if err != nil {
		t.Logf("remote dial error: %s\n", err)
		return
	}

	defer destConn.Close()

	t.Logf("connected to %s (2 of 2)\n", t.Destination.String())

	copyConn := func(writer, reader net.Conn) {
		if _, err := io.Copy(writer, reader); err != nil {
			t.Logf("io.Copy error: %s\n", err)
		}
	}

	go copyConn(c, destConn)
	copyConn(destConn, c)
}
