package room

import (
	"testing"
	"time"
)

func TestSegmentSchema_Start(t *testing.T) {
	s := new(SegmentSchema)
	startTime := time.Now()
	s.start(startTime)

	if s.startedAt != startTime {
		t.Errorf("Expected startedAt to be %v, but got %v", startTime, s.startedAt)
	}
}

func TestSegmentSchema_StartNow(t *testing.T) {
	s := new(SegmentSchema)
	s.startNow()

	if s.startedAt.IsZero() {
		t.Error("Expected startedAt to be set, but it is zero")
	}
}

func TestSegmentSchema_End(t *testing.T) {
	s := new(SegmentSchema)
	startTime := time.Now().Add(-time.Minute) // Simulate starting in the past
	s.start(startTime)
	s.End()

	if s.endedAt.IsZero() {
		t.Error("Expected endedAt to be set, but it is zero")
	}

	if s.elapsedTime <= 0 {
		t.Errorf("Expected elapsedTime to be greater than 0, but got %f", s.elapsedTime)
	}
}

func TestSegmentSchema_GetElapsedTime(t *testing.T) {
	s := new(SegmentSchema)
	startTime := time.Now().Add(-time.Minute) // Simulate starting in the past
	s.start(startTime)
	s.End()

	elapsedTime := s.GetElapsedTime()

	if elapsedTime <= 0 {
		t.Errorf("Expected GetElapsedTime to return a value greater than 0, but got %f", elapsedTime)
	}
}

func TestIntegration(t *testing.T) {
	s := StartSegmentNow()
	time.Sleep(time.Second) // Simulate some work being done
	s.End()

	elapsedTime := s.GetElapsedTime()

	if elapsedTime <= 0 {
		t.Errorf("Expected GetElapsedTime to return a value greater than 0, but got %f", elapsedTime)
	}
}
