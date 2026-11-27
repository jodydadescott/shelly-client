package util

import (
	"regexp"
	"strings"
)

// A Shelly access point is named "<Model>-<12 hex digit MAC>", for example "ShellyPlus1-A8032AB12345".
var shellySSID = regexp.MustCompile(`(?i)^shelly.*-([0-9a-f]{12})$`)

// DeviceIDForSSID returns the entry of ids (desired-config names) for the device whose access
// point is ssid: the one ending in "-<MAC>" with the MAC from the SSID, ignoring case.
func DeviceIDForSSID(ssid string, ids []string) (string, bool) {
	m := shellySSID.FindStringSubmatch(ssid)
	if m == nil {
		return "", false
	}
	suffix := "-" + strings.ToLower(m[1])
	for _, id := range ids {
		if strings.HasSuffix(strings.ToLower(id), suffix) {
			return id, true
		}
	}
	return "", false
}
