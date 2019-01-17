package main

import (
	"log"
	"time"

	"./services"
	"./upnp"
	"./utils"
)

func main() {
	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Fatalf("Error getting local IPv4 address: %v\n", err)
	}

	log.Printf("YodlCast server is starting on %s\n", ip.To4())

	go upnp.StartXMLServer(&ip)
	log.Printf("XML server ready...\n")

	go services.StartContentDirectorySvc(&ip)
	log.Printf("ContentDirectory service ready...\n")

	go upnp.RespondToSearch("urn:schemas-upnp-org:device:MediaServer:1", &ip)
	log.Printf("UPnP responder ready...\n")

	for true {
		time.Sleep(1 * time.Second)
	}
}
