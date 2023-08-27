package connection

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func OpenClient(clientCtx *ClientContext) (*ssh.Client, error) {
	clientCtx.Lock()
	defer clientCtx.Unlock()

	err := clientCtx.fillSshConnectionConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to prepare ssh connection config (%w)", err)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", clientCtx.ClientConfig.Host, clientCtx.ClientConfig.Port), clientCtx.sshConnectionConfig)
}
