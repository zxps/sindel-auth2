package prototype

import (
	cfg "auth2/config"
)

type ContainerOptions struct {
	Config   *cfg.Config
	Services *ContainerServices
}

type Container struct {
	config   *cfg.Config
	services *ContainerServices
}

func NewContainer(opts *ContainerOptions) *Container {
	return &Container{
		config:   opts.Config,
		services: opts.Services,
	}
}

func (s *Container) Config() *cfg.Config {
	return s.config
}

func (s *Container) Services() *ContainerServices {
	return s.services
}
