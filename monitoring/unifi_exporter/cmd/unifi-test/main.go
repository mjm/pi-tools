package main

import (
	"log"
	"os"
	"time"

	"github.com/mdlayher/unifi"
)

func main() {
	username := "mmoriarity"
	password := os.Getenv("UNIFI_PASSWORD")
	log.Printf("%s:%s", username, password)

	c, err := unifi.NewClient("https://10.0.0.1", unifi.InsecureHTTPClient(10*time.Second))
	if err != nil {
		log.Panicf("creating client: %v", err)
	}
	c.UnifiOS = true

	if err := c.Login(username, password); err != nil {
		log.Panicf("logging in: %v", err)
	}

	devices, err := c.Devices("default")
	if err != nil {
		log.Panicf("listing devices: %v", err)
	}

	for _, d := range devices {
		log.Printf("%s", d.Name)
	}
}
