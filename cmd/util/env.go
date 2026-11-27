package util

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var envRef = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// ExpandEnv replaces each ${NAME} in input with the environment variable NAME, so secrets such
// as Wi-Fi passwords can stay out of the config file. Only the ${NAME} form is expanded; a bare $
// is left alone. Every referenced variable must be set.
func ExpandEnv(input []byte) ([]byte, error) {
	missing := map[string]bool{}
	out := envRef.ReplaceAllFunc(input, func(m []byte) []byte {
		name := string(envRef.FindSubmatch(m)[1])
		v, ok := os.LookupEnv(name)
		if !ok {
			missing[name] = true
			return m
		}
		return []byte(v)
	})

	if len(missing) > 0 {
		var names []string
		for n := range missing {
			names = append(names, n)
		}
		sort.Strings(names)
		return nil, fmt.Errorf("config references unset environment variables: %s", strings.Join(names, ", "))
	}
	return out, nil
}
