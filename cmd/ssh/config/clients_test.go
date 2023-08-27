package config

import (
	"github.com/AWoelfel/go-multi-ssh/tests/assert"
	"github.com/kevinburke/ssh_config"

	"testing"
)

func TestListClients(t *testing.T) {

	t.Run("duplicates", func(t *testing.T) {
		hosts, err := listClients("testdata/ssh_config_duplicates")
		assert.NoError(t, err)
		assert.EqualArrayValues(t, []string{"unifi-manager.lan", "dns.lan"}, hosts)
	})

	t.Run("mixed", func(t *testing.T) {
		hosts, err := listClients("testdata/ssh_config_mixed")
		assert.NoError(t, err)
		assert.EqualArrayValues(t, []string{"fence.lan", "unifi-manager.lan", "dns.lan", "dns.iot"}, hosts)
	})

	t.Run("single", func(t *testing.T) {
		hosts, err := listClients("testdata/ssh_config_single")
		assert.NoError(t, err)
		assert.EqualArrayValues(t, []string{"fence.lan"}, hosts)
	})

	t.Run("comments", func(t *testing.T) {
		hosts, err := listClients("testdata/ssh_config_with_comments")
		assert.NoError(t, err)
		assert.EqualArrayValues(t, []string{"fence.lan"}, hosts)
	})

}

func TestBuildClientConfigs(t *testing.T) {

	sshClientConfig := &ssh_config.UserSettings{}
	sshClientConfig.WithConfigLocations(func() (string, error) { return "testdata/ssh_config", nil }, func() (string, error) { return "testdata/user_ssh_config_mock", nil })

	hosts, err := buildClientConfigs(sshClientConfig, []string{"fence.lan", "unifi-manager.lan", "dns.lan", "dns.iot"})
	assert.NoError(t, err)
	assert.ArrayLen(t, 4, hosts)

	expectedHostsConfig := []SshClientConfig{
		{
			Port:         22,
			Host:         "fence.lan",
			IdentityFile: "/root/.ssh/fence.lan_id_rsa",
			User:         "fence.user",
			Password:     "",
			Tags:         []string{"target"},
			Color:        1,
		},
		{
			Port:         4422,
			Host:         "unifi-manager.lan",
			IdentityFile: "/root/.ssh/unifi-manager.lan_id_rsa",
			User:         "root",
			Password:     "",
			Tags:         []string{"ubuntu", "vm"},
			Color:        2,
		},
		{
			Port:         4423,
			Host:         "dns.lan",
			IdentityFile: "/root/.ssh/dns.lan_id_rsa",
			User:         "root",
			Password:     "",
			Tags:         []string{"target", "dns"},
			Color:        5,
		},
		{
			Port:         22,
			Host:         "dns.iot",
			IdentityFile: "/root/.ssh/dns.iot_id_rsa",
			User:         "foo",
			Password:     "secure-login-passowrd",
			Tags:         []string{"dns"},
			Color:        6,
		},
	}
	assert.AssertObjectsEqual(t, expectedHostsConfig, hosts)
}
