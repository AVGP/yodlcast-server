package web

import (
	"log"
	"net"
	"net/http"
)

// StartMediaServer runs an HTTP server on the given IP & port to serve the media files
func StartMediaServer(ip *net.IP, rootPath *string) {
	log.Printf("Listening on %s:8000\n", ip)
	http.ListenAndServe(ip.To4().String()+":8000", http.FileServer(http.Dir(*rootPath)))
}
