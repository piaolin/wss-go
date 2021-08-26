package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var (
	addr     = flag.String("addr", "0.0.0.0:443", "https service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	newLineLength int
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.InfoLevel)

	goos := runtime.GOOS
	switch goos {
	case "windows":
		newLineLength = 2
	default:
		newLineLength = 1
	}
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

//write message
func pumpStdout(ws *websocket.Conn, done chan struct{}) {
	defer func() {
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Shell > ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-newLineLength]

		if text == "exit" {
			break
		}

		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, []byte(text)); err != nil {
			ws.Close()
			log.Errorln(err)
			break
		}
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

// receive message
func pumpStdin(ws *websocket.Conn, done chan struct{}) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			break
		case <-ticker.C:
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Fatalln(err)
			}
			if len(message) == 0 {
				message = []byte("Nothing\n")
			}
			fmt.Printf("\r\n%sShell > ", message)
		}
	}
}

func WebsocketHandle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	stdoutDone := make(chan struct{})
	go pumpStdout(c, stdoutDone)
	go ping(c, stdoutDone)

	pumpStdin(c, stdoutDone)

}

func main() {
	flag.Parse()

	log.Printf("listen at %s", *addr)
	http.HandleFunc("/", WebsocketHandle)
	log.Fatalln(http.ListenAndServeTLS(*addr, "certificate/ssl.crt", "certificate/ssl.key", nil))
}
