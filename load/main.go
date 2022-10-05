package main

import (
	"math"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/fluktuid/talk-k8s-autoscaling/config"
	"github.com/fluktuid/talk-k8s-autoscaling/metrics"
)

func getConfig() config.Load {
	var c config.Load
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/load")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&c)
	log.WithField("conf", c).Info("LoadConfig")
	if err != nil {
		log.Fatal(err)
	}
	return c
}
func initCallers(parallel, waitMillisBetween int, callUrl string) []chan struct{} {
	ch := make([]chan struct{}, parallel)
	for i := range ch {
		ch[i] = make(chan struct{})
	}

	for i := 0; i < parallel; i++ {
		log.Debug("creating caller func")
		wait := 1000 / parallel
		go func(ch <-chan struct{}, parallel, waitMillisBetween int, callUrl string) {
			time.Sleep(time.Duration(wait) * time.Millisecond)
			callerFunc(ch, callUrl, waitMillisBetween)
		}(ch[i], parallel, waitMillisBetween, callUrl)
	}
	return ch
}

func main() {
	c := getConfig()

	channel := initCallers(c.ParallelRequests, c.WaitMillisBetween, c.CallURL)

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.WithField("name", e.Name).Info("Config file changed")
		c = getConfig()
		go func() {
			wait := 1000 / len(channel)
			for _, ch := range channel {
				close(ch)
				time.Sleep(time.Duration(wait) * time.Millisecond)
			}
		}()
		channel = initCallers(c.ParallelRequests, c.WaitMillisBetween, c.CallURL)
	})
	viper.WatchConfig()

	metrics.Init()
}

func callerFunc(c <-chan struct{}, url string, waitMillisBetween int) {

	call := func() (elapsedMs int64, success bool) {
		log.Debug("calling")
		start := time.Now()
		metrics.LoadRecordRequest()

		client := http.Client{
			Timeout: time.Duration(math.Max(float64(waitMillisBetween)-5, 100)) * time.Millisecond,
		}
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := client.Do(req)
		elapsedMs = time.Since(start).Milliseconds()
		if err != nil {
			log.WithField("err", err).Warnf("something failed calling url %s", url)
			return elapsedMs, false
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			log.WithField("status", resp.StatusCode).Warnf("something failed calling url %s", url)
			return elapsedMs, false
		}
		return elapsedMs, true
	}

	for {
		select {
		case <-c:
			return
		default:
			elapsedMs, success := call()
			metrics.LoadRecordResponse(success, elapsedMs)

			time.Sleep(time.Duration(waitMillisBetween-int(elapsedMs)) * time.Millisecond)
		}
	}
}
