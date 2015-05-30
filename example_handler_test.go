package easyssh_test

import (
	"dev.justinjudd.org/justin/easyssh"

	"golang.org/x/crypto/ssh"
)

type testHandler struct{}

func (testHandler) HandleChannel(nCh ssh.NewChannel, ch ssh.Channel, reqs <-chan *ssh.Request, conn *ssh.ServerConn) {
	defer ch.Close()
	// Do something
}

func ExampleChannelsMux_HandleChannel() {
	handler := easyssh.NewChannelsMux()

	handler.HandleChannel("test", testHandler{})

	test2Handler := func(newChannel ssh.NewChannel, channel ssh.Channel, reqs <-chan *ssh.Request, sshConn *ssh.ServerConn) {
		defer channel.Close()
		ssh.DiscardRequests(reqs)
	}

	handler.HandleChannelFunc("test2", test2Handler)

	handler.HandleChannel("anotherTest2", easyssh.ChannelHandlerFunc(test2Handler))
}
