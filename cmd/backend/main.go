package main

import (
	"github.com/Eldarrin/devops-tools/pkg/api"
	"github.com/Eldarrin/devops-tools/pkg/conf"
	"github.com/Eldarrin/devops-tools/pkg/server"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

func main() {

	// loads configuration from env and configures logger
	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	svr, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bind api")
	}

	e := echo.New()
	// shut up
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	// add a version to the api
	g := e.Group("/v1")

	api.RegisterHandlers(g, svr)

	log.Info().Str("addr", cfg.Addr).Msg("starting http listener")
	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")
}
