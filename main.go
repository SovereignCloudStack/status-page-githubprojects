package main

import (
	"context"
	"flag"
	"os"

	"github.com/SovereignCloudStack/status-page-githubprojects/pkg/server"
	"github.com/SovereignCloudStack/status-page-openapi/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	addr := flag.String("addr", ":3000", "address to listen on")
	projectOwner := flag.String("github.project.owner", "joshmue", "user owning the project")
	projectOwnerIsOrg := flag.Bool("github.project.owner.is-org", false, "sets whether the owner of the github project is an org instead of an user")
	projectNumber := flag.Int64("github.project.number", 1, `project number, i. e. the "1" in https://github.com/users/$USER/projects/1`)
	lastPhase := flag.String("last-phase", "Done", "last phase of incidents")
	flag.Parse()

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	))
	server := &server.ServerImplementation{
		GithubV4Client:    githubv4.NewClient(httpClient),
		ProjectOwner:      *projectOwner,
		ProjectOwnerIsOrg: *projectOwnerIsOrg,
		ProjectNumber:     *projectNumber,
		LastPhase:         *lastPhase,
	}

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.Debugf("Obtaining Github Project ID...")
	if err := server.FillProjectID(); err != nil {
		e.Logger.Fatal(err)
	}
	e.Logger.Debugf("Ensuring Github Project configuration meets expectations...")
	if err := server.EnsureProjectConfiguration(); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Debugf("Registering handlers...")
	e.Use(middleware.Logger())
	api.RegisterHandlers(e, server)
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := api.GetSwagger()
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(500)
		}
		return c.JSON(200, swagger)
	})
	e.GET("/swagger/", serveSwagger)

	e.Logger.Debugf("Starting server...")
	e.Logger.Fatal(e.Start(*addr))
}
