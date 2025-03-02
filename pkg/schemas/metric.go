package schemas

import (
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricKindCoverage refers to the coerage of a job/pipeline.
	MetricKindCoverage MetricKind = iota

	// MetricKindDurationSeconds ..
	MetricKindDurationSeconds

	// MetricKindEnvironmentBehindCommitsCount ..
	MetricKindEnvironmentBehindCommitsCount

	// MetricKindEnvironmentBehindDurationSeconds ..
	MetricKindEnvironmentBehindDurationSeconds

	// MetricKindEnvironmentDeploymentCount ..
	MetricKindEnvironmentDeploymentCount

	// MetricKindEnvironmentDeploymentDurationSeconds ..
	MetricKindEnvironmentDeploymentDurationSeconds

	// MetricKindEnvironmentDeploymentJobID ..
	MetricKindEnvironmentDeploymentJobID

	// MetricKindEnvironmentDeploymentStatus ..
	MetricKindEnvironmentDeploymentStatus

	// MetricKindEnvironmentDeploymentTimestamp ..
	MetricKindEnvironmentDeploymentTimestamp

	// MetricKindEnvironmentInformation ..
	MetricKindEnvironmentInformation

	// MetricKindID ..
	MetricKindID

	// MetricKindJobArtifactSizeBytes ..
	MetricKindJobArtifactSizeBytes

	// MetricKindJobDurationSeconds ..
	MetricKindJobDurationSeconds

	// MetricKindJobID ..
	MetricKindJobID

	// MetricKindJobQueuedDurationSeconds ..
	MetricKindJobQueuedDurationSeconds

	// MetricKindJobRunCount ..
	MetricKindJobRunCount

	// MetricKindJobStatus ..
	MetricKindJobStatus

	// MetricKindJobTimestamp ..
	MetricKindJobTimestamp

	// MetricKindJobStartTime ..
	MetricKindJobStartTime

	// MetricKindQueuedDurationSeconds ..
	MetricKindQueuedDurationSeconds

	// MetricKindRunCount ..
	MetricKindRunCount

	// MetricKindStatus ..
	MetricKindStatus

	// MetricKindTimestamp ..
	MetricKindTimestamp

	// MetricKindStartTime ..
	MetricKindStartTime

	// MetricKindTestReportTotalTime ..
	MetricKindTestReportTotalTime

	// MetricKindTestReportTotalCount ..
	MetricKindTestReportTotalCount

	// MetricKindTestReportSuccessCount ..
	MetricKindTestReportSuccessCount

	// MetricKindTestReportFailedCount ..
	MetricKindTestReportFailedCount

	// MetricKindTestReportSkippedCount ..
	MetricKindTestReportSkippedCount

	// MetricKindTestReportErrorCount ..
	MetricKindTestReportErrorCount

	// MetricKindTestSuiteTotalTime ..
	MetricKindTestSuiteTotalTime

	// MetricKindTestSuiteTotalCount ..
	MetricKindTestSuiteTotalCount

	// MetricKindTestSuiteSuccessCount ..
	MetricKindTestSuiteSuccessCount

	// MetricKindTestSuiteFailedCount ..
	MetricKindTestSuiteFailedCount

	// MetricKindTestSuiteSkippedCount ..
	MetricKindTestSuiteSkippedCount

	// MetricKindTestSuiteErrorCount ..
	MetricKindTestSuiteErrorCount
)

// MetricKind ..
type MetricKind int32

// Metric ..
type Metric struct {
	Kind   MetricKind
	Labels prometheus.Labels
	Value  float64
}

// MetricKey ..
type MetricKey string

// Metrics ..
type Metrics map[MetricKey]Metric

// Key ..
func (m Metric) Key() MetricKey {
	key := strconv.Itoa(int(m.Kind))

	switch m.Kind {
	case MetricKindCoverage, MetricKindDurationSeconds, MetricKindID, MetricKindQueuedDurationSeconds, MetricKindRunCount, MetricKindStatus, MetricKindTimestamp, MetricKindStartTime, MetricKindTestReportTotalCount, MetricKindTestReportErrorCount, MetricKindTestReportFailedCount, MetricKindTestReportSkippedCount, MetricKindTestReportSuccessCount, MetricKindTestReportTotalTime:
		key += fmt.Sprintf("%v", []string{
			m.Labels["project"],
			m.Labels["kind"],
			m.Labels["ref"],
		})

	case MetricKindJobArtifactSizeBytes, MetricKindJobDurationSeconds, MetricKindJobID, MetricKindJobQueuedDurationSeconds, MetricKindJobRunCount, MetricKindJobStatus, MetricKindJobTimestamp, MetricKindJobStartTime:
		key += fmt.Sprintf("%v", []string{
			m.Labels["project"],
			m.Labels["kind"],
			m.Labels["ref"],
			m.Labels["stage"],
			m.Labels["job_name"],
		})

	case MetricKindEnvironmentBehindCommitsCount, MetricKindEnvironmentBehindDurationSeconds, MetricKindEnvironmentDeploymentCount, MetricKindEnvironmentDeploymentDurationSeconds, MetricKindEnvironmentDeploymentJobID, MetricKindEnvironmentDeploymentStatus, MetricKindEnvironmentDeploymentTimestamp, MetricKindEnvironmentInformation:
		key += fmt.Sprintf("%v", []string{
			m.Labels["project"],
			m.Labels["environment"],
		})

	case MetricKindTestSuiteErrorCount, MetricKindTestSuiteFailedCount, MetricKindTestSuiteSkippedCount, MetricKindTestSuiteSuccessCount, MetricKindTestSuiteTotalCount, MetricKindTestSuiteTotalTime:
		key += fmt.Sprintf("%v", []string{
			m.Labels["project"],
			m.Labels["kind"],
			m.Labels["ref"],
			m.Labels["test_suite_name"],
		})
	}

	// If the metric is a "status" one, add the status label
	switch m.Kind {
	case MetricKindJobStatus, MetricKindEnvironmentDeploymentStatus, MetricKindStatus:
		key += m.Labels["status"]
	}

	return MetricKey(strconv.Itoa(int(crc32.ChecksumIEEE([]byte(key)))))
}
