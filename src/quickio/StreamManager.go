package quickio

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"quickRadio/radioErrors"
	"strings"
	"sync"
)

// One day we may want multiple types of string managers here we will use a streamType like CloudFront in order to dictate how the tls stream should be managed
type StreamManager struct {
	mu         sync.Mutex
	streamType string
	client     *http.Client
	transport  *http.Transport
	userAgent  string
	accepts    string
}

var (
	instance *StreamManager
	once     sync.Once
)

// Fetch Content targets an active stream resource down the shared pipe
func (manager *StreamManager) FetchContent(url string) ([]byte, io.ReadCloser) {
	manager.mu.Lock()
	if manager.client == nil {
		manager.mu.Unlock()
		radioErrors.ErrorFail(fmt.Errorf("Cannot Fetch: Stream Manager Fun has been killed."))
		return nil, nil
	}
	localClient := manager.client
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		radioErrors.ErrorFail(err)
		return nil, nil
	}
	req.Header["user-agent"] = []string{manager.userAgent}
	req.Header["accept"] = []string{manager.accepts}
	resp, err := localClient.Do(req)
	if err != nil {
		// broke socket, drop internal classes and close it
		manager.transport.CloseIdleConnections()
		radioErrors.ErrorFail(err)
		return nil, nil
	}
	//defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		errorMessage := fmt.Errorf("shared pipe connection failed %d", resp.StatusCode)
		radioErrors.ErrorFail(errorMessage)
		return nil, nil
	}
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorLog(err)
	return byteValue, resp.Body
}

// Get Stream Manager yields the unified connection singleton
func GetStreamManager(streamType string, userAgent string, accepts string) *StreamManager {
	if strings.ToLower(streamType) == "cloudfront" {
		once.Do(func() {
			tlsConfig := tls.Config{
				MinVersion: tls.VersionTLS13,
				NextProtos: []string{"h2", "http:/1.1"},
			}
			//Create transport layer
			baseTransport := http.Transport{TLSClientConfig: &tlsConfig}
			baseTransport.Protocols.SetHTTP1(true)
			baseTransport.Protocols.SetHTTP2(true)
			instance = &StreamManager{
				streamType: streamType,
				client: &http.Client{
					Transport: &baseTransport,
				},
				transport: &baseTransport,
				userAgent: userAgent,
				accepts:   accepts,
			}
		})

	} else {
		errorMessage := fmt.Sprintf("Stream Type %s not supported", streamType)
		radioErrors.ErrorFail(errors.New(errorMessage))
		instance = nil
	}
	return instance
}

func (manager *StreamManager) ForceCloseIdleConnections() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.transport != nil {
		manager.transport.CloseIdleConnections()
	}
}

func (manager *StreamManager) KillFun() {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.transport != nil {
		manager.transport.CloseIdleConnections()
		manager.transport = nil
		manager.client = nil
	}
}

func NewStreamManager() *StreamManager {
	streamType := "cloudfront"
	//This collides with our instance, mat want to change user agent to safari or something we dont use
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36"
	accepts := "*/*"
	return GetStreamManager(streamType, userAgent, accepts)
}
