// Package devconfig holds a Shelly Gen2+ device config in the device's own shape and compares
// a desired config against the running one.
//
// A Config maps a component key, exactly as Shelly.GetConfig returns it ("sys", "wifi",
// "switch:0", "cover:0", ...), to that component's config object. The extra key "auth" holds
// {"enable": bool}; it is not part of Shelly.GetConfig and is applied with Shelly.SetAuth.
package devconfig

import (
	"encoding/json"
	"fmt"
	"sort"
)

// AuthKey is the Config key for device authentication.
const AuthKey = "auth"

// Config is a device config: component key -> component config.
type Config map[string]any

// UnmarshalYAML decodes YAML into JSON-compatible values (string keys, float64 numbers), so a
// Config read from YAML compares equal to the same config read from a device.
func (t *Config) UnmarshalYAML(unmarshal func(any) error) error {
	var raw map[string]any
	if err := unmarshal(&raw); err != nil {
		return err
	}
	v, err := normalize(raw)
	if err != nil {
		return err
	}
	*t = Config(v.(map[string]any))
	return nil
}

// FromJSON decodes a JSON object into a Config.
func FromJSON(data []byte) (Config, error) {
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return c, nil
}

// Clone returns a deep copy.
func (t Config) Clone() Config {
	if t == nil {
		return nil
	}
	return Config(cloneValue(map[string]any(t)).(map[string]any))
}

// Merge returns base with override layered on top: objects merge key by key, any other value in
// override replaces the one in base. Neither input is modified.
func Merge(base, override Config) Config {
	if base == nil {
		return override.Clone()
	}
	return Config(mergeValue(map[string]any(base), map[string]any(override)).(map[string]any))
}

// Result is the outcome of Compare.
type Result struct {
	// Drift has one line per desired value that differs on the device: "~ path: running -> desired".
	Drift []string
	// Unverified has one line per desired key the device does not report, such as write-only
	// passwords: "? path". They are sent when their component is set but never count as drift.
	Unverified []string
	// Changed lists the component keys that have drift, sorted.
	Changed []string
}

// Equal reports whether there is no drift.
func (t *Result) Equal() bool {
	return len(t.Drift) == 0
}

// Compare checks every value in desired against running. Keys that are only in running are
// ignored, so desired needs to name only what it manages. A desired component the device does
// not have at all is an error, not drift.
func Compare(running, desired Config) (*Result, error) {
	r := &Result{}

	for _, key := range sortedKeys(desired) {
		rv, ok := running[key]
		if !ok {
			return nil, fmt.Errorf("device has no component %q", key)
		}
		before := len(r.Drift)
		compareValue(key, rv, desired[key], r)
		if len(r.Drift) > before {
			r.Changed = append(r.Changed, key)
		}
	}

	sort.Strings(r.Drift)
	sort.Strings(r.Unverified)
	return r, nil
}

func compareValue(path string, running, desired any, r *Result) {
	if dm, ok := desired.(map[string]any); ok {
		if rm, ok := running.(map[string]any); ok {
			for _, k := range sortedKeys(dm) {
				p := path + "." + k
				rv, ok := rm[k]
				if !ok {
					r.Unverified = append(r.Unverified, "? "+p)
					continue
				}
				compareValue(p, rv, dm[k], r)
			}
			return
		}
	}

	if encode(running) != encode(desired) {
		r.Drift = append(r.Drift, "~ "+path+": "+encode(running)+" -> "+encode(desired))
	}
}

func mergeValue(base, override any) any {
	bm, bok := base.(map[string]any)
	om, ook := override.(map[string]any)
	if !bok || !ook {
		return cloneValue(override)
	}
	out := cloneValue(bm).(map[string]any)
	for k, v := range om {
		out[k] = mergeValue(out[k], v)
	}
	return out
}

func cloneValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			out[k] = cloneValue(e)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = cloneValue(e)
		}
		return out
	default:
		return v
	}
}

// normalize converts YAML-decoded values (map[any]any, int) into JSON-decoded ones
// (map[string]any, float64).
func normalize(v any) (any, error) {
	v, err := stringKeys(v)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	err = json.Unmarshal(b, &out)
	return out, err
}

func stringKeys(v any) (any, error) {
	switch x := v.(type) {
	case map[any]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			ks, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("key %v is not a string", k)
			}
			ev, err := stringKeys(e)
			if err != nil {
				return nil, err
			}
			out[ks] = ev
		}
		return out, nil
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			ev, err := stringKeys(e)
			if err != nil {
				return nil, err
			}
			out[k] = ev
		}
		return out, nil
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			ev, err := stringKeys(e)
			if err != nil {
				return nil, err
			}
			out[i] = ev
		}
		return out, nil
	default:
		return v, nil
	}
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func encode(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
