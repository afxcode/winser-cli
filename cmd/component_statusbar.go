package main

type StatusBar struct {
	Type StatusBarType
	Msg  string
}

func (s StatusBar) Render() string {
	switch s.Type {
	case StatusInfo:
		return info("INFO: " + s.Msg)
	case StatusSuccess:
		return success("SUCCESS: " + s.Msg)
	case StatusWarning:
		return warning("WARNING: " + s.Msg)
	case StatusError:
		return danger("ERROR: " + s.Msg)
	default:
		return secondary(s.Msg)
	}
}

func (s *StatusBar) Info(msg string) {
	s.Type = StatusInfo
	s.Msg = msg
}

func (s *StatusBar) Success(msg string) {
	s.Type = StatusSuccess
	s.Msg = msg
}

func (s *StatusBar) Warning(msg string) {
	s.Type = StatusWarning
	s.Msg = msg
}

func (s *StatusBar) Error(msg string) {
	s.Type = StatusError
	s.Msg = msg
}

type StatusBarType string

const (
	StatusInfo    StatusBarType = "info"
	StatusSuccess StatusBarType = "success"
	StatusWarning StatusBarType = "warning"
	StatusError   StatusBarType = "error"
)
