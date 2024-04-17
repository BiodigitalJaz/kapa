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

	deploymentService := NewKubernetesDeploymentService(k8Client)
	podService := NewKubernetesPodService(k8Client)

	return &APIServer{
		ListenAddr:         listenAddr,
		k8Client:           k8Client,
		DeploymentServices: deploymentService,
		PodServices:        podService,
	}
}

func (a *APIServer) Run() {
	gin.SetMode(gin.ReleaseMode)
	a.PrintServerBanner()
	router := gin.Default()
	router.Static("/static", ".static")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	a.HandleCoreGroupRoutes(router)
	a.HandleAppsGroupRoutes(router)

	router.Run(":" + a.ListenAddr)
}

func (a *APIServer) HandleAppsGroupRoutes(router *gin.Engine) {
	// Group routes for deployments under "/api/v1/deployments"
	deployments := router.Group("/api/v1/apps/deployments")
	{
		// Route for creating a new Deployment
		deployments.GET("/:namespace/:name", a.getDeploymentHandler)
		// Route for creating a new Deployment
		deployments.PUT("/:namespace", a.createDeploymentHandler)
		// Route for deleting a Deployment
		deployments.DELETE("/:namespace/:name", a.deleteDeploymentHandler)
		// Route for updating a Deployment
		deployments.PATCH("/:namespace", a.patchDeploymentHandler)
	}
}

func (a *APIServer) HandleCoreGroupRoutes(router *gin.Engine) {

	pods := router.Group("/api/v1/core/pods")
	{
		pods.GET("/get-pod-logs/:podName/:namespace", a.getPodLogsHandler)
		pods.GET("/get-pods:namespace", a.getPodsHandler)
	}
	namespaces := router.Group("/api/v1/core/namespaces")
	{
		namespaces.GET("/get-namespaces", a.getNamespacesHandler)
		namespaces.GET("/get-namespace-events:namespace", a.getNamespaceEventsHander)
	}
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
