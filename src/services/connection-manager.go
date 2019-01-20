package services

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ConnectionManagerService represents an instance of the UPnP ConnectionManagerService
type ConnectionManagerService struct {
	addr net.IP
}

// Listen starts the web server
// for eventing and SOAP actions of the ConnectionManager service
func (cms ConnectionManagerService) Listen(ip *net.IP) {
	cms.addr = *ip
	addr := fmt.Sprintf("%s:8060", ip)
	log.Printf("[CMS] ConnectionManager service starting on %s\n", addr)

	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/cms_evt", cms.eventHandler)
	serverMux.HandleFunc("/cms_ctrl", cms.ctrlHandler)
	http.ListenAndServe(addr, serverMux)
}

func (ConnectionManagerService) eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "SUBSCRIBE" {
		w.WriteHeader(501)
		w.Write([]byte("Unsupported method"))
		return
	}

	log.Printf("[CMS] Event msg received:\nPath: %s\n Callback: %s\nHeader: %s\n---\n", r.URL.Path, r.Header.Get("Callback"), r.Header)
	// Let's just reject subscription requests...
	//w.WriteHeader(500)
	// TODO: Implement subscriptions
	sid, _ := uuid.NewRandom()
	w.Header().Set("SID", "uuid:"+sid.String())
	w.Header().Set("Timeout", "Second-1800")
	w.Header().Set("Server", "unix/5.1 UPnP/1.0 YodlCastServer/1.0")
	//	w.Header().Set("transferMode.dlna.org", "Streaming")
	//	w.Header().Set("contentFeatures.dlna.org", "DLNA.ORG_OP=01;DLNA.ORG_CI=0;DLNA.ORG_FLAGS=01700000000000000000000000000000")
	w.WriteHeader(200)
}

func (cms ConnectionManagerService) ctrlHandler(w http.ResponseWriter, r *http.Request) {
	action := strings.ToLower(
		regexp.MustCompile("\".*#(.+)\"").FindStringSubmatch(
			r.Header.Get("SOAPACTION"))[1])
	log.Printf("[CMS] Action received: %s\n", action)

	//	w.Header().Set("transferMode.dlna.org", "Streaming")
	//	w.Header().Set("contentFeatures.dlna.org", "DLNA.ORG_OP=01;DLNA.ORG_CI=0;DLNA.ORG_FLAGS=01700000000000000000000000000000")
	w.WriteHeader(200)

	switch action {
	case "getprotocolinfo":
		cms.getProtocolInfo(&w)
	default:
		log.Printf("[CMS] Unknown action '%s'\n", action)
	}
}

func (ConnectionManagerService) getProtocolInfo(w *http.ResponseWriter) {
	response := `<s:Envelope xmlns="urn:schemas-upnp-org:service-1-0" xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
		<s:Body>
			<u:GetProtocolInfoResponse xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1">
				<Source>http</Source>
				<Sink>http</Sink>
			</u:GetProtocolInfoResponse>
		</s:Body>
	</s:Envelope>`

	(*w).Write([]byte(response))
}

func (ConnectionManagerService) getCurrentConnectionIDs(w *http.ResponseWriter) {
	response := `<s:Envelope xmlns="urn:schemas-upnp-org:service-1-0" xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
		<s:Body>
			<u:GetCurrentConnectionIDsResponse xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1">
				<ConnectionIDs>http</ConnectionIDs>
			</u:GetCurrentConnectionIDsResponse>
		</s:Body>
	</s:Envelope>`

	(*w).Write([]byte(response))
}
func (ConnectionManagerService) getCurrentConnectionInfo(w *http.ResponseWriter) {
	response := `<s:Envelope xmlns="urn:schemas-upnp-org:service-1-0" xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
		<s:Body>
			<u:GetCurrentConnectionInfoResponse xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1">
				<RcsID>-1</RcsID>
				<AVTransportID>-1</AVTransportID>
				<ProtocolInfo></ProtocolInfo>
				<PeerConnectionManager></PeerConnectionManager>
				<PeerConnectionID>-1</PeerConnectionID>
				<Direction>Output</Direction>
				<Status>OK</Status>
			</u:GetCurrentConnectionInfoResponse>
		</s:Body>
	</s:Envelope>`

	(*w).Write([]byte(response))
}
