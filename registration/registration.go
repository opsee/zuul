package registration

import (
	"github.com/opsee/portmapper"
)

const (
	ipFilePathDefault = "/gozer/state/ip"
	nsqdTopic         = "_.connected"
)

var (
	// The location of the file produced by OpenVPN containing the instance IP.
	IPFilePath string
)

func init() {
	IPFilePath = ipFilePathDefault
}

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

type ConnectedMessage struct {
	CustomerID string                `json:"customer_id"`
	BastionID  string                `json:"bastion_id"`
	InstanceID string                `json:"instance_id"`
	IPAddress  string                `json:"ip_address"`
	PublicIP   string                `json:"public_ip"`
	Services   []*portmapper.Service `json:"services"`
	Timestamp  int64                 `json:"timestamp"`
}

// Service provides registration with Opsee of components exposed by
// portmapper.
type Service interface {
	Start() error
	Stop() error
}
