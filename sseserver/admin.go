package sseserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type status struct {
	Node        string             `json:"node"`
	Status      string             `json:"status"`
	Reported    int64              `json:"reported_at"`
	Connections []connectionStatus `json:"connections"`
}

// add a statusReport to hub type
func (h *hub) statusReport() string {

	stat := status{
		Node:     fmt.Sprintf("%s-%s-%s", platform(), env(), dyno()),
		Status:   "OK",
		Reported: time.Now().Unix(),
	}

	stat.Connections = []connectionStatus{}
	for k := range h.connections {
		stat.Connections = append(stat.Connections, k.Status())
	}

	b, _ := json.MarshalIndent(stat, "", "  ")
	return string(b)
}

func platform() string {
	return "go"
}

func dyno() string {
	dyno := os.Getenv("DYNO")
	if dyno != "" {
		return dyno
	} else {
		return "dev.1"
	}
}

func env() string {
	env := os.Getenv("GO_ENV")
	if env != "" {
		return env
	} else {
		return "development"
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request, h *hub) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, h.statusReport())
	return
}
