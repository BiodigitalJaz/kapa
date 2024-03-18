package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewAPIServer(listenAddr string) *APIServer {

	k8Client, err := setupK8Access()
	if err != nil {
		log.Printf("error setting up k8's client: %s\n", err)
	}

	return &APIServer{
		ListenAddr: listenAddr,
		k8Client:   k8Client,
	}
}

func (a *APIServer) Run() {
	gin.SetMode(gin.ReleaseMode)
	a.PrintServerBanner()
	router := gin.Default()
	router.Static("/static", ".static")
	router.LoadHTMLGlob("templates/*.html")

	a.HandleRoutes(router)

	router.Run(":" + a.ListenAddr)
}

func (a *APIServer) HandleRoutes(router *gin.Engine) {

	router.GET("/k/get-pod-logs/:podName/:namespace", a.getPodLogsHandler)
	router.GET("/k/get-pods:namespace", a.getPodsHandler)
	router.GET("/k/get-namespaces", a.getNamespacesHandler)
	router.GET("/k/get-namespace-events:namespace", a.getNamespaceEventsHander)

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
}
func (a *APIServer) PrintServerBanner() {
	fmt.Printf(`
	Starting service on port %s
		╦╔═╔═╗╔═╗╔═╗
		╠╩╗╠═╣╠═╝╠═╣
		╩ ╩╩ ╩╩  ╩ ╩
	   Start Time: %s
	__________________________	
	`, a.ListenAddr, time.Now().Format("3:04PM"))
	fmt.Println()
}
