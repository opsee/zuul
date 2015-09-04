package registration

import (
	"github.com/opsee/portmapper"
)

const (
	ipFilePath = "/zuul/state/ip"
	nsqdTopic  = "connected"
)

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

type connectedMessage struct {
	CustomerID string                `json:"customer_id"`
	BastionID  string                `json:"bastion_id"`
	InstanceID string                `json:"instance_id"`
	IPAddress  string                `json:"ip_address"`
	Services   []*portmapper.Service `json:"services"`
	Timestamp  int64                 `json:"timestamp"`
}

// Service provides registration with Opsee of components exposed by
// portmapper.
type Service interface {
	Start() error
	Stop() error
}