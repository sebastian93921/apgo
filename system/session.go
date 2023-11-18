package system

type Session struct {
	Settings       *Settings
	GlobalSettings *GlobalSettings
}

// NewSession creates a new session
func NewSession(cfg *Settings) *Session {
	gc := &GlobalSettings{}
	s := &Session{
		Settings:       cfg,
		GlobalSettings: gc,
	}
	return s
}
