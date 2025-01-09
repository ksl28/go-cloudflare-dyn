package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type Response struct {
	Result  []DNSRecord `json:"result"`
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
}

func getDNS(zoneID string, apiKey string) ([]DNSRecord, error) {
	cloudUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", cloudUrl, nil)
	if err != nil {
		return nil, err
	}
	bearerKey := fmt.Sprintf("Bearer %s", apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearerKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned errors: %v", response.Errors)
	}

	return response.Result, nil
}

func getCurrentPublicIP() (string, error) {
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		return "", fmt.Errorf("error fetching public IP: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return strings.TrimSpace(string(body)), nil
}

func updateDNS(recordID, zoneID, recordType, recordName, recordContent string, recordTTL int, recordProxied bool, apiKey string) {
	cloudUrl := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	body := DNSRecord{
		Type:    recordType,
		Name:    recordName,
		Content: recordContent,
		TTL:     recordTTL,
		Proxied: recordProxied,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling body:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", cloudUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	bearerKey := fmt.Sprintf("Bearer %s", apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearerKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	var records []string
	var apiKey, zoneID string
	var refreshSeconds int

	flag.Func("records", "Specify one or more records (use -records multiple times)", func(value string) error {
		records = append(records, value)
		return nil
	})
	flag.StringVar(&apiKey, "apiKey", "", "Enter the Cloudflare API key")
	flag.StringVar(&zoneID, "zoneID", "", "Enter the Cloudflare zoneID")
	flag.IntVar(&refreshSeconds, "refreshSeconds", 300, "Enter the number of seconds that the program will sleep, before repeating again.")
	flag.Parse()

	fmt.Printf("The records entered was: %s \n", strings.Join(records, ", "))

	for {
		dnsRecords, err := getDNS(zoneID, apiKey)
		if err != nil {
			fmt.Println("Error fetching DNS records:", err)
			return
		}

		currentIP, err := getCurrentPublicIP()
		if err != nil {
			fmt.Println("Error fetching public IP:", err)
			return
		}

		for _, record := range dnsRecords {
			for _, validEntry := range records {
				if record.Name == validEntry {
					if record.Content != currentIP {
						log.Printf("The record %s is not correct", record.Name)
						updateDNS(record.ID, zoneID, record.Type, record.Name, currentIP, record.TTL, record.Proxied, apiKey)
					}
				}
			}
		}
		time.Sleep(time.Duration(refreshSeconds) * time.Second)
	}

}
