package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev    bool `json:"dev" envconfig:"DEV" desc:"Development mode"`
	Server struct {
		ServiceName         string        `json:"service_name" envconfig:"SERVICE_NAME" default:"api" desc:"Service name"`
		GatewayAddr         string        `json:"public_addr" envconfig:"ADDR" default:":8080" desc:"Gateway server listen address"`
		GRPCAddr            string        `json:"grpc_addr" envconfig:"GRPC_ADDR" default:":8090" desc:"GRPC server listen address"`
		Timeout             time.Duration `json:"timeout" envconfig:"TIMEOUT" default:"5s" desc:"Operation timeout"`
		ShutdownTimeout     time.Duration `json:"shutdown_timeout" envconfig:"SHUTDOWN_TIMEOUT" default:"10s" desc:"Shutdown timeout"`
		InstanceID          string        `json:"instance_id" envconfig:"INSTANCE_ID" default:"" desc:"Instance ID"`
		OpenCensusAgentHost string        `json:"oc_agent_host" envconfig:"OC_AGENT_HOST" default:"" desc:"OpenCensus agent host"`
	} `json:"server" envconfig:"SERVER"`
}

func parseConfig() (*Config, error) {
	var c Config

	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}
