package upnp

import (
	"fmt"
	"log"
	"net"

	"../utils"
)

// Announce sends SSDP discovery announcements for the media server root device and services
func Announce(localIP *net.IP) {
	log.Printf("Announcing to UPnP\n")
	addr, _ := net.ResolveUDPAddr("udp", UPNPMulticastAddr)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Cannot announce to UPnP: %v\n", err)
	}

	notifyRootDevice(conn, localIP)
	notifyService(conn, localIP, "urn:schemas-upnp-org:service:ContentDirectory:1")
	notifyService(conn, localIP, "urn:schemas-upnp-org:service:ConnectionManager:1")
}

func notifyRootDevice(conn *net.UDPConn, localIP *net.IP) {
	rootDeviceMsg := fmt.Sprintf("%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
		"NOTIFY * HTTP/1.1",
		"HOST: "+UPNPMulticastAddr,
		"CACHE-CONTROL: max-age=1800",
		"SERVER: unix/5.1 UPnP/1.0 YodlCastServer/1.0 DLNADOC/1.50",
		"NT: upnp:rootdevice",
		"NTS: ssdp:alive",
		"USN: uuid:"+utils.DeviceUUID+"::upnp:rootdevice",
		fmt.Sprintf("Location: http://%s:8040/xml/dms.xml", localIP.String()))

	deviceUUIDMsg := fmt.Sprintf("%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
		"NOTIFY * HTTP/1.1",
		"HOST: "+UPNPMulticastAddr,
		"CACHE-CONTROL: max-age=1800",
		"SERVER: unix/5.1 UPnP/1.0 YodlCastServer/1.0 DLNADOC/1.50",
		"NT: uuid:"+utils.DeviceUUID,
		"NTS: ssdp:alive",
		"USN: uuid:"+utils.DeviceUUID,
		fmt.Sprintf("Location: http://%s:8040/xml/dms.xml", localIP.String()))

	deviceTypeMsg := fmt.Sprintf("%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
		"NOTIFY * HTTP/1.1",
		"HOST: "+UPNPMulticastAddr,
		"CACHE-CONTROL: max-age=1800",
		"SERVER: unix/5.1 UPnP/1.0 YodlCastServer/1.0 DLNADOC/1.50",
		"NT: urn:schemas-upnp-org:device:MediaServer:1",
		"NTS: ssdp:alive",
		"USN: uuid:::urn:schemas-upnp-org:device:MediaServer:1",
		fmt.Sprintf("Location: http://%s:8040/xml/dms.xml", localIP.String()))

	fmt.Fprintf(conn, rootDeviceMsg)
	fmt.Fprintf(conn, deviceUUIDMsg)
	fmt.Fprintf(conn, deviceTypeMsg)
}

func notifyService(conn *net.UDPConn, localIP *net.IP, serviceName string) {
	serviceMsg := fmt.Sprintf("%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n",
		"NOTIFY * HTTP/1.1",
		"HOST: "+UPNPMulticastAddr,
		"CACHE-CONTROL: max-age=1800",
		"SERVER: unix/5.1 UPnP/1.0 YodlCastServer/1.0 DLNADOC/1.50",
		"NT: "+serviceName,
		"NTS: ssdp:alive",
		"USN: uuid:"+utils.DeviceUUID+"::"+serviceName,
		fmt.Sprintf("Location: http://%s:8040/xml/dms.xml", localIP.String()))

	fmt.Fprintf(conn, serviceMsg)

}
