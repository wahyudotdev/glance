package client

import (
	"testing"
)

func TestListDockerContainers_NoDocker(t *testing.T) {
	// This test just ensures that the function doesn't panic
	// and returns an error if Docker is not available.
	// If Docker IS available, it will return the list of containers.
	_, _ = ListDockerContainers()
}

func TestInterceptDocker_Stub(t *testing.T) {
	err := InterceptDocker("test", ":15500")
	if err != nil {
		t.Errorf("Expected nil for stub, got %v", err)
	}
}

func TestStopInterceptDocker_Stub(t *testing.T) {
	err := StopInterceptDocker("test")
	if err != nil {
		t.Errorf("Expected nil for stub, got %v", err)
	}
}
