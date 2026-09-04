package utils

import incus "github.com/lxc/incus/v6/client"

func ConnectIncus() (incus.InstanceServer, error) {
	return incus.ConnectIncusUnix("", nil)
}
