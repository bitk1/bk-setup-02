package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/exec"
)

var templates = template.Must(template.ParseFiles("index.html"))

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/add", addFileHandler)
    http.HandleFunc("/cat", catFileHandler)
    http.HandleFunc("/peers", peersHandler)

    fmt.Println("Server started at :8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        templates.ExecuteTemplate(w, "index.html", nil)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func addFileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        file, _, err := r.FormFile("file")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        tempFile, err := ioutil.TempFile("", "upload-*.tmp")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer os.Remove(tempFile.Name())

        fileBytes, err := ioutil.ReadAll(file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tempFile.Write(fileBytes)

        cmd := exec.Command("ipfs", "add", tempFile.Name())
        output, err := cmd.CombinedOutput()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"result": string(output)})
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func catFileHandler(w http.ResponseWriter, r *http.Request) {
    hash := r.URL.Query().Get("hash")
    if hash == "" {
        http.Error(w, "Missing hash parameter", http.StatusBadRequest)
        return
    }

    cmd := exec.Command("ipfs", "cat", hash)
    output, err := cmd.CombinedOutput()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"result": string(output)})
}

func peersHandler(w http.ResponseWriter, r *http.Request) {
    cmd := exec.Command("ipfs", "swarm", "peers")
    output, err := cmd.CombinedOutput()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"result": string(output)})
}

