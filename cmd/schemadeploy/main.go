package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hirokazumiyaji/schemadeploy"
)

var (
	dryrun = flag.Bool("dryrun", false, "dryrun")
	dsn    = flag.String("dsn", os.Getenv("MYSQL_DSN"), "mysql database dsn")
	schema = flag.String("schema", "", "path to schema file")
)

func main() {
	flag.Parse()
	if dsn == nil || *dsn == "" {
		fmt.Println("dsn option is required")
		flag.PrintDefaults()
		return
	}
	if schema == nil || *schema == "" {
		fmt.Println("schema option is required")
		flag.PrintDefaults()
		return
	}
	if err := run(); err != nil {
		log.Fatal(err)
	} else {
		if !*dryrun {
			fmt.Println("success")
		}
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-sig:
			cancel()
			return
		}
	}()

	r := &schemadeploy.Runner{
		Dryrun: *dryrun,
		DSN:    *dsn,
		Schema: *schema,
	}
	return r.Run(ctx)
}
