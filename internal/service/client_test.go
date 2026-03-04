package service

import (
	"testing"
)

func TestClientService_All(t *testing.T) {
	svc := NewClientService()

	// These will fail or return empty because we don't have adb/jps etc.
	// But we can verify they don't panic and we can mock the client package later if needed.
	_ = svc.LaunchChromium(":15500")
	_, _ = svc.ListJavaProcesses()
	_ = svc.InterceptJava("123", ":15500")
	_, _ = svc.ListAndroidDevices()
	_ = svc.InterceptAndroid("dev1", ":15500")
	_ = svc.InterceptAndroid("dev1", "localhost") // Test invalid address to hit port == "" fallback
	_ = svc.ClearAndroid("dev1", ":15500")
	_ = svc.ClearAndroid("dev1", "invalid")
	_ = svc.PushAndroidCert("dev1")

	script := svc.GetTerminalSetupScript(":15500", "")
	if script == "" {
		t.Error("Expected script")
	}
}
