package daemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var datacenters = []string{"loop-api"}

func getLowestLatency() string {
	bestLatency := time.Hour
	bestDomain := "loop-api"
	client := &http.Client{Timeout: 3 * time.Second}

	for _, dc := range datacenters {
		url := fmt.Sprintf("https://%s.agglabs.com", dc)
		start := time.Now()
		_, err := client.Get(url)
		if err == nil {
			latency := time.Since(start)
			if latency < bestLatency {
				bestLatency = latency
				bestDomain = dc
			}
		}
	}
	return bestDomain
}

type DatacenterResponse struct {
	Data struct {
		Domain string `json:"domain"`
		Public string `json:"public"`
		Ws     string `json:"ws"`
	} `json:"data"`
}

type LocationResponse struct {
	Location struct {
		IP string `json:"ip"`
	} `json:"location"`
}

type AppInfoResponse struct {
	Data struct {
		NormalizedName string `json:"normalized_name"`
	} `json:"data"`
}

func getPublicIP(bestDomain string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://%s.agglabs.com/location", bestDomain)
	resp, err := client.Get(url)
	if err != nil {
		return "127.0.0.1"
	}
	defer resp.Body.Close()

	var loc LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&loc); err == nil && loc.Location.IP != "" {
		return loc.Location.IP
	}
	return "127.0.0.1"
}

func getApplicationInfo(session, appId, bestDomain string) (string, error) {
	url := fmt.Sprintf("https://%s.agglabs.com/handler?event=applications_info_event&applicationId=%s", bestDomain, appId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+session)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result AppInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.NormalizedName, nil
}

func getWebsocketURI(session, appId, bestDomain string) (string, error) {
	url := fmt.Sprintf("https://%s.agglabs.com/handler?event=applications_datacenter_event&applicationId=%s", bestDomain, appId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+session)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result DatacenterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	data := result.Data
	if data.Ws != "" {
		return fmt.Sprintf("wss://%s", data.Ws), nil
	} else if data.Public != "" {
		return fmt.Sprintf("ws://%s:3332", data.Public), nil
	}

	return "", fmt.Errorf("no websocket endpoint found in response")
}

func RunBackgroundDaemon(session, appId, scheme, addr string, inspect bool, autoRetry bool) error {
	bestDomain := getLowestLatency()

	wsURI, err := getWebsocketURI(session, appId, bestDomain)
	if err != nil {
		return fmt.Errorf("failed to resolve datacenter: %v", err)
	}

	publicIP := getPublicIP(bestDomain)
	normalizedName, err := getApplicationInfo(session, appId, bestDomain)
	if err != nil {
		normalizedName = "unknown"
	}

	localURI := fmt.Sprintf("%s://127.0.0.1:%s", scheme, addr)

	appIdPrefix := appId
	if parts := strings.Split(appId, "-"); len(parts) > 0 {
		appIdPrefix = parts[0]
	}
	publicURL := fmt.Sprintf("https://%s-%s.lf.agglabs.com", normalizedName, appIdPrefix)

	return runLiveWS(appId, wsURI, localURI, publicURL, publicIP, session, inspect, autoRetry)
}
