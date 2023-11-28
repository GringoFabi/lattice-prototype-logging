package analyze

type LogsRequest struct {
	Logs []Log `json:"logs"`
}

type SingleLogRequest struct {
	Log Log `json:"log"`
}

type Log struct {
	Action   string   `json:"action"`
	NodeData NodeData `json:"node"`
	Time     int      `json:"time"`
}

type NodeData struct {
	Node      int        `json:"node"`
	Position  []float64  `json:"position"`
	Edges     [][]string `json:"edges"`
	Toplabel  []string   `json:"toplabel"`
	Botlabel  []string   `json:"botlabel"`
	Valuation string
}

type Logs []Log
type ActionCount map[string]int

func (ls Logs) CountActions() ActionCount {
	countPerAction := make(ActionCount)
	for _, l := range ls {
		if _, ok := countPerAction[l.Action]; !ok {
			countPerAction[l.Action] = 1
		} else {
			countPerAction[l.Action] += 1
		}
	}

	return countPerAction
}
