package stats

type Snapshot struct {
	TotalMeetings       int
	TotalBornAgents     int
	TotalDeadAgents     int
	SameIslandMeetings  int
	CrossIslandMeetings int
	BestScore           float64
	BestIsland          int

	AgentsPerIsland   map[int]int
	MeetingsPerIsland map[int]int
}

func (s Snapshot) SameIslandRatio() float64 {
	total := s.SameIslandMeetings + s.CrossIslandMeetings
	if total == 0 {
		return 0
	}
	return float64(s.SameIslandMeetings) / float64(total)
}
