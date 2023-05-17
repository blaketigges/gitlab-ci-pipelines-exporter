package controller

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/blaketigges/gitlab-ci-pipelines-exporter/pkg/config"
	"github.com/blaketigges/gitlab-ci-pipelines-exporter/pkg/schemas"
	"github.com/stretchr/testify/assert"
)

func TestPullProjectsFromWildcard(t *testing.T) {
	ctx, c, mux, srv := newTestController(config.Config{})
	defer srv.Close()

	mux.HandleFunc("/api/v4/projects",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":1,"path_with_namespace":"foo","jobs_enabled":false},{"id":2,"path_with_namespace":"bar","jobs_enabled":true}]`)
		})

	w := config.NewWildcard()
	assert.NoError(t, c.PullProjectsFromWildcard(ctx, w))

	projects, _ := c.Store.Projects(ctx)
	p1 := schemas.NewProject("bar")

	expectedProjects := schemas.Projects{
		p1.Key(): p1,
	}
	assert.Equal(t, expectedProjects, projects)
}
