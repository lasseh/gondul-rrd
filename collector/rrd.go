package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ziutek/rrd"
)

const (
	rrdStep      = 60
	rrdHeartbeat = 10 * rrdStep
)

// UpdateRRD should return errors
// TODO implement return errors
func UpdateRRD(rrdPath, device, iface string, in, out uint64) {
	// Prettify interface names
	iface = strings.Replace(iface, "/", "", -1)

	// rrd paths
	folder := fmt.Sprintf("%s%s/", rrdPath, device)
	filename := fmt.Sprintf("%s%s/%s.rrd", rrdPath, device, iface)

	// Create folder if it doesn't exist
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		log.Println("Creating folder:", folder)
		os.MkdirAll(folder, 0744)
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		CreateRRD(filename)
		// exit loop since we just created the file
		// insert the data next iteration
		return
	}

	// Update rrd with data
	u := rrd.NewUpdater(filename)
	err := u.Update(time.Now(), in, out)
	if err != nil {
		log.Println("Error updating rrd:", err)

	}
	log.Println("Updated file", filename)
}

// CreateRRD should return errors
// TODO implement return errors
func CreateRRD(filename string) {
	c := rrd.NewCreator(filename, time.Now(), rrdStep)
	c.DS("traffic_in", "COUNTER", rrdHeartbeat, 0, 1250000000000)
	c.DS("traffic_out", "COUNTER", rrdHeartbeat, 0, 1250000000000)
	c.RRA("MIN", 0, 360, 576)
	c.RRA("MIN", 0, 30, 576)
	c.RRA("MIN", 0, 7, 576)
	c.RRA("MIN", 0, 1, 576)
	c.RRA("AVERAGE", 0, 360, 576)
	c.RRA("AVERAGE", 0, 30, 576)
	c.RRA("AVERAGE", 0, 7, 576)
	c.RRA("AVERAGE", 0, 1, 576)
	c.RRA("MAX", 0, 360, 576)
	c.RRA("MAX", 0, 30, 576)
	c.RRA("MAX", 0, 7, 576)
	c.RRA("MAX", 0, 1, 576)
	// Change to c.Create(true) if you want to overwrite the file
	// Only if you want to keep data but change rrd layout
	err := c.Create(false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created file", filename)
	return
}
