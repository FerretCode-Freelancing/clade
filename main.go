package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/ferretcode-freelancing/clade/containers"
	mw "github.com/ferretcode-freelancing/clade/middleware"
	"github.com/ferretcode-freelancing/clade/registry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func setupLogger(homeDir string) {
	appDir := homeDir + "/.clade"
	appDir = filepath.FromSlash(appDir)

	logsDir := appDir + "/logs"
	logsDir = filepath.FromSlash(logsDir)

	if _, err := os.Stat(appDir); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}

		err := os.Mkdir(appDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		err = os.Mkdir(logsDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err := os.Create(filepath.FromSlash(
		logsDir + "/" + strconv.Itoa(int(time.Now().Unix())) + ".log"))
	if err != nil {
		log.Fatal(err)
	}

	// log.SetOutput(file)

	log.Println("Starting clade...")
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	setupLogger(homeDir)

	if os.Getenv("AUTH_SECRET") == "" {
		log.Fatal(errors.New("the auth secret for the http server must be set before running clade"))
	}

	registry.InitRegistry()

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	err = containers.Tick(cli, ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 5)

		for range ticker.C {
			err = containers.Tick(cli, ctx)
			if err != nil {
				log.Printf("Err performing orchestrator tick: %v\n", err)
			}
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.Auth())

	r.Route("/registry", func(r chi.Router) {
		r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
			request := registry.Request{}

			err := processRequest(r.Body, &request)
			handleError(w, err)

			err = registry.Add(request, homeDir)
			handleError(w, err)
		})

		r.Post("/remove", func(w http.ResponseWriter, r *http.Request) {
			request := registry.Request{}

			err := processRequest(r.Body, &request)
			handleError(w, err)

			err = registry.Remove(request, homeDir)
			handleError(w, err)
		})
	})

	r.Route("/containers", func(r chi.Router) {
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			request := containers.Request{}

			err := processRequest(r.Body, &request)
			handleError(w, err)

			err = containers.Create(request, cli, ctx)
			handleError(w, err)
		})

		r.Post("/delete", func(w http.ResponseWriter, r *http.Request) {
			request := containers.Request{}

			err := processRequest(r.Body, &request)
			handleError(w, err)

			err = containers.Delete(request, cli)
			handleError(w, err)
		})

		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			request := containers.Request{}

			err := processRequest(r.Body, &request)
			handleError(w, err)

			err = containers.Update(request, cli)
			handleError(w, err)
		})
	})

	if os.Getenv("CLADE_CERT_FILE") != "" && os.Getenv("CLADE_KEY_FILE") != "" {
		http.ListenAndServeTLS(":3002", os.Getenv("CLADE_CERT_FILE"), os.Getenv("CLADE_KEY_FILE"), r)
	} else {
		http.ListenAndServe(":3002", r)
	}
}

func handleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(200)
		w.Write([]byte("the operation was successful"))
	}
}

func processRequest(body io.ReadCloser, request interface{}) error {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &request); err != nil {
		log.Printf("Err parsing request: %v\n", err)
		return err
	}

	return nil
}
