package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"os"
  cesw "github.com/jmcarbo/consul-externalservice"
  "os/signal"
  "time"
)

func main() {
	app := cli.NewApp()
	app.Name = "consul-externalservice"
	app.Usage = "manage consul external services!"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "consul-externalservice version",
			Action: func(c *cli.Context) {
				fmt.Println(app.Version)
			},
		},
		{
			Name:      "start",
			ShortName: "s",
			Usage:     "start service watcher",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "node",
					Value: "node1",
					Usage: "node name",
				},
      },
			Action: func(c *cli.Context) {
				client := cesw.Connect()
				watcher := cesw.NewExternalServiceWatcher(client, c.String("node"))
				if watcher != nil {
					log.Printf("Starting external service watcher for node %s ...\n", c.String("node"))
          stopCh := make(chan struct{})
          go func(){
TRY_LEADERSHIP:
					  watcher.Run()
            if watcher.IsLeader() {
              log.Info("I am the leader now ...")
            }
WAIT_FOR_EVENT:
            select {
              case <-stopCh:
                watcher.Destroy()
                return
              case <-time.After(10 * time.Second):
                if watcher.IsLeader() {
                  log.Info("I am still the leader ...")
                  goto WAIT_FOR_EVENT
                } else {
                  log.Info("Trying to be leader ...")
                  goto TRY_LEADERSHIP
                }
            }
          }()

          // Wait for termination
          signalCh := make(chan os.Signal, 1)
          signal.Notify(signalCh, os.Interrupt, os.Kill)
          select {
            case <-signalCh:
              log.Warn("Received signal, stopping service watch ...")
              close(stopCh)
          }
				} else {
          log.Error("Error starting external service watcher. Check consul agent is running on localhost:8500. Exiting ...")
        }
			},
		},
	}
	app.Run(os.Args)
}
