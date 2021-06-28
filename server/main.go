package server

import (
	cfg "auth2/config"
	"auth2/container/bundle"
	"auth2/container/prototype"
	pb "auth2/proto"
	"auth2/utils"
	"auth2/validator"
	"context"
	"fmt"
	"github.com/tkanos/gonfig"
	"runtime"
)

func (s *AuthServer) Restart(ctx context.Context, in *pb.EmptyRequest) (*pb.EmptyResponse, error) {
	configPath := s.getConfig().ConfigPath

	config := cfg.Config{}
	gonfig.GetConf(configPath, &config)

	config.ConfigPath = configPath

	if err := validator.ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation error (%s)", err.Error())
	}

	var services prototype.ContainerServices = *bundle.CreateServices(&config)

	container := prototype.NewContainer(&prototype.ContainerOptions{
		Config:   &config,
		Services: &services,
	})

	s.reloadContainer(container)

	return &pb.EmptyResponse{}, nil
}

func (s *AuthServer) Stats(ctx context.Context, in *pb.EmptyRequest) (*pb.StatsResponse, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var response pb.StatsResponse

	response.Alloc = utils.FormatBytesAsText(m.Alloc)
	response.TotalAlloc = utils.FormatBytesAsText(m.TotalAlloc)
	response.Sys = utils.FormatBytesAsText(m.Sys)
	response.NumGc = utils.FormatBytesAsText(uint64(m.NumGC))

	return &response, nil
}
