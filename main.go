package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	// http.HandleFunc("/containers", getContainers)
	// http.Handle("/", http.FileServer(http.Dir("./static")))
	// log.Fatal(http.ListenAndServe(":8081", nil))
	mux := http.NewServeMux()
	mux.HandleFunc("/containers", getContainers)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// Wrap the mux with the CORS middleware
	handler := handleCORS(mux)

	http.ListenAndServe(":8081", handler)
}

// Middleware to handle CORS
func handleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getContainers(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		// fmt.Println(" err.Error() +>>>>>", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		// fmt.Println("+>>>>>", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		var containerDetails []map[string]interface{}
		details := map[string]interface{}{
			"StatusCode": http.StatusInternalServerError,
			"Message":    "Docker is not UP.",
		}
		containerDetails = append(containerDetails, details)
		jsonResponse, _ := json.Marshal(containerDetails)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
		return
	}

	// Create a slice to hold container details
	var containerDetails []map[string]interface{}

	for _, container := range containers {
		// Collecting container details in a map
		ports := []map[string]interface{}{}
		for _, port := range container.Ports {
			ports = append(ports, map[string]interface{}{
				"PrivatePort": port.PrivatePort,
				"PublicPort":  port.PublicPort,
			})
		}
		details := map[string]interface{}{
			"ID":         container.ID[:12],
			"Image":      container.Image,
			"Names":      container.Names,
			"Status":     container.Status,
			"State":      container.State,
			"Ports":      ports,
			"Labels":     container.Labels,
			"Created":    container.Created,
			"StatusCode": 200,
			"Message":    "Container List.",
		}
		containerDetails = append(containerDetails, details)
	}

	// Convert the container details slice to JSON
	jsonResponse, err := json.Marshal(containerDetails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set content type to application/json and write the responsegit
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
