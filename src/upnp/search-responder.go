package upnp

import (
	"fmt"
	"log"
	"net"
	"regexp"
)

// UPNPMulticastAddr is the IP that is used for UPNP multicast messages, e.g. M-SEARCH
const UPNPMulticastAddr = "239.255.255.250:1900"
const maxDatagramSize = 8192

var ssdpAllMatcher = regexp.MustCompile("(?i)ST: ssdp:all")
var ssdpSearchMatcher = regexp.MustCompile("M-SEARCH")

// RespondToSearch responds to M-SEARCH requests for ssdp:all and the given searchType
// and responds with the given IP and port for the XML descriptors
func RespondToSearch(searchType string, localIP *net.IP) {
	var searchTypeMatcher = regexp.MustCompile("(?i)ST: " + searchType)
	addr, err := net.ResolveUDPAddr("udp", UPNPMulticastAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetReadBuffer(maxDatagramSize)

	for {
		buffer := make([]byte, maxDatagramSize)
		_, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		msg := string(buffer)
		if ssdpSearchMatcher.MatchString(msg) &&
			(ssdpAllMatcher.MatchString(msg) || searchTypeMatcher.MatchString(msg)) {
			log.Printf("M-SEARCH from %v: %v\n", src, string(buffer))
			sendSearchResponse(conn, src, localIP)
		}
	}
}

func sendSearchResponse(conn *net.UDPConn, dstAddr *net.UDPAddr, localIP *net.IP) {
	announcementMsg := fmt.Sprintf("%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
		"HTTP/1.1 200 OK",
		"CACHE-CONTROL: max-age=1800",
		"EXT:",
		"SERVER: unix/5.1 UPnP/1.1 YodlCastServer/1.0",
		"ST: urn:schemas-upnp-org:device:MediaServer:1",
		"USN: uuid:::urn:schemas-upnp-org:device:MediaServer:1",
		fmt.Sprintf("Location: http://%s:8040/xml/dms.xml", localIP.String()))

	conn.WriteTo([]byte(announcementMsg), dstAddr)
}
