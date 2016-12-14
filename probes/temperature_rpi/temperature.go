package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	influx_username string //Username to connect to the server
	influx_password string //Password to connect to the server
	influx_host     string //Host to connect to
	influx_port     int    //Port to connect to
	influx_db       string //Database to connect to the server
	influx_delay    int    //Delay in seconds between each GET
	/***** SECURITY BY DEFAULT *****/
	influx_disable_ssl bool //Use HTTP instead of HTTPS for requests, default use HTTPS
	influx_unsafessl   bool //Set this when connecting to the cluster using HTTPS and not use SSL verification, default check certificate
	/*******************************/

	default_host  = "localhost"
	default_port  = 8086
	default_delay = 10

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLogs(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle, "TRACE: ", log.Lshortfile)
	Info = log.New(infoHandle, "INFO: ", 0)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "ERROR: ", log.Lshortfile)
}

func loadVars() {
	influx_username = os.Getenv("INFLUX_USERNAME")
	influx_password = os.Getenv("INFLUX_PASSWORD")
	influx_host = os.Getenv("INFLUX_HOST")
	if influx_host == "" {
		influx_host = default_host
	}
	influx_port, _ = strconv.Atoi(os.Getenv("INFLUX_PORT"))
	if influx_port == 0 {
		influx_port = default_port
	}
	influx_db = os.Getenv("INFLUX_DB")
	influx_delay, _ = strconv.Atoi(os.Getenv("INFLUX_DELAY"))
	if influx_delay == 0 {
		influx_delay = default_delay
	}

	flag.StringVar(&influx_username, "username", influx_username, "Username to connect to the server")
	flag.StringVar(&influx_password, "password", influx_password, "Password to connect to the server")
	flag.StringVar(&influx_host, "host", influx_host, "Host to connect to")
	flag.IntVar(&influx_port, "port", influx_port, "Port to connect to")
	flag.StringVar(&influx_db, "name", influx_db, "Database to connect to the server")
	flag.BoolVar(&influx_disable_ssl, "disable-ssl", false, "Use HTTP instead of HTTPS for requests, default use HTTPS")
	flag.BoolVar(&influx_unsafessl, "unsafeSsl", false, "Set this when connecting to the cluster using https and not use SSL verification")
	flag.IntVar(&influx_delay, "delay", influx_delay, "Delay in seconds between each GET")
	flag.Parse()

	if influx_unsafessl && influx_disable_ssl {
		Error.Fatal("Can't avoid certificate checks if ssl is disabled.")
	}
	if influx_db == "" {
		Error.Fatal("Please specify database name.")
	}
	if influx_username == "" && influx_password == "" {
		Info.Println("Try to connect without authentication.")
	}
}

func main() {
	InitLogs(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	loadVars()

	fmt.Println("influx_username \t= ", influx_username)
	//fmt.Println("influx_password \t= ", influx_password)
	fmt.Println("influx_unsafessl \t= ", influx_unsafessl)
	fmt.Println("influx_delay \t\t= ", influx_delay)

	secure_s := "s"
	if influx_disable_ssl {
		secure_s = ""
	}
	fmt.Printf("Connect via http%s://%s:%d/write?db=%s\n", secure_s, influx_host, influx_port, influx_db)

	var tls_config *tls.Config = nil
	if influx_unsafessl {
		tls_config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	httpClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:      fmt.Sprintf("http%s://%s:%d", secure_s, influx_host, influx_port),
		Username:  influx_username,
		Password:  influx_password,
		TLSConfig: tls_config,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: influx_db,
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	tags := map[string]string{
		"probe":  "rpi3",
		"room":   "sala",
		"sensor": "TMP102",
	}
	fields := map[string]interface{}{
		"value": 22.5,
	}

	p, err := client.NewPoint(
		"temperature",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(p)

	err = httpClient.Write(bp)
	if err != nil {
		log.Fatalln("Error: ", err)
	}
}
