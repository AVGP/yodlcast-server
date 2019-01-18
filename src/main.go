package main

import (
	"flag"
	"log"
	"time"

	"./services"
	"./upnp"
	"./utils"
	"./web"
)

func main() {
	var mediaRoot string
	var deviceName string

	flag.StringVar(&mediaRoot, "rootdir", "./media", "The root directory for media files")
	flag.StringVar(&deviceName, "name", "YodlCast", "The name of the media server that other devices will display")

	flag.Parse()

	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Fatalf("Error getting local IPv4 address: %v\n", err)
	}

	log.Printf("YodlCast server is starting on %s\n", ip.To4())

	go upnp.StartXMLServer(&ip, &deviceName)
	log.Printf("XML server ready...\n")

	go web.StartMediaServer(&ip, &mediaRoot)
	log.Printf("XML server ready...\n")

	go services.StartContentDirectorySvc(&ip, &mediaRoot)
	log.Printf("ContentDirectory service ready...\n")

	go upnp.RespondToSearch("urn:schemas-upnp-org:device:MediaServer:1", &ip)
	log.Printf("UPnP responder ready...\n")

	for true {
		time.Sleep(1 * time.Second)
	}
}
