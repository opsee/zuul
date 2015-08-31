package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

const (
	version = "0.0.1"
)

func validate(c *cli.Context) {
	authURL := c.String("url")

	upFile := c.Args().First()
	up, err := ioutil.ReadFile(upFile)
	if err != nil {
		log.Println("Error reading password file:", upFile)
		log.Fatal(err.Error())
	}

	split := strings.Split(string(up), "\n")
	username := split[0]
	password := split[1]

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	authStr := []byte(fmt.Sprintf(`{"id":"%s","password":"%s"}`, username, password))
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(authStr))
	if err != nil {
		log.Println("Unable to create request.")
		log.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	fmt.Println(req.Body)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error contacting auth service.")
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		os.Exit(0)
	} else {
		log.Println(resp)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "auth"
	app.Version = version
	app.Usage = "Query the authentication service to validate a username/password hash."
	app.Action = validate
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Value: os.Getenv("AUTH_URL"),
		},
	}

	app.Run(os.Args)
}
