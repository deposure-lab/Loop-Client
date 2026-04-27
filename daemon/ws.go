package daemon

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync" // <--- DODAJ TO
	"time"

	"github.com/gorilla/websocket"
)

var writeMu sync.Mutex

type InitMessage struct {
	Type      string `json:"type"`
	AppID     string `json:"appId"`
	Port      string `json:"port"`
	Token     string `json:"token"`
	Host      string `json:"host"`
	IPAddress string `json:"ip_address"`
}

type HttpRequestMsg struct {
	Type      string            `json:"type"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	RequestID string            `json:"requestId"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
}

type HttpResponseMsg struct {
	Type       string            `json:"type"`
	RequestID  string            `json:"requestId"`
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func getSystemName() string {
	return runtime.GOOS
}

func runLiveWS(appId, wsURI, localURI, publicURL, publicIP, session string, inspect bool, autoRetry bool) error {
	reconnectDelay := 3 * time.Second

	for {
		err := connectAndListen(appId, wsURI, localURI, publicURL, publicIP, session, inspect)
		if err != nil {
			fmt.Printf("WebSocket error: %v\n", err)
		}

		if !autoRetry {
			fmt.Println("Connection closed. Auto-retry is disabled (-r=false). Exiting.")
			return err
		}

		fmt.Printf("Reconnecting in %v...\n", reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}

func connectAndListen(appId, wsURI, localURI, publicURL, publicIP, session string, inspect bool) error {
	ws, _, err := websocket.DefaultDialer.Dial(wsURI, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	if inspect {
		fmt.Print("\033[H\033[2J")
		fmt.Println("──────────────────────────────────────────────────")
		fmt.Println("WEBSOCKET >> DATA")
		fmt.Println("Status:       \x1b[32m● Live\x1b[0m")
		fmt.Printf("Local URL:    \x1b[37m%s\x1b[0m\n", localURI)
		fmt.Printf("Public URL:   \x1b[95m\x1b[4m%s\x1b[0m\n", publicURL)
		fmt.Println("──────────────────────────────────────────────────")
		fmt.Println("TRAFFIC INSPECTOR")
		fmt.Println("──────────────────────────────────────────────────")
	}

	port := strings.Split(localURI, ":")[2]

	initMsg := InitMessage{
		Type:      "init",
		AppID:     appId,
		Port:      port,
		Token:     session,
		Host:      getSystemName(),
		IPAddress: publicIP, // Zaktualizowano! Korzysta z rzeczywistego IP pobranego wcześniej
	}

	writeMu.Lock()
	ws.WriteJSON(initMsg)
	writeMu.Unlock()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return err
		}

		var baseMsg struct {
			Type string `json:"type"`
		}
		json.Unmarshal(message, &baseMsg)

		if baseMsg.Type == "connection_exit" {
			fmt.Println("Connection terminated by server.")
			os.Exit(0)
		} else if baseMsg.Type == "http-request" {
			var reqMsg HttpRequestMsg
			json.Unmarshal(message, &reqMsg)

			go handleRequest(reqMsg, localURI, ws)
		}
	}
}

func handleRequest(msg HttpRequestMsg, localURI string, ws *websocket.Conn) {
	targetURL := fmt.Sprintf("%s/%s", localURI, strings.TrimPrefix(msg.Path, "/"))

	var reqBody io.Reader
	if msg.Body != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(msg.Body)
		if err != nil {
			fmt.Printf("Base64 decode error: %v\n", err)
			reqBody = bytes.NewBufferString(msg.Body)
		} else {
			reqBody = bytes.NewReader(decodedBytes)
		}
	}

	req, err := http.NewRequest(msg.Method, targetURL, reqBody)
	if err != nil {
		sendErrorResponse(ws, msg.RequestID, err)
		return
	}

	skipHeaders := map[string]bool{
		"host": true, "content-length": true, "transfer-encoding": true,
		"upgrade-insecure-requests": true, "connection": true,
	}

	for k, v := range msg.Headers {
		if !skipHeaders[strings.ToLower(k)] {
			req.Header.Set(k, v)
		}
	}

	req.Host = strings.TrimPrefix(localURI, "http://")
	req.Host = strings.TrimPrefix(req.Host, "https://")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s %s \x1b[31m502 (Local Refused)\x1b[0m\n", msg.Method, msg.Path)
		sendErrorResponse(ws, msg.RequestID, err)
		return
	}
	defer resp.Body.Close()

	statusColor := "\x1b[32m"
	if resp.StatusCode >= 400 {
		statusColor = "\x1b[31m"
	} else if resp.StatusCode >= 300 {
		statusColor = "\x1b[33m"
	}
	fmt.Printf("%s %s %s%d\x1b[0m\n", msg.Method, msg.Path, statusColor, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		respHeaders[strings.ToLower(k)] = v[0]
	}

	respMsg := HttpResponseMsg{
		Type:       "http-response",
		RequestID:  msg.RequestID,
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       base64.StdEncoding.EncodeToString(bodyBytes),
	}

	writeMu.Lock()
	err = ws.WriteJSON(respMsg)
	writeMu.Unlock()
}

func sendErrorResponse(ws *websocket.Conn, reqID string, err error) {
	bodyB64 := base64.StdEncoding.EncodeToString([]byte(err.Error()))
	respMsg := HttpResponseMsg{
		Type:       "http-response",
		RequestID:  reqID,
		StatusCode: 502,
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       bodyB64,
	}
	writeMu.Lock()
	ws.WriteJSON(respMsg)
	writeMu.Unlock()
}
