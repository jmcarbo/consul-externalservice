package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	cesw "github.com/jmcarbo/consul-externalservice"
	"os"
	"os/signal"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "consul-externalservice"
	app.Usage = "manage consul external services!"
	app.Version = "0.0.3"
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
					go func() {
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
		{
			Name:      "export",
			ShortName: "e",
			Usage:     "export service definitions",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "export.yaml",
					Usage: "export file name",
				},
			},
			Action: func(c *cli.Context) {
				client := cesw.Connect()
        log.Infof("Exporting services to %s", c.String("file"))
        cesw.BackupExternalServicesToYAML(client, c.String("file"))
      },
    },
		{
			Name:      "import",
			ShortName: "i",
			Usage:     "import service definitions",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "export.yaml",
					Usage: "import file name",
				},
			},
			Action: func(c *cli.Context) {
				client := cesw.Connect()
        log.Infof("Importing services from %s", c.String("file"))
        cesw.RestoreExternalServicesFromYAML(client, c.String("file"))
      },
    },
	}
	app.Run(os.Args)
}
