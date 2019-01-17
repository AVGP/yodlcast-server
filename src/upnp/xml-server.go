package upnp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
)

// DeviceUUID is the UPnP UUID we are broadcasting and announcing
const DeviceUUID = "c000ffee-cafe-c0c0-dead-c000ffffeeee"

// StartXMLServer starts an HTTP server to serve the XML files for UPnP
func StartXMLServer(ip *net.IP) {
	addr := fmt.Sprintf("%s:8040", ip.To4())
	xmlServeMux := http.NewServeMux()
	xmlServeMux.HandleFunc("/", makeXMLHandler(ip))
	log.Printf("XML server listening at %s\n", ip)
	http.ListenAndServe(addr, xmlServeMux)
}

func makeXMLHandler(ip *net.IP) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving XML %s\n", r.URL.Path)
		tplContent, err := ioutil.ReadFile("." + r.URL.Path)
		if err != nil {
			log.Printf("XML file read failed for %s: %v\n", r.URL.Path, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		content := regexp.MustCompile("\\$UUID").ReplaceAllLiteralString(string(tplContent), DeviceUUID)
		content = regexp.MustCompile("\\$IP").ReplaceAllLiteralString(content, ip.To4().String())
		w.Write([]byte(content))
	}
}
