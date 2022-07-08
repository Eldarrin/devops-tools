package github

import (
	"context"
	"flag"
	ghClient "github.com/google/go-github/v45/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

// CODE FROM MY GITFARM PROJECT

const (
	path    = "/jobs"
	health  = "/health"
	ready   = "/ready"
	workDir = "/runner"
	// Personal Access Token created in GitHub that allows us to make calls into GitHub.

)

type JobInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CallingURL   string `json:"calling_url"`
	Labels       string `json:"labels"`
	RepoName     string `json:"repo_name"`
	RunnerGroup  string `json:"runner_group"`
	Owner        string `json:"owner"`
	Organization string `json:"organization"`
}

var available = true

func getRunnerToken(ctx context.Context, jobInfo *JobInfo) string {
	personalAccessToken := os.Getenv(personalAccessTokenKey)
	if personalAccessToken == "" {
		log.Print("Unauthorized: No token present")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: personalAccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := ghClient.NewClient(tc)

	runnerToken, _, err := client.Actions.CreateRegistrationToken(ctx, jobInfo.Owner, jobInfo.RepoName)
	if err != nil {
		log.Print(err)
	}

	if jobInfo.Organization != "" {
		runnerToken, _, err = client.Actions.CreateOrganizationRegistrationToken(ctx, jobInfo.Organization)
		if err != nil {
			log.Print(err)
		}
	}

	runnerTokenValue := *runnerToken.Token

	return runnerTokenValue
}

func configureRunner(ctx context.Context, jobInfo *JobInfo) {
	githubUrl := "https://github.com/" + jobInfo.Owner + "/" + jobInfo.Name

	// do enterprise stuff
	if !strings.Contains(jobInfo.CallingURL, "api.github.com") {
		githubUrl = "this will break"
	}

	// clear and mash up labels
	// strip non-essentials
	labels := jobInfo.Labels

	runnerName := labels

	configApp := workDir + "/config.sh"

	// do organisation stuff
	if jobInfo.Organization != "" {
		githubUrl = "https://github.com/" + jobInfo.Organization + "/" + jobInfo.Name
	}

	runnerTokenValue := getRunnerToken(ctx, jobInfo)
	args := []string{configApp,
		"--unattended",
		"--replace",
		"--name", runnerName,
		"--url", githubUrl,
		"--token", runnerTokenValue,
		"--labels", labels,
		"--work", workDir,
		"--ephemeral",
		"--disableupdate"}

	if jobInfo.RunnerGroup != "" {
		args = append(args, "--runnergroups", jobInfo.RunnerGroup)
	}

	cmdConfig := &exec.Cmd{
		Path:   configApp,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	log.Print(cmdConfig.String())

	if err := cmdConfig.Run(); err != nil {
		log.Print(err)
	}
}

func HandleWorkflowJob(ctx context.Context, jobInfo *JobInfo, ch chan<- string) {
	log.Print("Handling Workflow Job")
	log.Print(ch)

	configureRunner(ctx, jobInfo)

	cmdRun := &exec.Cmd{
		Path:   workDir + "/run.sh",
		Args:   []string{workDir + "/run.sh"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	log.Print(ch)
	log.Print(cmdRun.String())

	if err := cmdRun.Run(); err != nil {
		log.Print(err)
	}
	log.Print(ch)
	available = true

}

func newJob(name string) *JobInfo {
	j := JobInfo{Name: name}
	j.ID = 1
	j.Labels = "main"
	j.CallingURL = "https://api.github.com"
	j.Owner = "eldarrin"
	j.RepoName = "knative-gitfarm"
	return &j
}

func handler(w http.ResponseWriter, _ *http.Request) {
	log.Print("in handler")

	ctx := context.Background()

	ch := make(chan string)

	if available {
		// block this call so knative thinks its processing and doesn't kill mid-job
		available = false
		HandleWorkflowJob(ctx, newJob("knative-gitfarm"), ch)
	} else {
		// accept no more requests so it spawns a new agent
		w.WriteHeader(503)
		w.Write([]byte("Server is active, send it somewhere else"))
	}

}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	if available {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("KO"))
	}
}

func main() {
	flag.Parse()
	log.Print("gitrunner started.")
}
