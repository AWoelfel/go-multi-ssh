package config

import (
	"bytes"
	"fmt"
	"github.com/kevinburke/ssh_config"
	"github.com/muesli/termenv"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
)

type SshClientConfig struct {
	Port         int
	Host         string
	IdentityFile string
	User         string
	Password     string
	Tags         []string
	Color        termenv.ANSI256Color
}

const _USER_HOME_SHORTCUT = "~"
const _CONST_PERCENTAGE = "%%"
const _USER_HOME = "%d"
const _REMOTE_HOST = "%h"
const _USER_ID = "%i"
const _LOCAL_HOST = "%L"
const _COMMAND_LINE_HOSTNAME = "%n"
const _REMOTE_PORT = "%p"
const _REMOTE_USER = "%r"
const _LOCAL_USER = "%u"

type sshClientStrategy func(targetHost string, target *SshClientConfig, setting string) (string, error)

var sshClientStrategies = make(map[string]sshClientStrategy)
var sshClientStrategiesSetup = sync.Once{}
var sshClientStrategiesOrder = []string{_USER_HOME_SHORTCUT, _LOCAL_USER, _REMOTE_USER, _REMOTE_PORT, _COMMAND_LINE_HOSTNAME, _LOCAL_HOST, _USER_ID, _REMOTE_HOST, _USER_HOME, _CONST_PERCENTAGE}

func init() {

	sshClientStrategiesSetup.Do(func() {
		sshClientStrategies[_USER_HOME_SHORTCUT] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			return strings.Replace(setting, _USER_HOME_SHORTCUT, _USER_HOME, -1), nil
		}
		sshClientStrategies[_CONST_PERCENTAGE] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			return strings.Replace(setting, _CONST_PERCENTAGE, "%", -1), nil
		}
		sshClientStrategies[_USER_HOME] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			userHome, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			return strings.Replace(setting, _USER_HOME, userHome, -1), nil
		}
		sshClientStrategies[_REMOTE_HOST] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			return strings.Replace(setting, _REMOTE_HOST, target.Host, -1), nil
		}
		sshClientStrategies[_USER_ID] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			user, err := user.Current()
			if err != nil {
				return "", err
			}
			return strings.Replace(setting, _USER_ID, user.Uid, -1), nil
		}
		sshClientStrategies[_LOCAL_HOST] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			hostname, err := os.Hostname()
			if err != nil {
				return "", err
			}
			return strings.Replace(setting, _LOCAL_HOST, hostname, -1), nil
		}
		sshClientStrategies[_COMMAND_LINE_HOSTNAME] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			return strings.Replace(setting, _COMMAND_LINE_HOSTNAME, targetHost, -1), nil
		}
		sshClientStrategies[_REMOTE_PORT] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			return strings.Replace(setting, _REMOTE_PORT, strconv.Itoa(target.Port), -1), nil
		}
		sshClientStrategies[_REMOTE_USER] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			return strings.Replace(setting, _REMOTE_USER, target.User, -1), nil
		}
		sshClientStrategies[_LOCAL_USER] = func(targetHost string, target *SshClientConfig, setting string) (string, error) {
			user, err := user.Current()
			if err != nil {
				return "", err
			}
			return strings.Replace(setting, _USER_ID, user.Username, -1), nil
		}

	})

}

func fillSshClientConfig(targetHost string, target *SshClientConfig) error {

	var err error
	for i := 0; i < len(sshClientStrategiesOrder); i++ {
		target.IdentityFile, err = sshClientStrategies[sshClientStrategiesOrder[i]](targetHost, target, target.IdentityFile)
		if err != nil {
			return fmt.Errorf("unable to execute strat %q (%w)", sshClientStrategiesOrder[i], err)
		}
	}

	return nil
}

func listClients(sshConfigFile string) ([]string, error) {

	data, err := os.ReadFile(sshConfigFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read ssh config file %s (%w)", sshConfigFile, err)
	}

	cfg, err := ssh_config.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("unable to decode ssh config file %s (%w)", sshConfigFile, err)
	}

	var result []string
	allKeys := make(map[string]bool)

	for hIdx, configHost := range cfg.Hosts {

		for i := 0; i < len(configHost.Patterns); i++ {
			//check whether this pattern is any wildcard pattern...

			p := configHost.Patterns[i].String()

			//TODO: this is faulty; Patterns do not change as "pattern.String()" returns the original pattern not the Regex Pattern...
			if p == "*" {
				continue
			}

			compiledPattern, err := ssh_config.NewPattern(p)
			if err != nil {
				return nil, fmt.Errorf("unable to compile pattern %d (%s) of host %d (%w)", i, p, hIdx, err)
			}
			if compiledPattern.String() == p {
				//no wildcards

				if _, found := allKeys[p]; !found {
					allKeys[p] = true
					result = append(result, p)
				}
			}
		}
	}

	return result, nil
}

func buildClientConfigs(sshClientConfig *ssh_config.UserSettings, targetHosts []string) ([]SshClientConfig, error) {

	result := make([]SshClientConfig, len(targetHosts), len(targetHosts))

	for i, targetHost := range targetHosts {
		userValue := sshClientConfig.Get(targetHost, "User")
		identityFileValue := sshClientConfig.Get(targetHost, "IdentityFile")
		passwordValue := sshClientConfig.Get(targetHost, "Password")
		tags := sshClientConfig.GetAll(targetHost, "Tag")

		portValue := sshClientConfig.Get(targetHost, "Port")
		portIntValue, err := strconv.ParseInt(portValue, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %q as int  for use as 'Port' (%w)", portValue, err)
		}

		colorValue := sshClientConfig.Get(targetHost, "Color")
		var colorIntValue int64
		if len(colorValue) > 0 {
			colorIntValue, err = strconv.ParseInt(colorValue, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("unable to parse %q as int for use as 'Color' (%w)", portValue, err)
			}
		}
		if colorIntValue == 0 {
			colorIntValue = int64(rand.Intn(256))
		}

		resultConfig := SshClientConfig{
			Port:         int(portIntValue),
			Host:         targetHost,
			IdentityFile: identityFileValue,
			User:         userValue,
			Password:     passwordValue,
			Tags:         tags,
			Color:        termenv.ANSI256Color(int(colorIntValue)),
		}

		err = fillSshClientConfig(targetHost, &resultConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to prefill ssh client config (%w)", err)
		}

		result[i] = resultConfig
	}

	return result, nil
}

func ReadClients(sshConfigFile string) ([]SshClientConfig, error) {

	targets, err := listClients(sshConfigFile)
	if err != nil {
		return nil, fmt.Errorf("unable to list target hosts (%w)", err)
	}

	sshClientConfig := &ssh_config.UserSettings{}

	sshClientConfig.
		WithConfigLocations(func() (string, error) { return sshConfigFile, nil }).
		AddConfigLocations(ssh_config.DefaultConfigFileFinders...)

	return buildClientConfigs(sshClientConfig, targets)
}
