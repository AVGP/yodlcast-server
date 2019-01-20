package upnp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"

	"../utils"
)

// StartXMLServer starts an HTTP server to serve the XML files for UPnP
func StartXMLServer(ip *net.IP, deviceName *string) {
	addr := fmt.Sprintf("%s:8040", ip.To4())
	xmlServeMux := http.NewServeMux()
	xmlServeMux.HandleFunc("/", makeXMLHandler(ip, deviceName))
	log.Printf("XML server listening at %s\n", ip)
	http.ListenAndServe(addr, xmlServeMux)
}

func makeXMLHandler(ip *net.IP, deviceName *string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving XML %s\n", r.URL.Path)
		tplContent, err := ioutil.ReadFile("." + r.URL.Path)
		if err != nil {
			log.Printf("XML file read failed for %s: %v\n", r.URL.Path, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		content := regexp.MustCompile("\\$UUID").ReplaceAllLiteralString(string(tplContent), utils.DeviceUUID)
		content = regexp.MustCompile("\\$IP").ReplaceAllLiteralString(content, ip.To4().String())
		content = regexp.MustCompile("\\$NAME").ReplaceAllLiteralString(content, *deviceName)
		w.Write([]byte(content))
	}
}
