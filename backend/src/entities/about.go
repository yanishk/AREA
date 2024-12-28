package entities

// About
type About struct {
	Client AboutClient `json:"client"`
	Server AboutServer `json:"server"`
}

// About Client
type AboutClient struct {
	Host string `json:"host"`
}

// About Server
type AboutServer struct {
	CurrentTime int64          `json:"current_time"`
	Services    []AboutService `json:"services"`
}

type AboutService struct {
	Name      string          `json:"name"`
	Actions   []AboutAction   `json:"actions"`
	Reactions []AboutReaction `json:"reactions"`
}

type AboutAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AboutReaction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
