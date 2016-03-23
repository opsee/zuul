package registration

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/nsqio/go-nsq"
	"github.com/opsee/portmapper"
	stdlog "log"
	"os"
)

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

const (
	routesPath = "/opsee.co/routes"
	// RegistrationTTL defines the number of seconds that a registration will be
	// valid.
	RegistrationTTL   = 180
	RequestTimeoutSec = 30
)

type consumerService struct {
	etcdClient   client.Client
	consumer     *nsq.Consumer
	stopChan     chan struct{}
	lookupdHosts []string
	maxRetries   int
}

// NewConsumer creates a new consumer service connected to the "connected" topic
// in NSQ.
func NewConsumer(consumerName, etcdHost string, nsqLookupdHosts []string, concurrency int, maxRetries int) (*consumerService, error) {
	consumer, err := nsq.NewConsumer("_.connected", consumerName, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	consumer.SetLogger(stdlog.New(os.Stdout, "[nsqd] ", stdlog.LstdFlags), nsq.LogLevelDebug)

	cfg := client.Config{
		Endpoints: []string{etcdHost},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	// initialize a new etcd client
	etcdClient, err := client.New(cfg)
	if err != nil {
		panic(err)
	}

	svc := &consumerService{
		etcdClient,
		consumer,
		make(chan struct{}, 1),
		nsqLookupdHosts,
		maxRetries,
	}

	svc.consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) error {
		cMsg := &ConnectedMessage{}
		if err := json.Unmarshal(message.Body, cMsg); err != nil {
			log.WithError(err).Error("error unmarshaling connected message: ", cMsg)
			return nil
		}

		logger := log.WithFields(log.Fields{"customer_id": cMsg.CustomerID, "bastion_id": cMsg.BastionID})
		logger.Info("handling message")

		svcMap := make(map[string]*portmapper.Service)

		for _, svc := range cMsg.Services {
			svc.Hostname = cMsg.IPAddress
			svcMap[svc.Name] = svc
		}

		key := fmt.Sprintf("/opsee.co/routes/%s/%s", cMsg.CustomerID, cMsg.BastionID)
		mapBytes, err := json.Marshal(svcMap)
		if err != nil {
			logger.WithError(err).Error("error marshaling service map")
			return nil
		}

		kAPI := client.NewKeysAPI(etcdClient)
		ctx, cancel := context.WithTimeout(context.Background(), RequestTimeoutSec*time.Second)
		defer cancel()

		_, err = kAPI.Set(ctx, key, string(mapBytes), &client.SetOptions{TTL: RegistrationTTL * time.Second})
		if err != nil {
			logger.WithError(err).Error("couldn't register with etcd")
			return nil
		}

		logger.Info("successfully registered service with etcd")
		message.Finish()
		return nil
	}), concurrency)

	return svc, nil
}

func (c *consumerService) Start() error {
	go func() {
		for {
			stats := c.consumer.Stats()
			isStarved := c.consumer.IsStarved()
			log.Infof("(NSQ) Received:%d, Finished:%d, Requeued:%d, Connections: %d, Starved: %t", stats.MessagesReceived, stats.MessagesFinished, stats.MessagesRequeued, stats.Connections, isStarved)
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()

	return c.consumer.ConnectToNSQLookupds(c.lookupdHosts)
}

func (c *consumerService) Stop() error {
	c.consumer.Stop()
	return nil
}
