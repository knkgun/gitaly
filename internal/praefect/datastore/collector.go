package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitaly/v14/internal/praefect/datastore/glsql"
)

var (
	// This is kept for backwards compatibility as some alerting rules depend on this.
	// The unavailable repositories is a more accurate description for the metric and
	// is exported below so we can migrate to it.
	descReadOnlyRepositories = prometheus.NewDesc(
		"gitaly_praefect_read_only_repositories",
		"Number of repositories in read-only mode within a virtual storage.",
		[]string{"virtual_storage"},
		nil,
	)

	descUnavailableRepositories = prometheus.NewDesc(
		"gitaly_praefect_unavailable_repositories",
		"Number of repositories that have no healthy, up to date replicas.",
		[]string{"virtual_storage"},
		nil,
	)

	descriptions = []*prometheus.Desc{descReadOnlyRepositories, descUnavailableRepositories}
)

// RepositoryStoreCollector collects metrics from the RepositoryStore.
type RepositoryStoreCollector struct {
	log             logrus.FieldLogger
	db              glsql.Querier
	virtualStorages []string
	timeout         time.Duration
}

// NewRepositoryStoreCollector returns a new collector.
func NewRepositoryStoreCollector(log logrus.FieldLogger, virtualStorages []string, db glsql.Querier, timeout time.Duration) *RepositoryStoreCollector {
	return &RepositoryStoreCollector{
		log:             log.WithField("component", "RepositoryStoreCollector"),
		db:              db,
		virtualStorages: virtualStorages,
		timeout:         timeout,
	}
}

//nolint: revive,stylecheck // This is unintentionally missing documentation.
func (c *RepositoryStoreCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range descriptions {
		ch <- desc
	}
}

//nolint: revive,stylecheck // This is unintentionally missing documentation.
func (c *RepositoryStoreCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.TODO(), c.timeout)
	defer cancel()

	unavailableCounts, err := CountUnavailableRepositories(ctx, c.db, c.virtualStorages)
	if err != nil {
		c.log.WithError(err).Error("failed collecting read-only repository count metric")
		return
	}

	for _, vs := range c.virtualStorages {
		for _, desc := range descriptions {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(unavailableCounts[vs]), vs)
		}
	}
}

// CountUnavailableRepositories queries the number of unavailable repositories from the database.
// A repository is unavailable when it has no replicas that can act as a primary, indicating
// they are either unhealthy or out of date.
func CountUnavailableRepositories(ctx context.Context, db glsql.Querier, virtualStorages []string) (map[string]int, error) {
	rows, err := db.QueryContext(ctx, `
SELECT virtual_storage, COUNT(*)
FROM repositories
WHERE NOT EXISTS (
	SELECT FROM valid_primaries
	WHERE valid_primaries.repository_id = repositories.repository_id
) AND repositories.virtual_storage = ANY($1)
GROUP BY virtual_storage
	`, virtualStorages)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	vsUnavailable := make(map[string]int)
	for rows.Next() {
		var vs string
		var count int

		if err := rows.Scan(&vs, &count); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		vsUnavailable[vs] = count
	}

	return vsUnavailable, rows.Err()
}
