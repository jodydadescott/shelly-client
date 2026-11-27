package util

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Diff returns one line per JSON path where running and desired differ, sorted:
//
//	~ path: <running> -> <desired>   value differs
//	+ path: <desired>                set in desired only
//	- path: <running>                set on the device only
//
// Arrays are compared by index. Values are written as JSON.
func Diff(running, desired any) ([]string, error) {
	r, err := toJSONValue(running)
	if err != nil {
		return nil, fmt.Errorf("running config: %w", err)
	}
	d, err := toJSONValue(desired)
	if err != nil {
		return nil, fmt.Errorf("desired config: %w", err)
	}

	var lines []string
	diffValue("", r, d, &lines)
	sort.Strings(lines)
	return lines, nil
}

func toJSONValue(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	err = json.Unmarshal(b, &out)
	return out, err
}

func diffValue(path string, r, d any, lines *[]string) {
	switch rv := r.(type) {
	case map[string]any:
		if dv, ok := d.(map[string]any); ok {
			keys := map[string]bool{}
			for k := range rv {
				keys[k] = true
			}
			for k := range dv {
				keys[k] = true
			}
			for k := range keys {
				diffChild(join(path, k), rv, dv, k, lines)
			}
			return
		}
	case []any:
		if dv, ok := d.([]any); ok {
			for i := 0; i < max(len(rv), len(dv)); i++ {
				p := fmt.Sprintf("%s[%d]", path, i)
				switch {
				case i >= len(dv):
					*lines = append(*lines, "- "+p+": "+encode(rv[i]))
				case i >= len(rv):
					*lines = append(*lines, "+ "+p+": "+encode(dv[i]))
				default:
					diffValue(p, rv[i], dv[i], lines)
				}
			}
			return
		}
	}

	if encode(r) != encode(d) {
		*lines = append(*lines, "~ "+path+": "+encode(r)+" -> "+encode(d))
	}
}

func diffChild(path string, r, d map[string]any, k string, lines *[]string) {
	rv, inR := r[k]
	dv, inD := d[k]
	switch {
	case !inD:
		*lines = append(*lines, "- "+path+": "+encode(rv))
	case !inR:
		*lines = append(*lines, "+ "+path+": "+encode(dv))
	default:
		diffValue(path, rv, dv, lines)
	}
}

func join(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}

func encode(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
