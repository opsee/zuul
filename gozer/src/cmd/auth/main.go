package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

const (
	version = "0.0.1"
	authURL = "https://vape.opsy.co/bastions/authenticate"
)

func validate(c *cli.Context) {
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

	resp, err := client.PostForm(authURL, url.Values{
		"id":       []string{username},
		"password": []string{password},
	})
	if err != nil {
		log.Println("Error contacting auth service.")
		log.Fatal(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "auth"
	app.Version = version
	app.Usage = "Query the authentication service to validate a username/password hash."
	app.Action = validate

	app.Run(os.Args)
}
