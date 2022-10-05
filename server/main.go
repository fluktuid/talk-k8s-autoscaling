package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/fluktuid/talk-k8s-autoscaling/config"
	"github.com/fluktuid/talk-k8s-autoscaling/metrics"
	"github.com/fluktuid/talk-k8s-autoscaling/util"
)

func main() {
	var c config.Server
	err := envconfig.Process("", &c)
	log.WithField("config", c).Info("Import Config")
	if err != nil {
		log.Fatal(err)
	}
	metrics.InitAsync()
	hostname := util.Hostname()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello there! 👋")
		defer metrics.ServerRecordRequest()
		time.Sleep(time.Duration(c.DelayResponse) * time.Millisecond)
		fmt.Fprintf(w, "Hello there!\nI'm %s.", hostname)
	})

	log.Info("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
