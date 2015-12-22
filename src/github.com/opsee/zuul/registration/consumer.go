package registration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/nsqio/go-nsq"
	"github.com/opsee/portmapper"
)

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

const (
	routesPath = "/opsee.co/routes"
	// RegistrationTTL defines the number of seconds that a registration will be
	// valid.
	RegistrationTTL = 150
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

// /opsee.co/routes/customer_id/instance_id
func (c *consumerService) registerConnection(msg *nsq.Message) error {
	cMsg := &ConnectedMessage{}
	if err := json.Unmarshal(msg.Body, cMsg); err != nil {
		log.Println("Error unmarshaling connected message:", msg)
		return err
	}
	log.Println("Handling message:", string(msg.Body))
	svcMap := make(map[string]*portmapper.Service)

	for _, svc := range cMsg.Services {
		svc.Hostname = cMsg.IPAddress
		svcMap[svc.Name] = svc
	}

	key := fmt.Sprintf("/opsee.co/routes/%s/%s", cMsg.CustomerID, cMsg.BastionID)
	mapBytes, err := json.Marshal(svcMap)
	if err != nil {
		return err
	}

	if _, err := c.etcdClient.Set(key, string(mapBytes), RegistrationTTL); err != nil {
		return err
	}

	return nil
}

func (c *consumerService) Start() error {
	return c.consumer.ConnectToNSQLookupds(c.lookupdHosts)
}

func (c *consumerService) Stop() error {
	c.consumer.Stop()
	c.etcdClient.Close()
	return nil
}
