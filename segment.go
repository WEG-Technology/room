package room

import "time"

type SegmentSchema struct {
	startedAt   time.Time
	endedAt     time.Time
	elapsedTime float64
}

type ISegment interface {
	start(t time.Time) ISegment
	startNow() ISegment
	End() ISegment
	GetElapsedTime() float64
}

func (s *SegmentSchema) start(t time.Time) ISegment {
	s.startedAt = t
	return s
}

func (s *SegmentSchema) startNow() ISegment {
	return s.start(time.Now())
}

func (s *SegmentSchema) End() ISegment {
	s.endedAt = time.Now()
	s.elapsedTime = time.Since(s.startedAt).Seconds()
	return s
}

func (s *SegmentSchema) GetElapsedTime() float64 {
	return s.elapsedTime
}

func StartSegmentNow() ISegment {
	s := new(SegmentSchema)
	return s.startNow()
}
