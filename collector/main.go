package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Device is the sqlite structure
type Device struct {
	Timestamp     time.Time `json:"timestamp"`
	Device        string    `json:"device"`
	Interface     string    `json:"interface"`
	IfHCInOctets  uint64    `json:"ifHCInOctets"`
	IfHCOutOctets uint64    `json:"ifHCOutOctets"`
	IfAlias       string    `json:"ifAlias"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("sqlite3", "devices.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.AutoMigrate(&Device{})
	db.LogMode(true)

	// Flags
	URL := flag.String("url", "http://gondul.lasse.cloud/", "URL to gondul")
	rrdPath := flag.String("path", "rrd/", "Path to rrd files")
	Username := flag.String("username", "tech", "username to gondul")
	Password := flag.String("password", "rules", "password to gondul")
	sleep := flag.Int("sleep", 10, "Poller sleep time")
	flag.Parse()

	// Map for switch time
	switchTime := make(map[string]string)

	g := NewGondul(*URL, *Username, *Password)

	for {
		err := g.PollData()
		if err != nil {
			log.Println("Failed to poll data from API:", err)
			os.Exit(1)
		}
		for k, sw := range g.Switches {

			// Check if switch exists in map
			value, exist := switchTime[k]
			if !exist {
				switchTime[k] = g.Switches[k].Time
				continue
			}

			// skip switch if the time is the same
			if sw.Time == value {
				continue
			}

			// Switch has new timestamp
			// Update rrd
			for name, iface := range sw.Ifs {
				if iface.IfHCInOctets != 0 || iface.IfHCOutOctets != 0 {
					UpdateRRD(*rrdPath, k, name, iface.IfHCInOctets, iface.IfHCOutOctets)
					UpdateDB(k, name, iface.IfAlias, iface.IfHCInOctets, iface.IfHCOutOctets)
				}
			}

			if sw.Uplinks.IfHCInOctets != 0 || sw.Uplinks.IfHCOutOctets != 0 {
				UpdateRRD(*rrdPath, k, "uplinks", sw.Uplinks.IfHCInOctets, sw.Uplinks.IfHCOutOctets)
				UpdateDB(k, "uplinks", "Uplinks", sw.Clients.IfHCInOctets, sw.Clients.IfHCOutOctets)
			}

			if sw.Clients.IfHCInOctets != 0 || sw.Clients.IfHCOutOctets != 0 {
				UpdateRRD(*rrdPath, k, "clients", sw.Clients.IfHCInOctets, sw.Clients.IfHCOutOctets)
				UpdateDB(k, "clients", "Clients", sw.Clients.IfHCInOctets, sw.Clients.IfHCOutOctets)
			}

			if sw.Totals.IfHCInOctets != 0 || sw.Totals.IfHCOutOctets != 0 {
				UpdateRRD(*rrdPath, k, "totals", sw.Totals.IfHCInOctets, sw.Totals.IfHCOutOctets)
				UpdateDB(k, "totals", "Totals", sw.Totals.IfHCInOctets, sw.Totals.IfHCOutOctets)
			}

			// VCP Interfaces
			for ifaceID, subiface := range sw.Vcp.VcpIntIn { // Interface id: 0, 1, 2
				for vcpName, vcpInOctet := range subiface { // vcp-255/0/25
					// Get OutOctets from map
					vcpOutOctet := sw.Vcp.VcpIntOut[ifaceID][vcpName]
					// Convert string to uint64
					uintInOctets, _ := strconv.ParseUint(vcpInOctet, 10, 64)
					uintOutOctets, _ := strconv.ParseUint(vcpOutOctet, 10, 64)

					if uintInOctets != 0 || uintOutOctets != 0 {
						filename := ifaceID + "-" + vcpName
						filename = strings.Replace(filename, "/", "_", -1)
						UpdateRRD(*rrdPath, k, filename, uintInOctets, uintOutOctets)
						UpdateDB(k, vcpName, vcpName, uintInOctets, uintOutOctets)
					}
				}
			}

			// Update timestamp map
			switchTime[k] = g.Switches[k].Time
		}

		// Gondul API time
		gondulTime := time.Unix(g.Time, 0)
		log.Println("Gondul timestamp:", gondulTime)

		// Sleep
		time.Sleep(time.Duration(*sleep) * time.Second)
	}
}

// UpdateDB should return errors
// TODO implement return errors
func UpdateDB(device, iface, alias string, in, out uint64) {
	var d Device
	db.Where(Device{
		Device:    device,
		Interface: iface,
	}).Assign(
		Device{
			IfHCInOctets:  in,
			IfHCOutOctets: out,
			IfAlias:       alias,
			Timestamp:     time.Now(),
		}).FirstOrCreate(&d)
}
