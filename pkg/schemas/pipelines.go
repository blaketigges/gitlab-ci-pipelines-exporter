package schemas

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	goGitlab "github.com/xanzy/go-gitlab"
)

// Pipeline ..
type Pipeline struct {
	ID                    int
	Coverage              float64
	Timestamp             float64
	StartTime             float64
	DurationSeconds       float64
	QueuedDurationSeconds float64
	Status                string
	Variables             string
	TestReport            TestReport
}

// TestReport ..
type TestReport struct {
	TotalTime    float64
	TotalCount   int
	SuccessCount int
	FailedCount  int
	SkippedCount int
	ErrorCount   int
	TestSuites   []TestSuite
}

// TestSuite ..
type TestSuite struct {
	Name         string
	TotalTime    float64
	TotalCount   int
	SuccessCount int
	FailedCount  int
	SkippedCount int
	ErrorCount   int
}

// NewPipeline ..
func NewPipeline(ctx context.Context, gp goGitlab.Pipeline) Pipeline {
	var (
		coverage  float64
		err       error
		timestamp float64
		starttime float64
	)

	if gp.Coverage != "" {
		coverage, err = strconv.ParseFloat(gp.Coverage, 64)
		if err != nil {
			log.WithContext(ctx).
				WithField("error", err.Error()).
				Warnf("could not parse coverage string returned from GitLab API '%s' into Float64", gp.Coverage)
		}
	}

	if gp.UpdatedAt != nil {
		timestamp = float64(gp.UpdatedAt.Unix())
	}

	if gp.StartedAt != nil {
		starttime = float64(gp.StartedAt.Unix())
	}

	return Pipeline{
		ID:                    gp.ID,
		Coverage:              coverage,
		Timestamp:             timestamp,
		StartTime:             starttime,
		DurationSeconds:       float64(gp.Duration),
		QueuedDurationSeconds: float64(gp.QueuedDuration),
		Status:                gp.Status,
	}
}

// NewTestReport ..
func NewTestReport(gtr goGitlab.PipelineTestReport) TestReport {
	testSuites := []TestSuite{}

	for _, x := range gtr.TestSuites {
		testSuites = append(testSuites, NewTestSuite(x))
	}

	return TestReport{
		TotalTime:    gtr.TotalTime,
		TotalCount:   gtr.TotalCount,
		SuccessCount: gtr.SuccessCount,
		FailedCount:  gtr.FailedCount,
		SkippedCount: gtr.SkippedCount,
		ErrorCount:   gtr.ErrorCount,
		TestSuites:   testSuites,
	}
}

// NewTestSuite ..
func NewTestSuite(gts *goGitlab.PipelineTestSuites) TestSuite {
	return TestSuite{
		Name:         gts.Name,
		TotalTime:    gts.TotalTime,
		TotalCount:   gts.TotalCount,
		SuccessCount: gts.SuccessCount,
		FailedCount:  gts.FailedCount,
		SkippedCount: gts.SkippedCount,
		ErrorCount:   gts.ErrorCount,
	}
}
