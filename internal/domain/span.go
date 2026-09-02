package domain

import (
	"errors"
	"time"
)

type SpanStatus string

const (
	SpanOK    SpanStatus = "ok"
	SpanError SpanStatus = "error"
	SpanUnset SpanStatus = "unset"
)

type Span struct {
	TraceID    string            `json:"traceId"`
	SpanID     string            `json:"spanId"`
	ParentID   string            `json:"parentId,omitempty"`
	Service    string            `json:"service"`
	Operation  string            `json:"operation"`
	StartedAt  time.Time         `json:"startedAt"`
	Duration   time.Duration     `json:"duration"`
	Status     SpanStatus        `json:"status"`
	Attributes map[string]string `json:"attributes"`
}

func (s Span) Validate() error {
	if s.TraceID == "" || s.SpanID == "" || s.Service == "" || s.Operation == "" {
		return errors.New("trace, span, service, and operation identifiers are required")
	}
	if s.Duration < 0 {
		return errors.New("span duration cannot be negative")
	}
	if s.Status != SpanOK && s.Status != SpanError && s.Status != SpanUnset {
		return errors.New("unsupported span status")
	}
	return nil
}

func (s Span) DurationMillis() float64 { return float64(s.Duration) / float64(time.Millisecond) }
