// Package service implements the core business logic.
package service

import (
	"glance/internal/client"
	"glance/internal/model"
	"net"

	"github.com/elazarl/goproxy"
)

// ClientService defines the interface for interacting with external clients and devices.
type ClientService interface {
	LaunchChromium(proxyAddr string) error
	ListJavaProcesses() ([]model.JavaProcess, error)
	InterceptJava(pid string, proxyAddr string) error
	GetTerminalSetupScript(proxyAddr string) string
	ListAndroidDevices() ([]model.AndroidDevice, error)
	InterceptAndroid(deviceID string, proxyAddr string) error
	ClearAndroid(deviceID string, proxyAddr string) error
	PushAndroidCert(deviceID string) error
}

type clientService struct{}

// NewClientService creates a new ClientService.
func NewClientService() ClientService {
	return &clientService{}
}

func (s *clientService) LaunchChromium(proxyAddr string) error {
	return client.LaunchChromium(proxyAddr)
}

func (s *clientService) ListJavaProcesses() ([]model.JavaProcess, error) {
	return client.ListJavaProcesses()
}

func (s *clientService) InterceptJava(pid string, proxyAddr string) error {
	return client.BuildAndAttachAgent(pid, proxyAddr)
}

func (s *clientService) GetTerminalSetupScript(proxyAddr string) string {
	return client.GetTerminalSetupScript(proxyAddr)
}

func (s *clientService) ListAndroidDevices() ([]model.AndroidDevice, error) {
	return client.ListAndroidDevices()
}

func (s *clientService) InterceptAndroid(deviceID string, proxyAddr string) error {
	_, port, _ := net.SplitHostPort(proxyAddr)
	if port == "" {
		port = "8000"
	}
	return client.ConfigureAndroidProxy(deviceID, port)
}

func (s *clientService) ClearAndroid(deviceID string, proxyAddr string) error {
	_, port, _ := net.SplitHostPort(proxyAddr)
	if port == "" {
		port = "8000"
	}
	return client.ClearAndroidProxy(deviceID, port)
}

func (s *clientService) PushAndroidCert(deviceID string) error {
	return client.PushCertToDevice(deviceID, []byte(goproxy.CA_CERT))
}
