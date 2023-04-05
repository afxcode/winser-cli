package winser

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// Control
func Start(name string) (err error) {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("Could not access %s: %v", name, err)
	}
	defer s.Close()

	sStatus, err := s.Query()
	if err != nil {
		err = fmt.Errorf("Could not access %s: %v", name, err)
		return
	}

	status := statusFromSvc(sStatus.State)
	if status == StartPending || status == PausePending || status == StopPending || status == ContinuePending {
		err = fmt.Errorf("Service in pending state")
		return
	}
	if status == Running {
		return
	}

	if status == Stopped {
		if err = s.Start(); err != nil {
			err = fmt.Errorf("Could not start service %s: %v", name, err)
		}
		return
	}

	if status == Paused {
		if _, err = s.Control(svc.Continue); err != nil {
			err = fmt.Errorf("Could not continue service %s: %v", name, err)
		}
	}
	return
}

func Stop(name string) (err error) {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("Could not access %s: %v", name, err)
	}
	defer s.Close()

	sStatus, err := s.Query()
	if err != nil {
		err = fmt.Errorf("Could not access %s: %v", name, err)
		return
	}

	status := statusFromSvc(sStatus.State)
	if status == StartPending || status == PausePending || status == StopPending || status == ContinuePending {
		err = fmt.Errorf("Service in pending state")
		return
	}
	if status == Stopped {
		return
	}

	if _, err = s.Control(svc.Stop); err != nil {
		err = fmt.Errorf("Could not stop service %s: %v", name, err)
	}
	return
}

func Pause(name string) (err error) {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("Could not access %s: %v", name, err)
	}
	defer s.Close()

	sStatus, err := s.Query()
	if err != nil {
		err = fmt.Errorf("Could not access %s: %v", name, err)
		return
	}

	status := statusFromSvc(sStatus.State)
	if status == StartPending || status == PausePending || status == StopPending || status == ContinuePending {
		err = fmt.Errorf("Service in pending state")
		return
	}
	if status != Running || status == Stopped {
		err = fmt.Errorf("Service is not running")
		return
	}
	if status == Paused {
		return
	}

	if _, err = s.Control(svc.Pause); err != nil {
		err = fmt.Errorf("Could not pause service %s: %v", name, err)
	}
	return
}
