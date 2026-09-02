package domain

import (
	"errors"
	"sort"
	"time"
)

type Trace struct {
	ID          string        `json:"id"`
	Spans       []Span        `json:"spans"`
	StartedAt   time.Time     `json:"startedAt"`
	Duration    time.Duration `json:"duration"`
	RootService string        `json:"rootService"`
	HasError    bool          `json:"hasError"`
}

func BuildTrace(id string, spans []Span) (Trace, error) {
	if id == "" || len(spans) == 0 {
		return Trace{}, errors.New("trace id and spans are required")
	}
	copySpans := append([]Span(nil), spans...)
	for _, span := range copySpans {
		if err := span.Validate(); err != nil {
			return Trace{}, err
		}
		if span.TraceID != id {
			return Trace{}, errors.New("all spans must belong to the trace")
		}
	}
	sort.Slice(copySpans, func(i, j int) bool { return copySpans[i].StartedAt.Before(copySpans[j].StartedAt) })
	end := copySpans[0].StartedAt.Add(copySpans[0].Duration)
	hasError, root := false, copySpans[0].Service
	for _, span := range copySpans {
		if candidate := span.StartedAt.Add(span.Duration); candidate.After(end) {
			end = candidate
		}
		if span.Status == SpanError {
			hasError = true
		}
		if span.ParentID == "" {
			root = span.Service
		}
	}
	return Trace{ID: id, Spans: copySpans, StartedAt: copySpans[0].StartedAt, Duration: end.Sub(copySpans[0].StartedAt), RootService: root, HasError: hasError}, nil
}
