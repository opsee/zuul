package registration

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/nsqio/go-nsq"
	"github.com/opsee/portmapper"
)

// /opsee.co/routes/customer_id/instance_id/svcname = ip:port

const (
	routesPath = "/opsee.co/routes"
	// RegistrationTTL defines the number of seconds that a registration will be
	// valid.
	RegistrationTTL   = 150
	RequestTimeoutSec = 120
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

	cfg := client.Config{
		Endpoints: []string{etcdHost},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	// initialize a new etcd client
	c, err := client.New(cfg)
	if err != nil {
		panic(err)
	}

	svc := &consumerService{
		c,
		consumer,
		make(chan struct{}, 1),
		nsqLookupdHosts,
		maxRetries,
	}

	svc.consumer.AddConcurrentHandlers(nsq.HandlerFunc(svc.registerConnection), concurrency)

	return svc, nil
}

// /opsee.co/routes/customer_id/instance_id
func (c *consumerService) registerConnection(msg *nsq.Message) error {
	cMsg := &ConnectedMessage{}
	if err := json.Unmarshal(msg.Body, cMsg); err != nil {
		logrus.WithFields(logrus.Fields{"module": "register", "event": "registerConnection", "Error": err}).Error("Error unmarshaling connected message: ", msg)
		return err
	}

	logrus.WithFields(logrus.Fields{"module": "register", "event": "registerConnection"}).Info("Handling message: ", string(msg.Body))
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

	kAPI := client.NewKeysAPI(c.etcdClient)
	for try := 0; try < c.maxRetries; try++ {
		ctx, cancel := context.WithTimeout(context.Background(), RequestTimeoutSec*time.Second)
		defer cancel()

		_, err := kAPI.Set(ctx, key, string(mapBytes), &client.SetOptions{TTL: RegistrationTTL})
		if err != nil {
			// handle error
			if err == context.DeadlineExceeded {
				logrus.WithFields(logrus.Fields{
					"module":  "register",
					"event":   "registerConnection",
					"key":     key,
					"val":     string(mapBytes),
					"attempt": try,
					"errstr":  err.Error(),
				}).Warn("Service path set exceeded context deadline. Retrying")
			} else {
				logrus.WithFields(logrus.Fields{
					"module":  "register",
					"event":   "registerConnection",
					"key":     key,
					"val":     string(mapBytes),
					"attempt": try,
					"errstr":  err.Error(),
				}).Error("Service path set failed.")
				return err
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"module":  "register",
				"event":   "registerConnection",
				"key":     key,
				"val":     string(mapBytes),
				"attempt": try,
				"errstr":  err.Error(),
			}).Info("Successfully registered service with etcd")
			break
		}

		time.Sleep(2 << uint(try) * time.Millisecond)
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
