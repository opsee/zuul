package registration

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/nsqio/go-nsq"
)

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

const (
	routesPath = "/opsee.co/routes"
)

type consumerService struct {
	etcdClient *etcd.Client
	consumer   *nsq.Consumer
	stopChan   chan struct{}
}

// NewConsumer creates a new consumer service connected to the "connected" topic
// in NSQ.
func NewConsumer(consumerName, etcdHost, nsqLookupdHost string) (*consumerService, error) {
	consumer, err := nsq.NewConsumer("connected", consumerName, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	svc := &consumerService{
		etcd.NewClient([]string{etcdHost}),
		consumer,
		make(chan struct{}, 1),
	}

	return svc, nil
}

func (c *consumerService) Start() error {
	return nil
}

func (c *consumerService) Stop() error {
	return nil
}
