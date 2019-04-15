package main

import (
	"flag"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var rrdLol = ""
var dbLol = ""
var db *gorm.DB
var err error

// Device is the sqlite structure
type Device struct {
	Timestamp     time.Time `json:"timestamp"`
	Device        string    `json:"device"`
	Interface     string    `json:"interface"`
	IfHCInOctets  uint64    `json:"ifHCInOctets"`
	IfHCOutOctets uint64    `json:"ifHCOutOctets"`
	IfAlias       string    `json:"ifAlias"`
}

// Datalol is ..
type Datalol struct {
	Switches map[string]*Switch `json:"switches"`
}

// Switch is
type Switch struct {
	Ifs    map[string]*IfStat `json:"ifs" db:"interface"`
	Totals *Total             `json:"totals"`
}

// IfStat is
type IfStat struct {
	Timestamp     time.Time `json:"timestamp"`
	IfHCInOctets  uint64    `json:"ifHCInOctets"`
	IfHCOutOctets uint64    `json:"ifHCOutOctets"`
	IfAlias       string    `json:"ifAlias"`
}

// Total is ...
type Total struct {
	IfHCInOctets  uint64 `json:"ifHCInOctets"`
	IfHCOutOctets uint64 `json:"ifHCOutOctets"`
	Total         uint64 `json:"total"`
	Live          uint64 `json:"live"`
}

func main() {
	// Flags
	rrdPath := flag.String("path", "rrd/", "Path to rrd files")
	dbPath := flag.String("db", "/var/lib/gondul-rrd.sqlite", "Path to database")
	flag.Parse()
	rrdLol = *rrdPath
	dbLol = *dbPath

	// Database
	db, err = gorm.Open("sqlite3", *dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	//db.LogMode(true)

	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Graph
	r.GET("/graph", graphHandler)

	// Serve the interface tree list created by the collector
	r.GET("/switches", switchHandler)

	// Index
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"switches": "http://rrd.lasse.cloud/switches",
			"graph":    "https://rrd.lasse.cloud/graph?width=1700&height=600&legend=1&start=-24h&end=-1m&device=r1.tele&interface=ae0&title=Interwebz&theme=dark",
		})
	})

	r.Run(":8080")
}

func switchHandler(c *gin.Context) {
	var d []Device
	sw := NewDatalol()
	if err := db.Table("devices").Find(&d); err != nil {
		c.AbortWithStatus(404)
	}

	for _, v := range d {
		if _, ok := sw.Switches[v.Device]; !ok {
			sw.Switches[v.Device] = NewSwitch()
		}
		sw.Switches[v.Device].Ifs[v.Interface] = v.GetIfStat()
	}

	for _, v := range sw.Switches {
		t := &Total{}
		v.Totals = t
		for _, values := range v.Ifs {
			t.Total++
			t.IfHCOutOctets += values.IfHCOutOctets
			t.IfHCInOctets += values.IfHCInOctets
		}
	}

	c.JSON(200, sw)
}

// GetIfStat populates ..
func (d *Device) GetIfStat() *IfStat {
	return &IfStat{
		Timestamp:     d.Timestamp,
		IfHCInOctets:  d.IfHCInOctets,
		IfHCOutOctets: d.IfHCOutOctets,
		IfAlias:       d.IfAlias,
	}
}

// NewSwitch creates a new switch
func NewSwitch() *Switch {
	var s Switch
	s.Ifs = make(map[string]*IfStat)
	return &s
}

// NewDatalol is ..
func NewDatalol() *Datalol {
	var d Datalol
	d.Switches = make(map[string]*Switch)
	return &d
}
