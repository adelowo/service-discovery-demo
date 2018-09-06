package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/adelowo/service-discovery-demo/pkg/registry"
	"github.com/hashicorp/consul/api"
)

func main() {

	var discoveryURL = flag.String("discovery", "127.0.0.1:8500", "Consul service discovery url")
	var httpPort = flag.String("http", ":3000", "Port to run HTTP service at")

	flag.Parse()

	reg, err := registry.New(*discoveryURL)
	if err != nil {
		log.Fatalf("an error occurred while bootstrapping service discovery... %v", err)
	}

	var healthURL string

	ip, err := registry.IPAddr()
	if err != nil {
		log.Fatalf("could not determine IP address to register this service with... %v", err)
	}

	healthURL = "http://" + ip.String() + *httpPort + "/health"

	pp, err := strconv.Atoi((*httpPort)[1:]) // get rid of the ":" port
	if err != nil {
		log.Fatalf("could not discover port to register with consul.. %v", err)
	}

	svc := &api.AgentServiceRegistration{
		Name:    "cool_app",
		Address: ip.String(),
		Port:    pp,
		Tags:    []string{"urlprefix-/oops"},
		Check: &api.AgentServiceCheck{
			TLSSkipVerify: true,
			Method:        "GET",
			Timeout:       "20s",
			Interval:      "1m",
			HTTP:          healthURL,
			Name:          "HTTP check for cool app",
		},
	}

	id, err := reg.RegisterService(svc)
	if err != nil {
		log.Fatalf("Could not register service in consul... %v", err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		fmt.Println("Here")
		w.Write([]byte("home page"))
	})

	if err := http.ListenAndServe(*httpPort, nil); err != nil {
		reg.DeRegister(id)
	}
}
