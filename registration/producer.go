package registration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{"module": "registration", "event": "register", "Error": err}).Error("Error getting portmapper services")
		return
	}

	ip, err := ioutil.ReadFile(IPFilePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "registration", "event": "register", "Error": err}).Error("Error reading IP from file:", IPFilePath)
		return
	}

	if len(ip) == 0 {
		logrus.WithFields(logrus.Fields{"module": "registration", "event": "register"}).Warn("IP file empty: ", IPFilePath)
		return
	}

	ipStr := strings.TrimRight(string(ip), "\n")

	msg := &ConnectedMessage{
		s.customerID,
		s.bastionID,
		s.instanceID,
		ipStr,
		svcs,
		time.Now().Unix(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "registration", "event": "register", "Error": err}).Error("Unable to marshal message: ", string(msgBytes))
		return
	}

	logrus.WithFields(logrus.Fields{"module": "registration", "event": "register"}).Info("Publishing message: ", string(msgBytes))
	if err := s.producer.Publish(nsqdTopic, msgBytes); err != nil {
		logrus.WithFields(logrus.Fields{"module": "registration", "event": "register", "Error": err}).Error("Error Publishing message: ", string(msgBytes))
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
