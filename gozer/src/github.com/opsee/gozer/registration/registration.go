package registration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/opsee/portmapper"
)

const (
	ipFilePath = "/zuul/state/ip"
	nsqdTopic  = "connected"
)

// /opsee.co/routes/customer_id/instance_id/svcname=ip:port

type connectedMessage struct {
	CustomerID string              `json:"customer_id"`
	BastionID  string              `json:"bastion_id"`
	InstanceID string              `json:"instance_id"`
	IPAddress  string              `json:"ip_address"`
	Services   []*portmapper.Service `json:"services"`
	Timestamp  int64               `json:"timestamp"`
}

// Service provides registration with Opsee of components exposed by
// portmapper.
type Service interface {
	Start() error
	Stop() error
}

type nsqdService struct {
	producer             *nsq.Producer
	nsqdHost             string
	portmapPath          string
	stopChan             chan struct{}
	registrationInterval time.Duration
	customerID           string
	bastionID            string
	instanceID           string
}

// NewService creates a new registration.Service for NSQ using portmapper.
func NewService(interval time.Duration, nsqdHost string, customerID string, bastionID string, instanceID string) *nsqdService {
	svc := &nsqdService{
		nsqdHost:             nsqdHost,
		portmapPath:          portmapper.RegistryPath,
		stopChan:             make(chan struct{}),
		registrationInterval: interval,
		customerID:           customerID,
		bastionID:            bastionID,
		instanceID:           instanceID,
	}

	return svc
}

func (s *nsqdService) register() {
	svcs, err := portmapper.Services()
	if err != nil {
		log.Println(err.Error())
		return
	}

	ip, err := ioutil.ReadFile(ipFilePath)
	if err != nil {
		log.Println("Error reading IP from file:", ipFilePath)
		log.Println(err.Error())
		return
	}

	if len(ip) == 0 {
		log.Println("IP file empty:", ipFilePath)
		return
	}

	msg := &connectedMessage{
		s.customerID,
		s.bastionID,
		s.instanceID,
		string(ip),
		svcs,
		time.Now().Unix(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Unable to marshal message:", msg)
		log.Println(err.Error())
		return
	}

	log.Println("Publishing message:", msg)
	if err := s.producer.Publish(nsqdTopic, msgBytes); err != nil {
		log.Println(err.Error())
	}
}

func (s *nsqdService) registrationLoop() {
	ticker := time.NewTicker(s.registrationInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.register()
		case <-s.stopChan:
			s.producer.Stop()
			return
		}
	}
}

// Start the registration loop.
func (s *nsqdService) Start() error {
	producer, err := nsq.NewProducer(s.nsqdHost, nsq.NewConfig())
	if err != nil {
		return err
	}
	s.producer = producer

	go s.registrationLoop()
	return nil
}

// Stop the registration loop.
func (s *nsqdService) Stop() error {
	s.stopChan <- struct{}{}

	return nil
}
