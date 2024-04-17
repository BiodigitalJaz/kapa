package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (a *APIServer) deleteDeploymentHandler(ctx *gin.Context) {
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")
	name := strings.TrimPrefix(ctx.Param("name"), ":")

	if err := a.DeploymentServices.DeleteDeployment(namespace, name); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"deployment error": err.Error()})
		return

	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deployment deleted successfully"})
}

func (a *APIServer) getDeploymentHandler(ctx *gin.Context) {
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")
	name := strings.TrimPrefix(ctx.Param("name"), ":")

	deployment, err := a.DeploymentServices.GetDeployment(namespace, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"unable to retrieve deployment error": err.Error()})
		return
	}

	ctx.JSON(http.StatusFound, deployment)
}

func (a *APIServer) patchDeploymentHandler(ctx *gin.Context) {
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")

	// Read YAML from request body
	deploymentYAML, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	deploymentObj, err := convertYAMLToTypedDeployment(deploymentYAML)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"YAML conversion error": err.Error()})
		return
	}

	if err := a.DeploymentServices.UpdateDeployment(namespace, &deploymentObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"deployment error": err.Error()})
		return

	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deployment updated successfully"})
}

func (a *APIServer) createDeploymentHandler(ctx *gin.Context) {
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")

	// Read YAML from request body
	deploymentYAML, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	deploymentObj, err := convertYAMLToTypedDeployment(deploymentYAML)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"YAML conversion error": err.Error()})
		return
	}

	if err := a.DeploymentServices.CreateDeployment(namespace, &deploymentObj); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"deployment error": err.Error()})
		return

	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"message": "Deployment created successfully"})
}

func (a *APIServer) getPodLogsHandler(ctx *gin.Context) {
	podName := strings.TrimPrefix(ctx.Param("podName"), ":")
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")

	logs, err := a.PodServices.GetPodLogs(podName, namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}

func (a *APIServer) getNamespaceEventsHander(ctx *gin.Context) {
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")

	events, err := a.NamespaceServices.GetNamespaceEvents(namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

func (a *APIServer) getNamespacesHandler(ctx *gin.Context) {
	namespaces, err := a.NamespaceServices.GetNamespaces()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"namespaces": namespaces,
	})
}

func (a *APIServer) getPodsHandler(ctx *gin.Context) {
	namespace := strings.TrimPrefix(ctx.Param("namespace"), ":")

	pods, err := a.PodServices.GetPods(namespace)
	if err != nil {
		log.Printf("error getting top pods: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		namespace: pods,
	})
}
