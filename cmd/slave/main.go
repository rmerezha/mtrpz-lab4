package main

import (
	"flag"
	"github.com/rmerezha/mtrpz-lab4/listener"
	"github.com/rmerezha/mtrpz-lab4/runner"
	"log"
	"time"
)

var (
	masterUrl = flag.String("master", "", "master node url")
	host      = flag.String("host", "", "host node")
	interval  = flag.Duration("interval", 5*time.Second, "interval")
	token     = flag.String("token", "", "auth token")
)

func main() {
	flag.Parse()

	if *masterUrl == "" {
		log.Println("--master must be specified")
	}
	if *host == "" {
		log.Println("--host must be specified")
	}
	if *token == "" {
		log.Println("--token must be specified")
	}

	runner, err := runner.NewDockerRunner()
	if err != nil {
		log.Fatal(err)
	}
	store := listener.NewContainerStateStore()
	globalListener := listener.GlobalListener{
		Listeners: []listener.Listener{
			listener.NewPollingListener(*masterUrl, *host, runner, *interval, *token, store),
			listener.NewStateWatcherListener(*masterUrl, *host, runner, *interval, *token, store),
		},
	}

	stopCh := make(chan struct{})
	log.Println("slave node is starting")
	globalListener.Listen(stopCh)
}
