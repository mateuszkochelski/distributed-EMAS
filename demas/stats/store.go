package stats

import (
	"maps"
	"math"
	"sync"
)

type Store struct {
	mu sync.RWMutex

	AgentsPerIsland     map[int]int
	MeetingsPerIsland   map[int]int
	TotalMeetings       int
	SameIslandMeetings  int
	CrossIslandMeetings int
	BestScore           float64
}

func NewStore(initialAgentsPerIsland map[int]int) *Store {
	agentsCopy := make(map[int]int, len(initialAgentsPerIsland))
	for k, v := range initialAgentsPerIsland {
		agentsCopy[k] = v
	}

	return &Store{
		AgentsPerIsland:   agentsCopy,
		MeetingsPerIsland: make(map[int]int),
		BestScore:         math.Inf(1),
	}
}

func (s *Store) Apply(event Event) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	newBest := false

	switch event.EventType {
	case EventBorn:
		s.AgentsPerIsland[event.PrimaryIslandID]++

	case EventDead:
		s.AgentsPerIsland[event.PrimaryIslandID]--

	case EventMeeting:
		s.TotalMeetings++
		s.MeetingsPerIsland[event.PrimaryIslandID]++

		if event.SameIsland {
			s.SameIslandMeetings++
		} else {
			s.CrossIslandMeetings++
		}

		if event.Score < s.BestScore {
			s.BestScore = event.Score
			newBest = true
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
		SameIslandMeetings:  s.SameIslandMeetings,
		CrossIslandMeetings: s.CrossIslandMeetings,
		BestScore:           s.BestScore,
		AgentsPerIsland:     agentsCopy,
		MeetingsPerIsland:   meetingsCopy,
	}
}
