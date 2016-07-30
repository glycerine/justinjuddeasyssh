package easyssh

import (
	"net"

	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH Client
type Client struct {
	*ssh.Client
}

// Dial starts an ssh connection to the provided server
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	c, err := ssh.Dial(network, addr, config)
	return &Client{c}, err
}

// NewClient returns a new SSH Client.
func NewClient(c ssh.Conn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request) *Client {

	client := ssh.NewClient(c, chans, reqs)
	return &Client{client}

}

// LocalForward performs a port forwarding over the ssh connection - ssh -L. Client will bind to the local address, and will tunnel those requests to host addr
func (c *Client) LocalForward(laddr, raddr *net.TCPAddr) error {

	ln, err := net.ListenTCP("tcp", laddr) //tie to the client connection
	if err != nil {
		println(err.Error())
		return err
	}
	logger.Println("Listening on address: ", ln.Addr().String())

	quit := make(chan bool)

	go func() { // Handle incoming connections on this new listener
		for {
			select {
			case <-quit:

				return
			default:
				conn, err := ln.Accept()
				if err != nil { // Unable to accept new connection - listener likely closed
					continue
				}
				go func(conn net.Conn) {
					conn2, err := c.DialTCP("tcp", laddr, raddr)

					if err != nil {
						return
					}
					go func(conn, conn2 net.Conn) {

						close := func() {
							conn.Close()
							conn2.Close()

						}

						go CopyReadWriters(conn, conn2, close)

					}(conn, conn2)

				}(conn)
			}

		}
	}()

	c.Wait()

	ln.Close()
	quit <- true

	return nil
}

// RemoteForward forwards a remote port - ssh -R
func (c *Client) RemoteForward(remote, local string) error {
	ln, err := c.Listen("tcp", remote)
	if err != nil {
		return err
	}

	quit := make(chan bool)

	go func() { // Handle incoming connections on this new listener
		for {
			select {
			case <-quit:

				return
			default:
				conn, err := ln.Accept()
				if err != nil { // Unable to accept new connection - listener likely closed
					continue
				}

				conn2, err := net.Dial("tcp", local)
				if err != nil {
					continue
				}

				close := func() {
					conn.Close()
					conn2.Close()

				}

				go CopyReadWriters(conn, conn2, close)

			}

		}
	}()

	c.Wait()
	ln.Close()
	quit <- true

	return nil
}

// HandleOpenChannel requests that the remote end accept a channel request and if accepted,
// passes the newly opened channel and requests to the provided handler
func (c *Client) HandleOpenChannel(channelName string, handler ChannelMultipleRequestsHandler, data ...byte) error {
	ch, reqs, err := c.OpenChannel(channelName, data)
	if err != nil {
		return err
	}
	handler.HandleMultipleRequests(reqs, c.Conn, channelName, ch)
	return nil
}

// HandleOpenChannelFunc requests that the remote end accept a channel request and if accepted,
// passes the newly opened channel and requests to the provided handler function
func (c *Client) HandleOpenChannelFunc(channelName string, handler ChannelMultipleRequestsHandlerFunc, data ...byte) error {

	return c.HandleOpenChannel(channelName, ChannelMultipleRequestsHandlerFunc(handler), data...)
}
