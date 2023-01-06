# SCS Status Page - GitHub Projects Backend

This repository contains an implementation of the [SCS Status Page API](https://github.com/SovereignCloudStack/status-page-openapi) backed by GitHub Projects.
This means that all state relevant to the status page is stored in a [GitHub Project](https://docs.github.com/en/issues/planning-and-tracking-with-projects/learning-about-projects/about-projects) and fetched via the GitHub GraphQL API.

## Mapping of attributes

| SCS Status Page API | GitHub |
| --- | --- |
| Component | Labels (`LA_***`) with "component:" prefix |
| Incident | Project Items (`PVTI_***`) of type "ISSUE" |
| Incident.phase | Project Item Field "Status" (Single select, predefined) |
| Incident.impactType | Project Item Field "Impact Type" (Single select) |
| Incident.beganAt | Project Item Field "Began At" (Text) |
| Incident.endedAt | Project Item Field "Ended At" (Text) |

On startup, the server verifies that fields are configured accordingly in GitHub. If this is not the case, it will not start.

Set `GITHUB_TOKEN` as environment variable with all required permissions (if using PAT's, use classic tokens, [as fine grained tokens are not yet supported with GraphQL](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#about-personal-access-tokens)), see `--help` for all other parameters.

Example invocation for debugging:

```
GITHUB_TOKEN=... go run *.go --github.project.number=1 --github.project.owner=$USER
```