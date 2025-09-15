package types

import (
	"fmt"
	"time"
)

type StatusWithTimestamp struct {
	Value     string    `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func NewStatusWithTimestamp(value string) StatusWithTimestamp {
	return StatusWithTimestamp{
		Value:     value,
		Timestamp: time.Now(),
	}
}

func (s StatusWithTimestamp) IsEmpty() bool {
	return s.Value == ""
}

func (s StatusWithTimestamp) Age() time.Duration {
	return time.Since(s.Timestamp)
}

type GenericID struct {
	value string
}

func NewGenericID(value string) (GenericID, error) {
	if value == "" {
		return GenericID{}, fmt.Errorf("ID cannot be empty")
	}
	return GenericID{value: value}, nil
}

func (i GenericID) Value() string {
	return i.value
}

func (i GenericID) String() string {
	return i.value
}

func (i GenericID) IsEmpty() bool {
	return i.value == ""
}

type TimestampWrapper struct {
	value time.Time
}

func NewTimestampWrapper(t time.Time) TimestampWrapper {
	return TimestampWrapper{value: t}
}

func Now() TimestampWrapper {
	return TimestampWrapper{value: time.Now()}
}

func (t TimestampWrapper) Value() time.Time {
	return t.value
}

func (t TimestampWrapper) String() string {
	return t.value.Format(time.RFC3339)
}

func (t TimestampWrapper) IsZero() bool {
	return t.value.IsZero()
}

func (t TimestampWrapper) Age() time.Duration {
	return time.Since(t.value)
}
