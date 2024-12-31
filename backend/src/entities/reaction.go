package entities

type Reaction struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	ServiceId   string                   `json:"serviceid"`
	NbParam     int                      `json:"nbparam"`
	Parameters  []map[string]interface{} `json:"parameters"`
}
