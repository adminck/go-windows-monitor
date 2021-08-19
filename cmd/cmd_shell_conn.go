package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/iamacarpet/go-winpty"
	"go-windows-monitor/utils/log"
	"io"
	"unicode/utf8"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type wsPty struct {
	Pty *winpty.WinPTY
	ws  *websocket.Conn
}

func (wp *wsPty) Start() {
	var err error
	// If you want to use a location other than the same folder for the DLL and exe
	// specify the path as the first param, e.g. winpty.Open(`C:\MYAPP\support`, cmdFlag)
	wp.Pty, err = winpty.OpenDefault("", "cmd")
	if err != nil {
		log.Fatalf("Failed to start command: %s\n", err)
	}
	//Set the size of the pty
	wp.Pty.SetSize(200, 60)
}

func (wp *wsPty) Stop() {
	wp.Pty.Close()

	wp.ws.Close()
}

func (wp *wsPty) readPump() {
	defer wp.Stop()

	for {
		mt, payload, err := wp.ws.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Info("conn.ReadMessage failed: %s\n", err)
				return
			}
		}
		var msg Message
		switch mt {
		case websocket.BinaryMessage:
			log.Info("Ignoring binary message: %q\n", payload)
		case websocket.TextMessage:
			err := json.Unmarshal(payload, &msg)
			if err != nil {
				log.Info("Invalid message %s\n", err)
				continue
			}
			switch msg.Type {
			case "resize":
				var size []float64
				err := json.Unmarshal(msg.Data, &size)
				if err != nil {
					log.Info("Invalid resize message: %s\n", err)
				} else {
					wp.Pty.SetSize(uint32(size[0]), uint32(size[1]))
				}
			case "data":
				var dat string
				err := json.Unmarshal(msg.Data, &dat)
				if err != nil {
					log.Info("Invalid data message %s\n", err)
				} else {
					wp.Pty.StdIn.Write([]byte(dat))
				}
			default:
				log.Info("Invalid message type %d\n", mt)
				return
			}
		default:
			log.Info("Invalid message type %d\n", mt)
			return
		}
	}
}

func (wp *wsPty) writePump() {
	defer wp.Stop()

	buf := make([]byte, 8192)
	reader := bufio.NewReader(wp.Pty.StdOut)
	var buffer bytes.Buffer
	for {
		n, err := reader.Read(buf)
		if err != nil {
			log.Info("Failed to read from pty master: %s", err)
			return
		}
		//read byte array as Unicode code points (rune in go)
		bufferBytes := buffer.Bytes()
		runeReader := bufio.NewReader(bytes.NewReader(append(bufferBytes[:], buf[:n]...)))
		buffer.Reset()
		i := 0
		for i < n {
			char, charLen, e := runeReader.ReadRune()
			if e != nil {
				log.Info("Failed to read from pty master: %s", err)
				return
			}
			if char == utf8.RuneError {
				runeReader.UnreadRune()
				break
			}
			i += charLen
			buffer.WriteRune(char)
		}
		err = wp.ws.WriteMessage(websocket.TextMessage, buffer.Bytes())
		if err != nil {
			log.Info("Failed to send UTF8 char: %s", err)
			return
		}
		buffer.Reset()
		if i < n {
			buffer.Write(buf[i:n])
		}
	}
}
