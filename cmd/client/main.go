package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/Eldarrin/devops-tools/pkg/api"
	"github.com/Eldarrin/devops-tools/pkg/oidc"
	"github.com/alecthomas/kingpin"
	"github.com/motemen/go-loghttp"
	"github.com/rs/zerolog/log"
)

var (
	app   = kingpin.New("client", "A command-line client for migrator.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	endpoint = app.Flag("endpoint", "The endpoint address of the migrator api.").Required().String()

	//auth = app.Flag("auth", "Use client credentials to authenticate to the api.").Required().Enum("clientcredentials")

	// required if you are using auth type "clientcredentials"
	clientID     = app.Flag("client-id", "oauth2 client id used with openid.").Envar("OAUTH_CLIENT_ID").String()
	clientSecret = app.Flag("client-secret", "oauth2 client secret used with openid.").Envar("OAUTH_CLIENT_SECRET").String()

	// openid authentication server
	authServer = app.Flag("openid-server", "OpenID authentication server URL.").Envar("OPENID_PROVIDER_URL").String()

	listProjects = app.Command("projects", "List projects.")

	version = "unknown"
)

func main() {
	kingpin.Version(version)

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *debug {
		http.DefaultTransport = &loghttp.Transport{
			Transport: http.DefaultTransport,
			LogRequest: func(req *http.Request) {
				log.Info().Msgf("[%p] %s %s", req, req.Method, req.URL)
				data, _ := httputil.DumpRequest(req, true)
				log.Info().Msgf(string(data))
			},
			LogResponse: func(resp *http.Response) {
				log.Info().Msgf("[%p] %d %s", resp.Request, resp.StatusCode, resp.Request.URL)
			},
		}
	}

	sess, err := oidc.New(context.TODO(), oidc.SessionConfig{
		ProviderURL:  *authServer,
		ClientID:     *clientID,
		ClientSecret: *clientSecret,
		Scopes:       []string{"migrator/project.read", "migrator/project.write"},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to oidc endoint")
	}

	switch cmd {
	case listProjects.FullCommand():

		client := &api.Client{Server: *endpoint, Client: sess.Client(context.TODO())}

		res, err := client.Migration(context.TODO(), &api.MigrationParams{})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to list migrations")
		}

		migrationRes, err := api.ParseMigrationResponse(res)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to list projects")
		}

		log.Info().Str("status", migrationRes.Status()).Msg("list projects")
		if migrationRes.StatusCode() != 200 {
			log.Fatal().Err(err).Msg("failed to list projects")
		}

	}
}
