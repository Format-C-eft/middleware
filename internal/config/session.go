package config

// AgentInfo - agent info decription
type AgentInfo struct {
	IP             string   `json:"IP"`
	Engine         string   `json:"Engine"`
	Browser        string   `json:"Browser"`
	BrowserVersion string   `json:"BrowserVersion"`
	Mobile         bool     `json:"Mobile"`
	Platform       string   `json:"Platform"`
	OS             string   `json:"OS"`
	CreateTime     DateTime `json:"CreateTime"`
	LastTime       DateTime `json:"LastTime"`
	EndsTime       DateTime `json:"EndsTime"`
}

// SessionInfo - session info description
type SessionInfo struct {
	Сredentials struct {
		Login    string
		Password string
	}
	Session struct {
		ID     string
		Cookie string
	}
	Info AgentInfo
}

// GetAgentInfo - get agent info
func (a *SessionInfo) GetAgentInfo() AgentInfo {
	return a.Info
}
