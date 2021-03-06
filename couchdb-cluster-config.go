package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/gesellix/couchdb-cluster-config/v17/pkg"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "CouchDB Cluster Config"
	app.Usage = ""
	app.Description = "Setup a CouchDB 2.x cluster"
	app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	app.Commands = []*cli.Command{
		{
			Name:  "setup",
			Usage: "perform a cluster setup",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:  "nodes",
					Usage: "list of node ip addresses to participate in the CouchDB cluster",
				},
				&cli.DurationFlag{
					Name:  "delay",
					Usage: "time to wait before the cluster setup will be started",
					Value: 5 * time.Second,
				},
				&cli.DurationFlag{
					Name:  "timeout",
					Usage: "time until all nodes need to be available",
					Value: 20 * time.Second,
				},
				&cli.StringFlag{
					Name:  "username",
					Usage: "admin username - admin will be created, if missing",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "admin password",
				},
			},
			Action: func(c *cli.Context) error {
				nodes := c.StringSlice("nodes")
				if len(nodes) == 0 {
					return fmt.Errorf("please pass a list of node ip addresses")
				}

				if c.String("username") == "" {
					return fmt.Errorf("please provide an admin username")
				}
				if c.String("password") == "" {
					return fmt.Errorf("please provide an admin password")
				}

				ips := cluster_config.ToIpAddresses(nodes)
				delay := c.Duration("delay")
				timeout := c.Duration("timeout")

				fmt.Printf("Going to setup the following nodes as cluster, delayed by %fs\n%v\n", delay.Seconds(), ips)
				return cluster_config.SetupClusterNodes(
					cluster_config.ClusterSetupConfig{
						IpAddresses: ips,
						Delay:       delay,
						Timeout:     timeout,
					},
					cluster_config.BasicAuth{
						Username: c.String("username"),
						Password: c.String("password")},
					c.Bool("insecure"))
			},
		},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "insecure",
			Usage: "ignore server certificate if using https",
			Value: true,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
