package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/blaketigges/gitlab-ci-pipelines-exporter/pkg/schemas"
	"github.com/stretchr/testify/assert"
)

func TestListRefPipelineJobs(t *testing.T) {
	ctx, mux, server, c := getMockedClient()
	defer server.Close()

	ref := schemas.Ref{
		Project: schemas.NewProject("foo"),
		Name:    "yay",
	}

	// Test with no most recent pipeline defined
	jobs, err := c.ListRefPipelineJobs(ctx, ref)
	assert.NoError(t, err)
	assert.Len(t, jobs, 0)

	mux.HandleFunc("/api/v4/projects/foo/pipelines/1/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":10}]`)
		})

	mux.HandleFunc("/api/v4/projects/11/pipelines/2/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":20}]`)
		})

	mux.HandleFunc("/api/v4/projects/12/pipelines/3/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":30}]`)
		})

	mux.HandleFunc("/api/v4/projects/foo/pipelines/1/bridges",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":1,"downstream_pipeline":{"id":2, "project_id": 11}}]`)
		})

	mux.HandleFunc("/api/v4/projects/11/pipelines/2/bridges",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"id":1,"downstream_pipeline":{"id":3, "project_id": 12}}]`)
		})

	mux.HandleFunc("/api/v4/projects/12/pipelines/3/bridges",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[]`)
		})

	ref.LatestPipeline = schemas.Pipeline{
		ID: 1,
	}

	jobs, err = c.ListRefPipelineJobs(ctx, ref)
	assert.NoError(t, err)
	assert.Equal(t, []schemas.Job{
		{ID: 10},
		{ID: 20},
		{ID: 30},
	}, jobs)

	// Test invalid project id
	ref.Project.Name = "bar"
	_, err = c.ListRefPipelineJobs(ctx, ref)
	assert.Error(t, err)
}

func TestListPipelineJobs(t *testing.T) {
	ctx, mux, server, c := getMockedClient()
	defer server.Close()

	mux.HandleFunc("/api/v4/projects/foo/pipelines/1/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			expectedQueryParams := url.Values{
				"page":     []string{"1"},
				"per_page": []string{"100"},
			}
			assert.Equal(t, expectedQueryParams, r.URL.Query())
			fmt.Fprint(w, `[{"id":1},{"id":2}]`)
		})

	mux.HandleFunc("/api/v4/projects/bar/pipelines/1/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

	jobs, err := c.ListPipelineJobs(ctx, "foo", 1)
	assert.NoError(t, err)
	assert.Len(t, jobs, 2)

	// Test invalid project id
	_, err = c.ListPipelineJobs(ctx, "bar", 1)
	assert.Error(t, err)
}

func TestListPipelineBridges(t *testing.T) {
	ctx, mux, server, c := getMockedClient()
	defer server.Close()

	mux.HandleFunc("/api/v4/projects/foo/pipelines/1/bridges",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			expectedQueryParams := url.Values{
				"page":     []string{"1"},
				"per_page": []string{"100"},
			}
			assert.Equal(t, expectedQueryParams, r.URL.Query())
			fmt.Fprint(w, `[{"id":1,"pipeline":{"id":100}}]`)
		})

	mux.HandleFunc("/api/v4/projects/bar/pipelines/1/bridges",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

	bridges, err := c.ListPipelineBridges(ctx, "foo", 1)
	assert.NoError(t, err)
	assert.Len(t, bridges, 1)

	// Test invalid project id
	_, err = c.ListPipelineBridges(ctx, "bar", 1)
	assert.Error(t, err)
}

func TestListRefMostRecentJobs(t *testing.T) {
	ctx, mux, server, c := getMockedClient()
	defer server.Close()

	ref := schemas.Ref{
		Project: schemas.NewProject("foo"),
		Name:    "yay",
	}

	jobs, err := c.ListRefMostRecentJobs(ctx, ref)
	assert.NoError(t, err)
	assert.Len(t, jobs, 0)

	mux.HandleFunc("/api/v4/projects/foo/jobs",
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			expectedQueryParams := url.Values{
				"page":     []string{"1"},
				"per_page": []string{"100"},
			}
			assert.Equal(t, expectedQueryParams, r.URL.Query())
			fmt.Fprint(w, `[{"id":3,"name":"foo","ref":"yay"},{"id":4,"name":"bar","ref":"yay"}]`)
		})

	mux.HandleFunc(fmt.Sprintf("/api/v4/projects/bar/jobs"),
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

	ref.LatestJobs = schemas.Jobs{
		"foo": {
			ID:   1,
			Name: "foo",
		},
		"bar": {
			ID:   2,
			Name: "bar",
		},
	}

	jobs, err = c.ListRefMostRecentJobs(ctx, ref)
	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, 3, jobs[0].ID)
	assert.Equal(t, 4, jobs[1].ID)

	ref.LatestJobs["baz"] = schemas.Job{
		ID:   5,
		Name: "baz",
	}

	jobs, err = c.ListRefMostRecentJobs(ctx, ref)
	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, 3, jobs[0].ID)
	assert.Equal(t, 4, jobs[1].ID)

	// Test invalid project id
	ref.Project.Name = "bar"
	_, err = c.ListRefMostRecentJobs(ctx, ref)
	assert.Error(t, err)
}
