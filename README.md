# yodlcast-server
UPnP media server

## Current state

This is *incomplete* and highly *experimental* at the moment.

Currently implemented features:

* The basic UPnP functionality to respond to searches.
* HTTP server for the device and service XML descriptions.
* A very basic and rough UPnP ContentDirectory service without nested folders or error checking.

The most important missing bits right now:

* HTTP server to serve the media files.
* Proper directory listing with nesting for the ContentDirectory service.
* Error handling for the ContentDirectoryService

# The Goal

The goal of this project is to provide a solid UPnP media server and media renderer.
It should also be DLNA compatible and provide a convenient web interface to manage it.