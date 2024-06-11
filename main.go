package main

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type ContainerInfo struct {
	ID     string
	Names  []string
	Image  string
	State  string
	Ports  []types.Port
	Status string
}

func main() {
	r := gin.Default()

	// Load templates from the templates directory
	r.LoadHTMLGlob("templates/*")

	// Define the handler for the containers page
	r.GET("/admin/containers", func(c *gin.Context) {
		containers, err := getDockerContainers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// hostname, err := os.Hostname()
		// Render the template and pass the list of containers
		c.HTML(http.StatusOK, "containers.tmpl", gin.H{
			"Containers": containers,
			// "Hostname":   hostname,
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/admin/containers")
	})

	// Start the server
	r.Run(":8082")
}

// getDockerContainers retrieves the list of Docker containers
func getDockerContainers() ([]ContainerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	var containerInfos []ContainerInfo
	for _, container := range containers {
		containerInfos = append(containerInfos, ContainerInfo{
			ID:     container.ID[:10], // Display only the first 10 characters of the container ID
			Names:  container.Names,
			Image:  container.Image,
			State:  container.State,
			Ports:  container.Ports,
			Status: container.Status,
		})
	}

	return containerInfos, nil
}
