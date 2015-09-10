package registration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/opsee/portmapper"
)

type producerService struct {
	producer             *nsq.Producer
	nsqdHost             string
	portmapPath          string
	stopChan             chan struct{}
	registrationInterval time.Duration
	customerID           string
	bastionID            string
	instanceID           string
}

// NewProducer creates a new registration.Service for NSQ using portmapper.
func NewProducer(interval time.Duration, nsqdHost string, customerID string, bastionID string, instanceID string) *producerService {
	svc := &producerService{
		nsqdHost:             nsqdHost,
		portmapPath:          portmapper.RegistryPath,
		stopChan:             make(chan struct{}, 1),
		registrationInterval: interval,
		customerID:           customerID,
		bastionID:            bastionID,
		instanceID:           instanceID,
	}

	portmapper.EtcdHost = os.Getenv("ETCD_HOST")

	return svc
}

func (s *producerService) register() {
	svcs, err := portmapper.Services()
	if err != nil {
		log.Println(err.Error())
		return
	}

	ip, err := ioutil.ReadFile(IPFilePath)
	if err != nil {
		log.Println("Error reading IP from file:", IPFilePath)
		log.Println(err.Error())
		return
	}

	if len(ip) == 0 {
		log.Println("IP file empty:", IPFilePath)
		return
	}

	ipStr := strings.TrimRight(string(ip), "\n")

	msg := &connectedMessage{
		s.customerID,
		s.bastionID,
		s.instanceID,
		ipStr,
		svcs,
		time.Now().Unix(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Unable to marshal message:", msg)
		log.Println(err.Error())
		return
	}

	log.Println("Publishing message:", string(msgBytes))
	if err := s.producer.Publish(nsqdTopic, msgBytes); err != nil {
		log.Println(err.Error())
	}
}

func (s *producerService) registrationLoop() {
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
func (s *producerService) Start() error {
	producer, err := nsq.NewProducer(s.nsqdHost, nsq.NewConfig())
	if err != nil {
		return err
	}
	s.producer = producer

	go s.registrationLoop()
	return nil
}

// Stop the registration loop.
func (s *producerService) Stop() error {
	s.stopChan <- struct{}{}

	return nil
}
