package iaas_clients

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cloudfoundry-incubator/stembuild/iaas_cli"
)

type VcenterClient struct {
	Username      string
	Password      string
	Url           string
	credentialUrl string
	Runner        iaas_cli.CliRunner
}

func NewVcenterClient(username string, password string, url string, runner iaas_cli.CliRunner) *VcenterClient {
	urlWithCredentials := fmt.Sprintf("%s:%s@%s", username, password, url)
	return &VcenterClient{Username: username, Password: password, Url: url, credentialUrl: urlWithCredentials, Runner: runner}
}

func (c VcenterClient) ValidateUrl() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.Url})
	if errCode != 0 {
		return errors.New("invalid url")
	}
	return nil

}

func (c VcenterClient) ValidateCredentials() error {
	errCode := c.Runner.Run([]string{"about", "-u", c.credentialUrl})
	if errCode != 0 {
		return errors.New("invalid credentials")
	}

	return nil
}

func (c VcenterClient) FindVM(vmInventoryPath string) error {
	errCode := c.Runner.Run([]string{"find", "-maxdepth=0", "-u", c.credentialUrl, vmInventoryPath})
	if errCode != 0 {
		errorMsg := "invalid VM path"
		return errors.New(errorMsg)
	}

	return nil
}

func (c VcenterClient) ListDevices(vmInventoryPath string) ([]string, error) {
	o, exitCode, err := c.Runner.RunWithOutput([]string{"device.ls", "-vm", vmInventoryPath})

	if exitCode != 0 {
		return []string{}, fmt.Errorf("failed to list devices in vCenter, govc exit code %d", exitCode)
	}

	if err != nil {
		return []string{}, fmt.Errorf("failed to parse list of devices. Err: %s", err)
	}

	entries := strings.Split(o, "\n")
	devices := []string{}
	r, _ := regexp.Compile(`\S+`)
	for _, entry := range entries {
		if entry != "" {
			devices = append(devices, r.FindString(entry))
		}
	}
	return devices, nil
}
func (c VcenterClient) RemoveDevice(vmInventoryPath string, deviceName string) error {
	errCode := c.Runner.Run([]string{"device.remove", "-u", c.credentialUrl, "-vm", vmInventoryPath, deviceName})
	if errCode != 0 {
		return fmt.Errorf("%s could not be removed/not found", deviceName)
	}
	return nil
}

func (c VcenterClient) EjectCDRom(vmInventoryPath string, deviceName string) error {

	errCode := c.Runner.Run([]string{"device.cdrom.eject", "-u", c.credentialUrl, "-vm", vmInventoryPath, "-device", deviceName})
	if errCode != 0 {
		return fmt.Errorf("%s could not be ejected", deviceName)
	}
	return nil
}

func (c VcenterClient) ExportVM(vmInventoryPath string, destination string) error {
	_, err := os.Stat(destination)
	if err != nil {
		return errors.New(fmt.Sprintf("provided destination: %s does not exist", destination))
	}
	errCode := c.Runner.Run([]string{"export.ovf", "-u", c.credentialUrl, "-sha", "1", "-vm", vmInventoryPath, destination})
	if errCode != 0 {
		errorMsg := fmt.Sprintf(vmInventoryPath + " could not be exported")
		return errors.New(errorMsg)
	}
	return nil
}
