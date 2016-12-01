# Probe template
Because all probes have common parts, for example connection to InfluxDB server, this template is intended to serve as a start point for new probes.

Data needed to connect to the server can be passed in two way: by *environment variables* or directly through *command line*.

**Command line flags override environment variables.** Thats means if you set INFLUX_PORT=80 as environment variable, to specify the port of the server, then you pass -port 666 to the command line, the probe will connect the server by port 666.

Environment variables are the following:
```bash
INFLUX_USERNAME=""  # Username to connect to the server
INFLUX_PASSWORD=""  # Password to connect to the server
INFLUX_HOST=""      # Host to connect to
INFLUX_PORT=""      # Port to connect to
INFLUX_DB=""        # Database to connect to the server
INFLUX_DELAY=""     # Delay in seconds between each GET
```
With the command line you will have the same flags more tow addittional options
```
Usage of probe:
  -username string
        Username to connect to the server
  -password string
    	Password to connect to the server
  -host string
    	Host to connect to (default "localhost")
  -port int
    	Port to connect to (default 8086)
  -name string
    	Database to connect to the server
  -delay int
    	Delay in seconds between each GET (default 10)
  -disable-ssl
    	Use HTTP instead of HTTPS for requests, default use HTTPS
  -unsafeSsl
    	Set this when connecting to the cluster using https and not use SSL verification
```
as can you see, by command line you can disable HTTPS connection for HTTP by **--disable-ssl** flag, or you can set to **not check** trust of server certificate by **--unsafeSsl**, this last one can be usefull if you use self-signed certificate.
###### Security by default
By default probe expects to connect via HTTPS and receive a certificate signed by trusted CA (i suggest [Let's Encrypt](https://letsencrypt.org/))
