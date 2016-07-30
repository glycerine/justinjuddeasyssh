package easyssh_test

import (
	"fmt"
	"io/ioutil"

	"dev.justinjudd.org/justin/easyssh"

	"golang.org/x/crypto/ssh"
)

func ExampleServer_ListenAndServe() {
	s := easyssh.Server{Addr: ":2022"}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		// Failed to load private key (./id_rsa)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		// Failed to parse private key
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == "test" && string(pass) == "test" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %s", c.User())
		},
	}
	config.AddHostKey(private)

	s.Config = config

	handler := easyssh.NewStandardSSHServerHandler()
	channelHandler := easyssh.NewChannelsMux()

	channelHandler.HandleChannel(easyssh.SessionRequest, easyssh.SessionHandler())
	handler.MultipleChannelsHandler = channelHandler

	s.Handler = handler

	s.ListenAndServe()
}

func ExampleListenAndServe() {
	config := &ssh.ServerConfig{}
	easyssh.HandleChannel(easyssh.SessionRequest, easyssh.SessionHandler())
	easyssh.HandleChannel(easyssh.DirectForwardRequest, easyssh.DirectPortForwardHandler())
	easyssh.HandleRequestFunc(easyssh.RemoteForwardRequest, easyssh.TCPIPForwardRequest)

	easyssh.ListenAndServe(":2022", config, nil)
}

func ExampleChannelsMux_HandleChannelFunc() {
	handler := easyssh.NewChannelsMux()

	testHandler := func(newChannel ssh.NewChannel, channel ssh.Channel, reqs <-chan *ssh.Request, sshConn ssh.Conn) {
		defer channel.Close()
		ssh.DiscardRequests(reqs)
	}

	handler.HandleChannelFunc("test", testHandler)

}
