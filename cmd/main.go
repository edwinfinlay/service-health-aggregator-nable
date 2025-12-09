package main

import (
	"context"
	healthAgg "github.com/edwinfinlay/service-health-aggregator-nable/internal/app/service-health-aggregator"
	"github.com/edwinfinlay/service-health-aggregator-nable/internal/app/service-health-aggregator/server"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := healthAgg.NewConfig()
	service := healthAgg.NewHealthService(cfg)

	healthServer := server.NewServer(ctx, service)

	if err := healthServer.Serve(); err != nil {
		log.Fatal(err)
	}
}
