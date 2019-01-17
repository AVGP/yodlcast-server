package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// StartContentDirectorySvc starts the web server
// for eventing and SOAP actions of the ContentDirectory service
func StartContentDirectorySvc(ip *net.IP) {
	addr := fmt.Sprintf("%s:8060", ip)
	log.Printf("ContentDirectory service starting on %s\n", addr)

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/cds_evt", eventHandler)
	serverMux.HandleFunc("/cds_ctrl", makeCtrlHandler(ip))
	http.ListenAndServe(addr, serverMux)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "SUBSCRIBE" {
		w.WriteHeader(501)
		w.Write([]byte("Unsupported method"))
		return
	}

	log.Printf("Event msg received: %s\n", r.Header.Get("Callback"))
	// Let's just reject subscription requests...
	w.WriteHeader(500) // TODO: Implement subscriptions
}

func makeCtrlHandler(ip *net.IP) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		action := strings.ToLower(
			regexp.MustCompile("\".*#(.+)\"").FindStringSubmatch(
				r.Header.Get("SOAPACTION"))[1])
		log.Printf("Action received: %s\n", action)

		w.WriteHeader(200)

		switch action {
		case "browse":
			log.Printf("Browse...\n")
			body, _ := ioutil.ReadAll(r.Body)
			log.Printf("Request body: %s\n", body)
			browse(string(body), ip, w)
		default:
			log.Printf("Unknown action '%s'\n", action)
		}
	}
}

// TODO: Does not support recursion yet :(
func browse(request string, ip *net.IP, w http.ResponseWriter) {
	mediaBaseURL := fmt.Sprintf("http://%s:8000/", ip.String())
	objectID, _ := strconv.Atoi(regexp.MustCompile("<ObjectID>(.+)</ObjectID>").FindStringSubmatch(request)[1])
	numItems := 1
	itemsResponse := ""
	// return the root folder
	if objectID == 0 {
		itemsResponse = `&lt;container id=&quot;1&quot; parentID=&quot;0&quot; restricted=&quot;1&quot;&gt;
			&lt;dc:title&gt;Photos&lt;/dc:title&gt;
			&lt;upnp:class&gt;object.container&lt;/upnp:class&gt;
		&lt;/container&gt;`
	} else {
		entries, _ := ioutil.ReadDir("./media")
		for index, entry := range entries {
			// skip meta folders
			if entry.Name() == "." || entry.Name() == "." {
				continue
			}
			itemClass := "container"
			res := ""
			if !entry.IsDir() {
				itemClass = "item"
				itemURL, _ := url.Parse(mediaBaseURL + entry.Name())
				res = fmt.Sprintf("&lt;res&gt;%s&lt;/res&gt;", itemURL.String())
			}

			itemsResponse = fmt.Sprintf(`%s&lt;%s id=&quot;%d&quot; parentID=&quot;1&quot; restricted=&quot;1&quot;&gt;
				&lt;dc:title&gt;%s&lt;/dc:title&gt;
				&lt;upnp:class&gt;object.%s&lt;/upnp:class&gt;
				%s
			&lt;/%s&gt;`, itemsResponse, itemClass, (2 + index), entry.Name(), itemClass, res, itemClass)

			numItems = len(entries)
		}
	}

	response := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<s:Envelope xmlns="urn:schemas-upnp-org:service-1-0" xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
			<s:Body>
				<u:BrowseResponse xmlns:u="urn:schemas-upnp-org:service:ContentDirectory:1">
					<Result>&lt;DIDL-Lite 
					xmlns=&quot;urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/&quot; 
					xmlns:dc=&quot;http://purl.org/dc/elements/1.1/&quot; 
					xmlns:upnp=&quot;urn:schemas-upnp-org:metadata-1-0/upnp/&quot; 
					xmlns:dlna=&quot;urn:schemas-dlna-org:metadata-1-0/&quot; 
					xmlns:sec=&quot;http://www.sec.co.kr/dlna&quot; 
					xmlns:fm=&quot;urn:schemas-jasmin-upnp.net:filemetadata/&quot; 
					xmlns:mm=&quot;urn:schemas-jasmin-upnp.net:musicmetadata/&quot; 
					xmlns:mo=&quot;urn:schemas-jasmin-upnp.net:moviemetadata/&quot;&gt;%s&lt;/DIDL-Lite&gt;</Result>
					<NumberReturned>%d</NumberReturned>
					<TotalMatches>%d</TotalMatches>
					<UpdateID>1</UpdateID>					
				</u:BrowseResponse>
			</s:Body>
		</s:Envelope>`, itemsResponse, numItems, numItems)
	w.Write([]byte(response))
}
