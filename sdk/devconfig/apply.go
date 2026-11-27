package devconfig

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ApplyOrder returns the component keys to set, in the order to set them: sorted, with "wifi"
// last because changing it can drop the connection. AuthKey is left out; it is set separately,
// after everything else.
func ApplyOrder(keys []string) []string {
	var out []string
	wifi := false
	for _, k := range keys {
		switch k {
		case AuthKey:
		case "wifi":
			wifi = true
		default:
			out = append(out, k)
		}
	}
	sort.Strings(out)
	if wifi {
		out = append(out, "wifi")
	}
	return out
}

// SetConfigCall returns the RPC method and params that set component key to config. The method
// name is taken from methods (the device's Shelly.ListMethods), so its exact spelling ("WiFi",
// "BLE", "PM1", ...) comes from the device.
func SetConfigCall(key string, config any, methods []string) (string, map[string]any, error) {
	component, idPart, hasID := strings.Cut(key, ":")

	want := component + ".SetConfig"
	method := ""
	for _, m := range methods {
		if strings.EqualFold(m, want) {
			method = m
			break
		}
	}
	if method == "" {
		return "", nil, fmt.Errorf("device has no method %s for component %q", want, key)
	}

	params := map[string]any{"config": config}
	if hasID {
		id, err := strconv.Atoi(idPart)
		if err != nil {
			return "", nil, fmt.Errorf("component %q: id %q is not a number", key, idPart)
		}
		params["id"] = id
	}
	return method, params, nil
}
