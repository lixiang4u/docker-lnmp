package util

import (
	"github.com/docker/docker/client"
)

var _dockerClientInstance *client.Client

func ConnectDocker() (*client.Client, error) {
	if _dockerClientInstance != nil {
		return _dockerClientInstance, nil
	}
	// Error response from daemon: client version 1.42 is too new. Maximum supported API version is 1.41
	return client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
}
