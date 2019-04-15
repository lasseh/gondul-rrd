package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ziutek/rrd"
)

func graphHandler(c *gin.Context) {
	rrdStart := c.DefaultQuery("start", "-7200s")
	rrdEnd := c.DefaultQuery("end", "-60s")
	rrdLegend := c.DefaultQuery("legend", "1")
	rrdTitle := c.DefaultQuery("title", "Network Traffic")
	rrdWidth := c.DefaultQuery("width", "600")
	rrdHeight := c.DefaultQuery("height", "200")
	rrdDevice := c.Query("device")
	rrdInterface := c.Query("interface")
	rrdTheme := c.Query("theme")

	if len(rrdDevice) == 0 {
		//TODO Return a rrd image with "no data found"
		c.String(http.StatusOK, "Not found")
		return
	}
	if len(rrdInterface) == 0 {
		//TODO Return a rrd image with "no data found"
		c.String(http.StatusOK, "Not found")
		return
	}
	// Prettify interface names
	iface := strings.Replace(rrdInterface, "/", "", -1)
	filename := fmt.Sprintf("%s%s/%s.rrd", rrdLol, rrdDevice, iface)

	fmt.Println("Filename:", filename)

	// Convert to int
	width, _ := strconv.ParseInt(rrdWidth, 10, 64)
	width2 := uint(width)
	height, _ := strconv.ParseInt(rrdHeight, 10, 64)
	height2 := uint(height)

	// Graph
	g := rrd.NewGrapher()
	g.SetImageFormat("PNG")
	g.SetBase(1000)
	g.SetSize(width2, height2)
	g.SetTitle(rrdTitle)
	g.SetVLabel("bits per second")
	g.SetSlopeMode()
	g.SetRigid()
	g.SetAltAutoscaleMax()
	//g.SetWatermark("Nora <3")

	if rrdLegend == "0" {
		g.SetNoLegend()
	}
	// Nightmode
	// not finished
	if rrdTheme == "dark" {
		g.SetColor("BACK", "31313DFF")
		g.SetColor("FONT", "BEC2CCFF")
		g.SetColor("CANVAS", "31313DFF")
		g.SetColor("SHADEA", "31313DFF")
		g.SetColor("SHADEB", "31313DFF")
		g.SetColor("GRID", "31313DFF")
		g.SetColor("MGRID", "FF7E2730")
		g.SetColor("AXIS", "31313DFF")
		g.SetColor("ARROW", "A500A5FF")
	}
	if rrdTheme == "old" {
		//g.SetColor("BACK", "141474FF")
		//g.SetColor("FONT", "F3FF3FFF")
		//g.SetColor("CANVAS", "313131FF")
		//g.SetColor("SHADEA", "343834FF")
		//g.SetColor("SHADEB", "343434FF")
		//g.SetColor("GRID", "313231FF")
		//g.SetColor("MGRID", "292929FF")
		//g.SetColor("AXIS", "999999FF")
		//g.SetColor("FRAME", "999999FF")
		//g.SetColor("ARROW", "999999FF")
	}

	g.Def("a", filename, "traffic_out", "MAX")
	g.Def("b", filename, "traffic_in", "AVERAGE")
	g.Def("c", filename, "traffic_out", "AVERAGE")
	g.CDef("cdefd", "b,8,*")
	g.CDef("cdefe", "c,8,*")
	// Inbound
	if rrdTheme == "dark" {
		g.Area("cdefd", "A500A550", "")
		g.Line(2, "cdefd", "A500A5FF", "Inbound\t")
	} else if rrdTheme == "old" {
		g.Area("cdefd", "00CF0050", "")
		g.Line(1, "cdefd", "00CF00FF", "Inbound\t")
	} else {
		g.Area("cdefd", "00CF0050", "")
		g.Line(1, "cdefd", "00CF00FF", "Inbound\t")
	}
	g.VDef("c_in", "cdefd,LAST")
	g.VDef("a_in", "cdefd,AVERAGE")
	g.VDef("m_in", "cdefd,MAXIMUM")
	g.GPrint("c_in", "Current\\:%8.2lf%s\t")
	g.GPrint("a_in", "Average\\:%8.2lf%s\t")
	g.GPrint("m_in", "Maximum\\:%8.2lf%s\\n")
	// Outbound
	if rrdTheme == "dark" {
		g.Area("cdefe", "31c0f640", "")
		g.Line(2, "cdefe", "31c0f6FF", "Outbound\t")
	} else if rrdTheme == "old" {
		g.Line(1, "cdefe", "212125FF", "Outbound\t")
	} else {
		g.Area("cdefe", "002A9750", "")
		g.Line(1, "cdefe", "002A97FF", "Outbound\t")
	}
	g.VDef("c_out", "cdefe,LAST")
	g.VDef("a_out", "cdefe,AVERAGE")
	g.VDef("m_out", "cdefe,MAXIMUM")
	g.GPrint("c_out", "Current\\:%8.2lf%s\t")
	g.GPrint("a_out", "Average\\:%8.2lf%s\t")
	g.GPrint("m_out", "Maximum\\:%8.2lf%s\\n")

	now := time.Now()

	startTime, _ := time.ParseDuration(rrdStart)
	endTime, _ := time.ParseDuration(rrdEnd)

	_, rrdImage, _ := g.Graph(now.Add(startTime), now.Add(endTime))

	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Write(rrdImage)
}
