package main

import (
	"github.com/Eldarrin/devops-migrator/pkg/api"
	"github.com/Eldarrin/devops-migrator/pkg/conf"
	"github.com/Eldarrin/devops-migrator/pkg/server"
	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	users = map[int]*user{}
	seq   = 1
)

//----------
// Handlers
//----------

func createUser(c echo.Context) error {
	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

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
