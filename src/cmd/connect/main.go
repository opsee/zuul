package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/codegangsta/cli"
	"github.com/opsee/zuul/registration"
)

const (
	version = "0.0.1"
)

func connect(c *cli.Context) {
	etcd := c.String("etcd-address")
	nsq := c.StringSlice("nsqlookupd-tcp-address")

	svc, err := registration.NewConsumer("connected", etcd, nsq, c.Int("consumer-concurrency"))
	if err != nil {
		log.Println("Unable to create consumer: etcd =", etcd, "nsq = ", nsq)
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	svc.Start()

	s := <-sigs
	svc.Stop() // blocks
	log.Println("Got signal, exiting:", s)
}

func main() {
	app := cli.NewApp()
	app.Name = "connect"
	app.Version = version
	app.Usage = "Consume connected messages and persist the data to Etcd"
	app.Action = connect
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "etcd-address",
			Value: os.Getenv("ETCD_ADDRESS"),
		},
		cli.StringSliceFlag{
			Name:  "nsqlookupd-tcp-address",
			Value: &cli.StringSlice{"nsqlookupd-1.opsy.co", "nsqlookupd-2.opsy.co"},
		},
		cli.IntFlag{
			Name:  "consumer-concurrency",
			Value: 10,
		},
	}

	app.Run(os.Args)
}
