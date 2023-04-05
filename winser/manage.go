package winser

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/go-shellwords"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func GetServices() (services []Service, err error) {
	m, err := mgr.Connect()
	if err != nil {
		return
	}
	defer m.Disconnect()

	names, err := m.ListServices()
	if err != nil {
		return
	}

	for _, name := range names {
		s, e := m.OpenService(name)
		if e != nil {
			continue
		}
		defer s.Close()

		sConf, e := s.Config()
		if e != nil {
			continue
		}

		sStatus, e := s.Query()
		if e != nil {
			continue
		}

		services = append(services, Service{
			Name:        name,
			DisplayName: sConf.DisplayName,
			Status:      statusFromSvc(sStatus.State),
		})
	}
	return
}

func Find(name string) (service Service, err error) {
	m, err := mgr.Connect()
	if err != nil {
		return
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		err = fmt.Errorf("Could not access %s: %v", name, err)
		return
	}
	defer s.Close()

	sConf, err := s.Config()
	if err != nil {
		err = fmt.Errorf("Could not access config: %v", err)
		return
	}

	sStatus, err := s.Query()
	if err != nil {
		err = fmt.Errorf("Could not access %s: %v", name, err)
		return
	}

	args, _ := shellwords.Parse(sConf.BinaryPathName)
	executable := ""
	if len(args) > 0 {
		executable = args[0]
	}
	argument := ""
	if len(args) > 1 {
		argument = strings.Join(args[1:], " ")
	}

	service = Service{
		Name:        name,
		DisplayName: sConf.DisplayName,
		Description: sConf.Description,
		Executable:  executable,
		Argument:    argument,
		StartMode:   startModeFromSvc(sConf.StartType),
		Status:      statusFromSvc(sStatus.State),
	}
	return
}

// func (a *App) SelectExecutable(currentExecutable string) (selectedExecutable string, err error) {
// 	dir := filepath.Dir(currentExecutable)
// 	if dir == "." {
// 		dir = "c:\\"
// 	}
// 	filename := filepath.Base(currentExecutable)
// 	if filename == "." {
// 		filename = ""
// 	}
//
// 	selectedExecutable, err = runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
// 		DefaultDirectory: dir,
// 		DefaultFilename:  filename,
// 		Title:            "Select Executable",
// 		Filters: []runtime.FileFilter{
// 			{
// 				DisplayName: "Windows Executable Files",
// 				Pattern:     "*.com;*.exe;*.bat;*.cmd;*.vbs;*.vbe;*.js;*.jse;*.wsf;*.wsh;*.msc",
// 			},
// 			{
// 				DisplayName: "All Files",
// 				Pattern:     "*",
// 			},
// 		},
// 		ShowHiddenFiles:            true,
// 		CanCreateDirectories:       true,
// 		ResolvesAliases:            false,
// 		TreatPackagesAsDirectories: false,
// 	})
// 	if selectedExecutable == "" {
// 		selectedExecutable = currentExecutable
// 	}
// 	return
// }

func Install(service Service) error {
	fi, err := os.Stat(service.Executable)
	if err == nil {
		if fi.Mode().IsDir() {
			return fmt.Errorf("Executable is a directory")
		}
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(service.Name)
	if err == nil {
		s.Close()
		return fmt.Errorf("Service %s already exist", service.Name)
	}

	s, err = m.CreateService(
		service.Name,
		service.Executable,
		mgr.Config{
			StartType:   service.StartMode.toSvcStartType(),
			DisplayName: service.DisplayName,
			Description: service.Description,
		},
		strings.Split(service.Argument, " ")...,
	)
	if err != nil {
		return fmt.Errorf("Install service failed: %s", err)
	}
	defer s.Close()

	if err = eventlog.InstallAsEventCreate(service.Name, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
		s.Delete()
		return fmt.Errorf("Install event log failed: %s", err)
	}
	return nil
}

func Update(service Service) (err error) {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(service.Name)
	if err != nil {
		s.Close()
		return fmt.Errorf("Could not access %s: %v", service.Name, err)
	}

	sConf, err := s.Config()
	if err != nil {
		return fmt.Errorf("Could not access config: %v", err)
	}

	sConf.StartType = service.StartMode.toSvcStartType()
	sConf.DisplayName = service.DisplayName
	sConf.Description = service.Description
	if err = s.UpdateConfig(sConf); err != nil {
		return fmt.Errorf("Update service failed: %v", err)
	}
	return
}

func Remove(name string) (err error) {
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

	if err = s.Delete(); err != nil {
		return fmt.Errorf("Failed to remove %s: %v", name, err)
	}

	if err = eventlog.Remove(name); err != nil {
		return fmt.Errorf("Remove event log failed: %s", err)
	}
	return
}
