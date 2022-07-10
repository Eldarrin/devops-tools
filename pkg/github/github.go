package github

import (
	"context"
	ghClient "github.com/google/go-github/v45/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

const (
	personalAccessTokenKey = "GH_PAT"
)

func getClient(ctx context.Context) (*ghClient.Client, error) {
	personalAccessToken := os.Getenv(personalAccessTokenKey)
	if personalAccessToken == "" {
		log.Error().Msg("Unauthorized: No token present")
		return nil, http.ErrServerClosed
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: personalAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	return ghClient.NewClient(tc), nil
}

func DoesIDPGroupExist(name string, ctx context.Context) bool {

	client, err := getClient(ctx)
	if err != nil {
		log.Error().Msg("Err from getclient")
		return false
	}

	options := ghClient.ListExternalGroupsOptions{DisplayName: &name}

	idpList, _, err := client.Teams.ListExternalGroups(ctx, "Eldarrin", &options)
	if err != nil {
		log.Error().Msg("is an error")
		return false
	}

	if len(idpList.Groups) > 0 {
		return true
	} else {
		log.Warn().Msg("Group doesn't exist")
	}

	return false
}

func StartRepoMigration(ctx context.Context) bool {

	client, err := getClient(ctx)
	if err != nil {
		log.Error().Msg("Err from getclient")
		return false
	}

	repo, _, err := client.Repositories.Get(ctx, "owner", "repo")
	if err != nil {
		log.Error().Msg("get repo has an error")
		return false
	}
	if repo != nil {
		log.Warn().Msg("repository exists, halting")
		return false
	}

	var repos []string
	repos[0] = "repo"
	status, _, err := client.Migrations.StartMigration(ctx, "owner", repos, nil)
	if err != nil {
		log.Error().Msg("migrate repo has an error")
		return false
	}
	log.Info().Msg(status.String())

	return true
}

func ChangeVisibility(ctx context.Context) bool {

	client, err := getClient(ctx)
	if err != nil {
		log.Error().Msg("Err from getclient")
		return false
	}

	repo, _, err := client.Repositories.Get(ctx, "owner", "repo")
	if err != nil {
		log.Error().Msg("get repo has an error")
		return false
	}
	if repo != nil {
		log.Warn().Msg("repository exists, halting")
		return false
	}

	visibility := "internal"

	repo.Visibility = &visibility

	_, _, err = client.Repositories.Edit(ctx, "owner", "name", repo)
	if err != nil {
		log.Error().Msg("change repo visibility has an error")
		return false
	}

	return true
}

func CreateRepo() {}

func CreateTeam() {}
