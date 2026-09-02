package clickhousestore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/domain"
	"github.com/VictorOdoyo/traceforge-reliability-lab/internal/ports"
)

type Repository struct {
	conn     driver.Conn
	fallback ports.Repository
}

func New(ctx context.Context, dsn string, fallback ports.Repository) (*Repository, error) {
	options, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse ClickHouse DSN: %w", err)
	}
	connection, err := clickhouse.Open(options)
	if err != nil {
		return nil, fmt.Errorf("open ClickHouse client: %w", err)
	}
	if err := connection.Ping(ctx); err != nil {
		return nil, fmt.Errorf("connect to ClickHouse: %w", err)
	}
	repository := &Repository{conn: connection, fallback: fallback}
	if err := repository.ensureSchema(ctx); err != nil {
		connection.Close()
		return nil, err
	}
	return repository, nil
}
func (r *Repository) Close() error { return r.conn.Close() }
func (r *Repository) ensureSchema(ctx context.Context) error {
	return r.conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS telemetry_spans (
    started_at DateTime64(6, 'UTC'), trace_id String, span_id String, parent_span_id String,
    service LowCardinality(String), operation String, duration_ns UInt64,
    status LowCardinality(String), attributes_json String
) ENGINE = MergeTree PARTITION BY toDate(started_at) ORDER BY (service, started_at, trace_id, span_id)
TTL started_at + INTERVAL 30 DAY DELETE`)
}

func (r *Repository) SaveTrace(ctx context.Context, trace domain.Trace) error {
	batch, err := r.conn.PrepareBatch(ctx, "INSERT INTO telemetry_spans")
	if err != nil {
		return err
	}
	for _, span := range trace.Spans {
		attributes, _ := json.Marshal(span.Attributes)
		if err := batch.Append(span.StartedAt, span.TraceID, span.SpanID, span.ParentID, span.Service, span.Operation, uint64(span.Duration), string(span.Status), string(attributes)); err != nil {
			return err
		}
	}
	return batch.Send()
}

func (r *Repository) ListTraces(ctx context.Context, filter ports.TraceFilter) ([]domain.Trace, error) {
	conditions, args := []string{"1 = 1"}, []any{}
	if filter.Service != "" {
		conditions = append(conditions, "service = ?")
		args = append(args, filter.Service)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, string(filter.Status))
	}
	if filter.MinDuration > 0 {
		conditions = append(conditions, "duration_ns >= ?")
		args = append(args, uint64(filter.MinDuration))
	}
	if !filter.From.IsZero() {
		conditions = append(conditions, "started_at >= ?")
		args = append(args, filter.From)
	}
	if !filter.To.IsZero() {
		conditions = append(conditions, "started_at <= ?")
		args = append(args, filter.To)
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 200
	}
	query := fmt.Sprintf(`SELECT started_at, trace_id, span_id, parent_span_id, service, operation, duration_ns, status, attributes_json
FROM telemetry_spans WHERE trace_id IN (
 SELECT trace_id FROM telemetry_spans WHERE %s GROUP BY trace_id ORDER BY max(started_at) DESC LIMIT ?
) ORDER BY trace_id, started_at`, strings.Join(conditions, " AND "))
	args = append(args, limit)
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	grouped := map[string][]domain.Span{}
	order := []string{}
	for rows.Next() {
		var started time.Time
		var traceID, spanID, parentID, service, operation, status, attributesJSON string
		var duration uint64
		if err := rows.Scan(&started, &traceID, &spanID, &parentID, &service, &operation, &duration, &status, &attributesJSON); err != nil {
			return nil, err
		}
		if _, exists := grouped[traceID]; !exists {
			order = append(order, traceID)
		}
		attributes := map[string]string{}
		_ = json.Unmarshal([]byte(attributesJSON), &attributes)
		grouped[traceID] = append(grouped[traceID], domain.Span{TraceID: traceID, SpanID: spanID, ParentID: parentID, Service: service, Operation: operation, StartedAt: started, Duration: time.Duration(duration), Status: domain.SpanStatus(status), Attributes: attributes})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.Trace, 0, len(grouped))
	for _, id := range order {
		trace, err := domain.BuildTrace(id, grouped[id])
		if err != nil {
			return nil, err
		}
		result = append(result, trace)
	}
	return result, nil
}

func (r *Repository) ListServices(ctx context.Context) ([]domain.Service, error) {
	return r.fallback.ListServices(ctx)
}
func (r *Repository) ListSLOs(ctx context.Context) ([]domain.SLO, error) {
	return r.fallback.ListSLOs(ctx)
}
func (r *Repository) ListDeployments(ctx context.Context) ([]domain.Deployment, error) {
	return r.fallback.ListDeployments(ctx)
}
func (r *Repository) ListAlerts(ctx context.Context) ([]domain.Alert, error) {
	return r.fallback.ListAlerts(ctx)
}
func (r *Repository) SaveAlert(ctx context.Context, value domain.Alert) error {
	return r.fallback.SaveAlert(ctx, value)
}
func (r *Repository) ListIncidents(ctx context.Context) ([]domain.Incident, error) {
	return r.fallback.ListIncidents(ctx)
}
func (r *Repository) SaveIncident(ctx context.Context, value domain.Incident) error {
	return r.fallback.SaveIncident(ctx, value)
}
func (r *Repository) ListSuppressions(ctx context.Context) ([]domain.SuppressionRule, error) {
	return r.fallback.ListSuppressions(ctx)
}
