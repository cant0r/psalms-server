# psalms-server
Serve metadata over HTTP about the currently playing music on your computer.

## Usage
1. Start the server
2. Send `GET` requests to `localhost:PORT/get-playing-psalm`
    * `PORT` is `16666` by default

### Options
```bash
  -port int
        On what port the server should listen on? (default 16666)
  -target-player string
        Your preferred media player to attach to. (default "spotify")
```

## Supported Platforms (on the `x86_64` and `arm64` architectures)

### GNU/Linux

The `psalms` server communicates through D-Bus to retrieve metadata about the currently playing media.

1. Sets up a `SessionBus` connection
2. Searches all the owned D-Bus objects for a match with the given `target-player` input parameter
    * `target-player` defaults to `spotify` by default
    *  If no MPRIS players were found then the server fucking dies
3. Calls `org.freedesktop.DBus.Properties.Get` on the retrieved MPRIS object to get the metadata from the user program.

### Windows
\<soon\>

### Darwin
N\A