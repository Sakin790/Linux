package handlers

import (
	"ece2clone/types"
	"ece2clone/utils"
	"encoding/json"
	"net/http"

	"github.com/lxc/incus/v6/shared/api"
)

func ListInstancesHandler(w http.ResponseWriter, r *http.Request) {
	server, err := utils.ConnectIncus()
	if err != nil {
		http.Error(w, "failed to connect to incus daemon: "+err.Error(), http.StatusInternalServerError)
		return
	}

	instances, err := server.GetInstancesFull(api.InstanceTypeAny)
	if err != nil {
		http.Error(w, "failed to list instances: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]types.InstanceSummary, 0, len(instances))
	for _, inst := range instances {
		var ips []string
		if inst.State != nil {
			for _, net := range inst.State.Network {
				for _, addr := range net.Addresses {
					if addr.Family == "inet" && addr.Scope == "global" {
						ips = append(ips, addr.Address)
					}
				}
			}
		}

		result = append(result, types.InstanceSummary{
			Name:      inst.Name,
			Status:    inst.Status,
			Type:      inst.Type,
			IPv4:      ips,
			Profiles:  inst.Profiles,
			CreatedAt: inst.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
