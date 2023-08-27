package connection

import (
	"fmt"
	"github.com/AWoelfel/go-multi-ssh/cmd/ssh/config"
	"github.com/muesli/termenv"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"sync"
)

type ClientContext struct {
	sync.Mutex
	ClientConfig        config.SshClientConfig
	sshConnectionConfig *ssh.ClientConfig
}

func NewClientContext(clientConfig config.SshClientConfig) *ClientContext {
	return &ClientContext{ClientConfig: clientConfig}
}

func (cc *ClientContext) fillSshConnectionConfig() error {

	if cc.sshConnectionConfig != nil {
		return nil
	}

	cc.sshConnectionConfig = &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            cc.ClientConfig.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		BannerCallback:  nil,
	}

	if cc.ClientConfig.IdentityFile != "" {
		key, err := os.ReadFile(cc.ClientConfig.IdentityFile)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("unable to read private key file %s (%w)", cc.ClientConfig.IdentityFile, err)
		}

		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				return fmt.Errorf("unable to create signed keypair from private key file %s (%w)", cc.ClientConfig.IdentityFile, err)
			}

			cc.sshConnectionConfig.Auth = append(cc.sshConnectionConfig.Auth, ssh.PublicKeys(signer))
		}
	}

	if cc.ClientConfig.Password != "" {
		cc.sshConnectionConfig.Auth = append(cc.sshConnectionConfig.Auth, ssh.Password(cc.ClientConfig.Password))
	}

	return nil
}

func (cc *ClientContext) Label() string {
	return cc.ClientConfig.Host
}

func (cc *ClientContext) Color() termenv.ANSI256Color {
	return cc.ClientConfig.Color
}
