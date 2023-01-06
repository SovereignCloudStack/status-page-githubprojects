package server

import (
	"context"

	"github.com/SovereignCloudStack/status-page-openapi/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/shurcooL/githubv4"
)

func (s *ServerImplementation) GetImpacttypes(ctx echo.Context) error {
	var query struct {
		Node struct {
			ProjectV2 struct {
				Field struct {
					ProjectV2SingleSelectField struct {
						Options []struct {
							Name string
						}
					} `graphql:"... on ProjectV2SingleSelectField"`
				} `graphql:"field(name: \"Impact Type\")"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $projectid)"`
	}
	err := s.GithubV4Client.Query(
		context.Background(),
		&query,
		map[string]interface{}{
			"projectid": githubv4.ID(s.ProjectID),
		},
	)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(500)
	}
	impactTypes := []api.IncidentImpactType{}
	for _, phase := range query.Node.ProjectV2.Field.ProjectV2SingleSelectField.Options {
		impactTypes = append(impactTypes, phase.Name)
	}
	return ctx.JSON(200, impactTypes)
}
