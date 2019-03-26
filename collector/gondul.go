package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Gondul is ...
type Gondul struct {
	APIURL   string
	Username string
	Password string
	Hash     string               `json:"hash"`
	Switches map[string]*switches `json:"switches"`
	Time     int64                `json:"time,omitempty"`
}

// NewGondul returns a new Gondul
func NewGondul(apiurl, username, password string) *Gondul {
	return &Gondul{
		APIURL:   apiurl,
		Username: username,
		Password: password,
	}
}

// PollData populates the Gondul instance by polling
// the API
func (g *Gondul) PollData() error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", g.APIURL+"/api/public/switch-state", nil)

	req.SetBasicAuth(g.Username, g.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	// Don't forget to close the response body
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&g)
	if err != nil {
		return err
	}

	return nil
}

type switches struct {
	Time    string          `json:"time"`
	Temp    string          `json:"temp"`
	Ifs     map[string]*ifs `json:"ifs"`
	Clients clients         `json:"clients"`
	Vcp     vcp             `json:"vcp,omitempty"`
	Totals  totals          `json:"totals"`
	Uplinks uplinks         `json:"uplinks"`
}

type ifs struct {
	IfHCOutOctets uint64 `json:"ifHCOutOctets"`
	IfHCInOctets  uint64 `json:"ifHCInOctets"`
}

type vcp struct {
	VcpIntOut map[string]map[string]string `json:"jnxVirtualChassisPortOutOctets,omitempty"`
	VcpIntIn  map[string]map[string]string `json:"jnxVirtualChassisPortInOctets,omitempty"`
}

type clients struct {
	IfHCOutOctets uint64 `json:"ifHCOutOctets"`
	IfHCInOctets  uint64 `json:"ifHCInOctets"`
}

type totals struct {
	IfHCOutOctets uint64 `json:"ifHCOutOctets"`
	IfHCInOctets  uint64 `json:"ifHCInOctets"`
}

type uplinks struct {
	IfHCOutOctets uint64 `json:"ifHCOutOctets"`
	IfHCInOctets  uint64 `json:"ifHCInOctets"`
}
