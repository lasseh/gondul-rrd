package main

import (
	"github.com/gin-gonic/gin"
)

const gondulURL = "https://gondul.lasse.cloud"

func main() {

	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Graph
	r.GET("/graph", graphHandler)

	// Serve the interface tree list created by the collector
	// r.GET("/tree", treeHandler)

	// Web
	r.Static("/", "../web")

	r.Run(":8080")
}
