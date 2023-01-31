package server

import (
	"context"
	"strings"

	"github.com/SovereignCloudStack/status-page-openapi/pkg/api"
	"github.com/labstack/echo/v4"
	"github.com/shurcooL/githubv4"
)

func (l *projectLabel) ToComponent(lastPhase string) api.Component {
	affectedBy := []api.Id{}
	for issue := range l.Issues.Nodes {
		for projectItem := range l.Issues.Nodes[issue].ProjectItems.Nodes {
			if l.Issues.Nodes[issue].ProjectItems.Nodes[projectItem].FieldValueByName.ProjectV2ItemFieldSingleSelectValue.Name == lastPhase {
				continue
			}
			affectedBy = append(affectedBy, l.Issues.Nodes[issue].ProjectItems.Nodes[projectItem].Id)
		}
	}
	displayName := strings.TrimPrefix(l.Name, "component:")
	if l.Description != "" {
		displayName = l.Description
	}

	return api.Component{
		AffectedBy:  affectedBy,
		DisplayName: displayName,
		Id:          l.Id,
		Labels:      map[string]string{}, // TODO
	}
}

type projectLabel struct {
	Id          string
	Name        string
	Description string
	Issues      struct {
		Nodes []struct {
			ProjectItems struct {
				Nodes []struct {
					Id               string
					FieldValueByName struct {
						ProjectV2ItemFieldSingleSelectValue struct {
							Name string
						} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
					} `graphql:"fieldValueByName(name:\"Status\")"`
				}
			} `graphql:"projectItems(first:10)"`
		}
	} `graphql:"issues(first:10)"`
}

func (s *ServerImplementation) GetComponent(ctx echo.Context, componentId string) error {
	var query struct {
		Node struct {
			Label projectLabel `graphql:"... on Label"`
		} `graphql:"node(id: $labelid)"`
	}
	err := s.GithubV4Client.Query(
		context.Background(),
		&query,
		map[string]interface{}{
			"labelid": githubv4.ID(componentId),
		},
	)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(500)
	}
	return ctx.JSON(200, query.Node.Label.ToComponent(s.LastPhase))
}
func (s *ServerImplementation) GetComponents(ctx echo.Context) error {
	var query struct {
		Node struct {
			ProjectV2 struct {
				Repositories struct {
					Nodes []struct {
						Labels struct {
							Nodes []projectLabel
						} `graphql:"labels(first: 10)"`
					}
				} `graphql:"repositories(first: 10)"`
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
	components := []api.Component{}
	for repo := range query.Node.ProjectV2.Repositories.Nodes {
		for label := range query.Node.ProjectV2.Repositories.Nodes[repo].Labels.Nodes {
			if !strings.HasPrefix(query.Node.ProjectV2.Repositories.Nodes[repo].Labels.Nodes[label].Name, "component:") {
				continue
			}
			components = append(components, query.Node.ProjectV2.Repositories.Nodes[repo].Labels.Nodes[label].ToComponent(s.LastPhase))
		}
	}
	return ctx.JSON(200, components)
}
