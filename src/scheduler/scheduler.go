package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"scheduler/api"
	"scheduler/api/api_monitoring"
	"scheduler/api/api_peer"
	"scheduler/config"
	"scheduler/discovery"
	"scheduler/log"
	"scheduler/metrics"
	"scheduler/queue"
	"scheduler/scheduler"
	"strings"
	"sync"
)

import _ "net/http/pprof"

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	// init modules
	config.Start()
	scheduler.Start()
	discovery.Start()
	metrics.Start()

	go worker()
	go server()

	// Check if profiling should be enabled
	if strings.ToLower(os.Getenv(config.EnvProfiling)) == "true" {
		go pprof()
	}

	log.Log.Infof("Started p2p-fog scheduler v" + config.Version)

	wg.Wait()
}

func server() {
	log.Log.Debugf("Starting webserver thread")

	// init api
	router := mux.NewRouter()
	router.HandleFunc("/", api.Hello).Methods("GET")
	// OpenFaaS APIs
	router.HandleFunc("/system/functions", api.SystemFunctionsGet).Methods("GET")
	router.HandleFunc("/system/functions", api.SystemFunctionsPost).Methods("POST")
	router.HandleFunc("/system/functions", api.SystemFunctionsPut).Methods("PUT")
	router.HandleFunc("/system/functions", api.SystemFunctionsDelete).Methods("DELETE")
	router.HandleFunc("/system/function/{function}", api.SystemFunctionGet).Methods("GET")
	router.HandleFunc("/system/scale-function/{function}", api.SystemScaleFunctionPost).Methods("POST")
	router.HandleFunc("/function/{function}", api.FunctionPost).Methods("POST")
	router.HandleFunc("/function/{function}", api.FunctionGet).Methods("GET")
	// new APIs
	router.HandleFunc("/monitoring/load", api_monitoring.LoadGetLoad).Methods("GET")
	router.HandleFunc("/monitoring/scale-delay/{function}", api_monitoring.ScaleDelay).Methods("GET")
	router.HandleFunc("/peer/function/{function}", api_peer.FunctionExecute).Methods("POST")
	// prometheus
	router.Handle("/metrics", promhttp.Handler())
	// dev apis
	router.HandleFunc("/configuration", api.GetConfiguration).Methods("GET")
	router.HandleFunc("/configuration/scheduler", api.GetScheduler).Methods("GET")
	// TODO add auth check on configuration APIs
	// if config.Configuration.GetRunningEnvironment() == config.RunningEnvironmentDevelopment {
	router.HandleFunc("/configuration", api.SetConfiguration).Methods("POST")
	router.HandleFunc("/configuration/scheduler", api.SetScheduler).Methods("POST")
	// }

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.Configuration.GetListeningPort()),
		Handler: router,
	}

	log.Log.Infof("Started listening on %d", config.Configuration.GetListeningPort())
	err := server.ListenAndServe()

	log.Log.Fatalf("Error while starting server: %s", err)
	wg.Done()
}

func worker() {
	log.Log.Debugf("Starting queue worker thread")
	queue.Looper()
	wg.Done()
}

func pprof() {
	// pprof
	_ = http.ListenAndServe("0.0.0.0:16060", nil)
}
