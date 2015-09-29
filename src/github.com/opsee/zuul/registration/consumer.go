package registration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/nsqio/go-nsq"
)

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

const (
	routesPath = "/opsee.co/routes"
	// RegistrationTTL defines the number of seconds that a registration will be
	// valid.
	RegistrationTTL = 60
)

var (
	nsqLookupds = []string{
		"nsqlookupd-1.opsy.co",
		"nsqlookupd-2.opsy.co",
	}
)

type consumerService struct {
	etcdClient   *etcd.Client
	consumer     *nsq.Consumer
	stopChan     chan struct{}
	lookupdHosts []string
}

// NewConsumer creates a new consumer service connected to the "connected" topic
// in NSQ.
func NewConsumer(consumerName, etcdHost string, nsqLookupdHosts []string, concurrency int) (*consumerService, error) {
	consumer, err := nsq.NewConsumer("_.connected", consumerName, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	consumer.SetLogger(log.New(os.Stderr, "", log.LstdFlags), nsq.LogLevelInfo)

	svc := &consumerService{
		etcd.NewClient([]string{etcdHost}),
		consumer,
		make(chan struct{}, 1),
		nsqLookupdHosts,
	}

	svc.consumer.AddConcurrentHandlers(nsq.HandlerFunc(svc.registerConnection), concurrency)

	return svc, nil
}

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port
func (c *consumerService) registerConnection(msg *nsq.Message) error {
	cMsg := &ConnectedMessage{}
	if err := json.Unmarshal(msg.Body, cMsg); err != nil {
		log.Println("Error unmarshaling connected message:", msg)
		return err
	}
	log.Println("Handling message:", msg.Body)

	for _, connectedSvc := range cMsg.Services {
		log.Printf("Registering %s service for customer %s, bastion %s, at IP: %s, port: %s", connectedSvc.Name, cMsg.CustomerID, cMsg.InstanceID, cMsg.IPAddress, connectedSvc.Port)
		key := fmt.Sprintf("/opsee.co/routes/%s/%s/%s", cMsg.CustomerID, cMsg.InstanceID, connectedSvc.Name)
		value := fmt.Sprintf("%s:%d", cMsg.IPAddress, connectedSvc.Port)

		resp, err := c.etcdClient.Set(key, value, 60)
		if err != nil {
			return err
		}
		log.Printf("ETCD Response Node: %s", *resp.Node)
	}

	return nil
}

func (c *consumerService) Start() error {
	return c.consumer.ConnectToNSQLookupds(c.lookupdHosts)
}

func (c *consumerService) Stop() error {
	c.consumer.Stop()
	return nil
}