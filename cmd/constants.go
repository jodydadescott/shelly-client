package cmd

import (
	"os"
)

const (
	FilePerm       = os.FileMode(0644)
	DirPerm        = os.FileMode(0755)
	ExePerm        = os.FileMode(0755)
	SecureFilePerm = os.FileMode(0400)

	BinaryName = "shelly-cli"

	DebugEnvVar = "DEBUG"

	ShellyConfigEnvVar   = "SHELLY_CONFIG"
	ShellyHostnameEnvVar = "SHELLY_HOST"
	ShellyPasswordEnvVar = "SHELLY_PASS"
	ShellyOutputEnvVar   = "SHELLY_OUTPUT"
	ShellyTimeoutEnvVar  = "SHELLY_TIMEOUT"
	ShellyURLEnvVar      = "SHELLY_URL"

	ShellyOutputDefault = "prettyjson"
)
