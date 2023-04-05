package winser

import (
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type Service struct {
	Name        string
	DisplayName string
	Description string
	Executable  string
	Argument    string
	StartMode
	Status
}

type StartMode string

const (
	Auto             StartMode = "auto"
	Manual           StartMode = "manual"
	Disabled         StartMode = "disabled"
	StartModeUnknown StartMode = "unknown"
)

func (s StartMode) toSvcStartType() uint32 {
	switch s {
	case Auto:
		return mgr.StartAutomatic
	case Manual:
		return mgr.StartManual
	case Disabled:
		return mgr.StartDisabled
	}
	return 0
}

func startModeFromSvc(startType uint32) StartMode {
	switch startType {
	case mgr.StartAutomatic:
		return Auto
	case mgr.StartManual:
		return Manual
	case mgr.StartDisabled:
		return Disabled
	default:
		return StartModeUnknown
	}
}

type Status string

const (
	Running         Status = "running"
	Paused          Status = "paused"
	Stopped         Status = "stopped"
	StartPending    Status = "pending_start"
	PausePending    Status = "pending_pause"
	StopPending     Status = "pending_stop"
	ContinuePending Status = "pending_continue"
	StatusUnknown   Status = "unknown"
)

func statusFromSvc(state svc.State) Status {
	switch state {
	case svc.Running:
		return Running
	case svc.Paused:
		return Paused
	case svc.Stopped:
		return Stopped
	case svc.StartPending:
		return StartPending
	case svc.PausePending:
		return PausePending
	case svc.StopPending:
		return StopPending
	case svc.ContinuePending:
		return ContinuePending
	default:
		return StatusUnknown
	}
}
