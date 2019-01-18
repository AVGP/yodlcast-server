# yodlcast-server
UPnP media server

## Current state

This is *incomplete* and highly *experimental* at the moment.

Currently implemented features:

* The basic UPnP functionality to respond to searches.
* HTTP server for the device and service XML descriptions.
* A very basic and rough UPnP ContentDirectory service without nested folders or error checking.

The most important missing bits right now:

* Proper directory listing with nesting for the ContentDirectory service.
* Handle errors in ContentDirectory service.
* Implement events for ContentDirectory service.
* Implement the ConnectionManager service.
* Implement the AVTransport service.
* Allow to upload or remove files and folders via UPnP.

# The Goal

The goal of this project is to provide a solid UPnP media server and media renderer.
It should also be DLNA compatible and provide a convenient web interface to manage it.

# Building it
To build the server, run `go build -o yodlserver src/main.go`. This will create the `yodlserver` executable.

# Running it
By default, running `yodlserver` will create a UPnP media server with the name "YodlCast".
It will attempt to serve files from a `media` directory in the same directory as the `yodlserver` binary by default. The following configuration flags exist:

* `--name` sets the name that the server shows to other devices.
* `--rootdir` sets the directory to serve media files from.