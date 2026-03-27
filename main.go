package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const allowedDir = "./safe-files"

func main() {
	http.HandleFunc("/readfile", readFileHandler)
	http.HandleFunc("/exec", execHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

const safeDir = "/app/safe-files"

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("file")

	absPath, err := filepath.Abs(filepath.Join(safeDir, path))
	if err != nil || !strings.HasPrefix(absPath, safeDir) {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}
	w.Write(data)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		http.Error(w, "Command failed", 500)
		return
	}
	w.Write(out)
}
