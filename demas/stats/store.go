package stats

import (
	"maps"
	"sync"
)

type ScoreComparator func(newScore, currentBest float64) bool

type Store struct {
	mu sync.RWMutex

	AgentsPerIsland     map[int]int
	MeetingsPerIsland   map[int]int
	TotalMeetings       int
	TotalBornAgents     int
	TotalDeadAgents     int
	SameIslandMeetings  int
	CrossIslandMeetings int
	BestScore           float64
	BestIsland          int

	isBetter ScoreComparator
}

func NewStore(initialAgentsPerIsland map[int]int, maximum float64, isBetter ScoreComparator) *Store {
	agentsCopy := make(map[int]int, len(initialAgentsPerIsland))
	maps.Copy(agentsCopy, initialAgentsPerIsland)

	return &Store{
		AgentsPerIsland:   agentsCopy,
		MeetingsPerIsland: make(map[int]int),
		BestScore:         maximum,
		isBetter:          isBetter,
	}
}

func (s *Store) Apply(event Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	newBest := false

	switch event.EventType {
	case EventBorn:
		s.AgentsPerIsland[event.PrimaryIslandID]++
		s.TotalBornAgents++

		if s.isBetter(event.Score, s.BestScore) {
			s.BestScore = event.Score
			s.BestIsland = event.PrimaryIslandID
			newBest = true
		}
	case EventDead:
		s.AgentsPerIsland[event.PrimaryIslandID]--
		s.TotalDeadAgents++

	case EventMeeting:
		s.TotalMeetings++
		s.MeetingsPerIsland[event.PrimaryIslandID]++

		if event.SameIsland {
			s.SameIslandMeetings++
		} else {
			s.CrossIslandMeetings++
		}

	}

	return newBest
}

func (s *Store) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	agentsCopy := make(map[int]int, len(s.AgentsPerIsland))
	maps.Copy(agentsCopy, s.AgentsPerIsland)

	meetingsCopy := make(map[int]int, len(s.MeetingsPerIsland))
	maps.Copy(meetingsCopy, s.MeetingsPerIsland)

	return Snapshot{
		TotalMeetings:       s.TotalMeetings,
		TotalBornAgents:     s.TotalBornAgents,
		TotalDeadAgents:     s.TotalDeadAgents,
		SameIslandMeetings:  s.SameIslandMeetings,
		CrossIslandMeetings: s.CrossIslandMeetings,
		BestScore:           s.BestScore,
		BestIsland:          s.BestIsland,
		AgentsPerIsland:     agentsCopy,
		MeetingsPerIsland:   meetingsCopy,
	}
}
