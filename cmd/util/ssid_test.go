package util

import "testing"

func TestDeviceIDForSSID(t *testing.T) {
	ids := []string{"common", "S1G3", "shelly1g3-a8032ab12345", "shellyplus1-00112233aabb"}

	tests := []struct {
		ssid   string
		want   string
		wantOK bool
	}{
		{"Shelly1Gen3-A8032AB12345", "shelly1g3-a8032ab12345", true},
		{"ShellyPlus1-00112233AABB", "shellyplus1-00112233aabb", true},
		{"Shelly1Gen3-FFFFFFFFFFFF", "", false}, // no config for this device
		{"HomeWiFi", "", false},
		{"NotShelly-A8032AB12345", "", false},
		{"Shelly1Gen3-A8032AB1234", "", false}, // 11 hex digits
	}

	for _, tt := range tests {
		got, ok := DeviceIDForSSID(tt.ssid, ids)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("%s: got %q %v, want %q %v", tt.ssid, got, ok, tt.want, tt.wantOK)
		}
	}
}
