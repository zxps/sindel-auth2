package main

import (
	cfg "auth2/config"
	"auth2/container/bundle"
	"auth2/container/prototype"
	pb "auth2/proto"
	"auth2/server"
	"auth2/utils"
	"auth2/validator"
	"flag"
	"fmt"
	"github.com/tkanos/gonfig"
	"google.golang.org/grpc"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

const ConfigFilename = "config.json"

func main() {
	configFile := flag.String("c", getDefaultConfigPath(), "Configuration file path")
	listenAddress := flag.String("l", "", "Listen server address")

	flag.Parse()

	if !utils.IsFileExists(*configFile) {
		logrus.Errorf("config file not found %s", *configFile)
		panic("config file not found")
	}

	logrus.Infof("using config file %s", *configFile)

	config := cfg.Config{}
	gonfig.GetConf(*configFile, &config)

	if err := validator.ValidateConfig(&config); err != nil {
		panic(err.Error())
	}

	config.ConfigPath = *configFile

	var services prototype.ContainerServices = *bundle.CreateServices(&config)

	container := prototype.NewContainer(&prototype.ContainerOptions{
		Config:   &config,
		Services: &services,
	})

	if listenAddress == nil {
		listenAddress = &config.ListenAddress
	}

	start(config.ListenAddress, container)
}

func start(address string, container *prototype.Container) {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		logrus.Fatal("Failed to listen ", address)
	}

	var authServer pb.Auth2Server = server.New(container)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuth2Server(grpcServer, authServer)

	grpcServer.Serve(lis)
}

func getDefaultConfigPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s/%s", path, ConfigFilename)
}

func validateConfig(config *cfg.Config) {
	if len(config.ListenAddress) < 1 {
		panic("server address not specified")
	}

	if len(config.DBConnection) < 1 {
		panic("database connection not specified")
	}

	if len(config.UsersTable) < 1 {
		panic("users table name not specified")
	}
}
