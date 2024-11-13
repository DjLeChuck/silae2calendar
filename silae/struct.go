package silae

type ApiResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type Company struct {
	Id int `json:"id"`
}

type CurrentCollaborator struct {
	Id      int     `json:"id"`
	Company Company `json:"company"`
}

type UserData struct {
	Firstname           string              `json:"firstname"`
	Lastname            string              `json:"lastname"`
	Token               string              `json:"token"`
	RefreshToken        string              `json:"refresh_token"`
	CurrentCollaborator CurrentCollaborator `json:"current_collaborator_base"`
}

type Request struct {
	Status string `json:"status"`
}

type Freedays struct {
	Abbr             string  `json:"abbreviation"`
	DateStart        string  `json:"date_start"`
	DateStartDayType string  `json:"date_start_day_type"`
	DateEnd          string  `json:"date_end"`
	DateEndDayType   string  `json:"date_end_day_type"`
	Request          Request `json:"request"`
}

type CollaboratorFreedays struct {
	Freedays []Freedays `json:"freedays"`
}

type FreedaysData struct {
	CollaboratorFreedays []CollaboratorFreedays `json:"collaborator_freedays"`
}

type Criteria interface{}

type DateRangeCriteria struct {
	Type string `json:"_type"`
	Min  string `json:"min"`
	Max  string `json:"max"`
}

type StringValueCriteria struct {
	Type  string `json:"_type"`
	Value int    `json:"value"`
}

type ListStringValueCriteria struct {
	Type   string `json:"_type"`
	Values []int  `json:"values"`
}

type Filter struct {
	Type     string   `json:"_type"`
	Name     string   `json:"name"`
	Criteria Criteria `json:"criteria"`
}

type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type RequestPayload struct {
	Filters []Filter `json:"filters"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Sort    Sort     `json:"sort"`
}
