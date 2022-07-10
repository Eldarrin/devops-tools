package server

import (
	"context"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/Eldarrin/devops-tools/pkg/api"
	"github.com/Eldarrin/devops-tools/pkg/conf"
	migGithub "github.com/Eldarrin/devops-tools/pkg/github"
	"github.com/labstack/echo/v4"
)

const (
	// DefaultMigrationID TODO NOT THIS
	DefaultMigrationID = "a4a777ff-fd47-42ab-84b4-1cca19a51f8f"
)

// Server represents all server handlers.
type Server struct {
	cfg *conf.Config
}

// NewServer new api server
func NewServer(cfg *conf.Config) (*Server, error) {
	return &Server{cfg: cfg}, nil
}

// Migration Get a list of migration. (GET /migrate)
func (sv *Server) Migration(ctx echo.Context, params api.MigrationParams) error {

	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	return ctx.JSON(http.StatusOK, &api.Migration{Id: DefaultMigrationID})
}

// NewMigration Create a migration. (PUT /migrate)
func (sv *Server) NewMigration(ctx echo.Context) error {

	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	newMig := new(api.NewMigration)
	if err := ctx.Bind(newMig); err != nil {
		return err
	}

	/*
		resMig, err := sv.stores.Customers.Create(ctx.Request().Context(), newCust)
		if err != nil {
			if err == store.ErrCustomerNameAlreadyExists {
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			}
			return err
		}
	*/

	return ctx.JSON(http.StatusCreated, newMig)
}

// GetMigration (GET /migration/{id})
func (sv *Server) GetMigration(ctx echo.Context, id string) error {

	// Validate access token.
	//
	// 🚨 SECURITY: It's important we check for the correct scopes to know what this token
	// is allowed to do.
	if !userHasAccess(ctx) {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope")
	}

	/*
		resMig, err := sv.stores.Customers.GetByID(ctx.Request().Context(), id)
		if err != nil {
			if _, ok := err.(*store.MigrationNotFoundError); ok {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return err
		}
	*/

	return ctx.JSON(http.StatusOK, &api.Migration{Id: DefaultMigrationID})
}

// ExecuteMigration Runs a migration. (POST /migrate)
func (sv *Server) ExecuteMigration(ctx echo.Context) error {
	ctxB := context.Background()

	// check for existing migration and create if none
	log.Info().Msg("executing migration")

	// pre-check IDP groups exist in GitHub or fail fast
	groupName := ctx.Param("contributorgroup")

	if !migGithub.DoesIDPGroupExist(groupName, ctxB) {
		log.Warn().Msg("Group Doesn't Exist")
		return ctx.JSON(http.StatusNotModified, &api.Migration{Id: "cool"})
	}

	// execute repo copy from ADO to GitHub
	migGithub.StartRepoMigration(ctxB)

	// convert private repository to internal repository
	migGithub.ChangeVisibility(ctxB)

	// create contributor team if doesn't exist
	//       (assign IDP AAD Group, fail if doesn't exist)
	//       (trim GBL ROL IT DevOps, convert to lcase, dash for spaces
	// add repo to team - set as write permission
	// create manager team if doesn't exist
	//       (assign IDP AAD Group, fail if doesn't exist)
	//       (trim GBL ROL IT DevOps, convert to lcase, dash for spaces
	// add repo to team - set as maintain
	// create dummy actions to match pipelines
	// get branch polices for the repo if defined
	// create empty branch policies in github for each defined branch policy

	return ctx.JSON(http.StatusOK, &api.Migration{Id: DefaultMigrationID})
}

func (sv *Server) SetupProject(ctx echo.Context) error {

	/*
		setup
			SonaType / Artifactory
			SonarQube
			ADO
			GitHub
	*/

	return ctx.JSON(http.StatusOK, &api.Migration{Id: DefaultMigrationID})
}

func userHasAccess(ctx echo.Context) bool {

	/*
		user, err := auth.LoadUserFromContext(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to load user from context")
			return false
		}

		scopes, err := auth.LoadOperationScopesFromContext(ctx, auth.OpenIDScopes)
		if err != nil {
			log.Error().Err(err).Msg("failed to load scopes from context")
			return false
		}

		log.Info().Strs("Scopes", scopes).Object("User", &user).Msg("Scopes check")

		return user.HasScope(scopes)
	*/

	return true
}
