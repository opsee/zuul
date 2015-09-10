package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/codegangsta/cli"
	"github.com/opsee/gozer/registration"
)

const (
	version = "0.0.1"
)

func register(c *cli.Context) {
	interval := time.Second * 30

	customerID := c.String("customer-id")
	bastionID := c.String("bastion-id")
	instanceID := c.String("instance-id")
	nsqdHost := c.String("nsqd-host")
	registration.IPFilePath = c.String("ip-file-path")
	//func NewService(interval time.Duration, nsqdHost string, customerID string, bastionID string, instanceID string) *nsqdService {
	svc := registration.NewProducer(interval, nsqdHost, customerID, bastionID, instanceID)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	svc.Start()

	s := <-sigs
	svc.Stop() // blocks
	log.Println("Got signal, exiting:", s)
}

func main() {
	app := cli.NewApp()
	app.Name = "register"
	app.Usage = "Announce that you have connected by sending a message to NSQD"
	app.Version = version
	app.Action = register
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "customer-id, c",
			Value: os.Getenv("CUSTOMER_ID"),
		},
		cli.StringFlag{
			Name:  "bastion-id, b",
			Value: os.Getenv("BASTION_ID"),
		},
		cli.StringFlag{
			Name:  "instance-id, i",
			Value: os.Getenv("AWS_INSTANCE_ID"),
		},
		cli.StringFlag{
			Name:  "nsqd-host, n",
			Value: os.Getenv("NSQD_HOST"),
		},
		cli.StringFlag{
			Name:  "ip-file-path, f",
			Value: "/gozer/state/ip",
		},
	}

	app.Run(os.Args)
}
