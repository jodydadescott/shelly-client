package devconfig

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func mustJSON(t *testing.T, s string) Config {
	t.Helper()
	c, err := FromJSON([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestCompare(t *testing.T) {
	running := mustJSON(t, `{
		"sys": {"device": {"name": "porch", "eco_mode": false}, "sntp": {"server": "time.google.com"}},
		"wifi": {"ap": {"enable": true}, "sta": {"ssid": "ShellyAP", "is_open": true, "enable": false}},
		"switch:0": {"id": 0, "name": null, "in_mode": "follow"},
		"auth": {"enable": false}
	}`)

	t.Run("equal when desired is a subset", func(t *testing.T) {
		r, err := Compare(running, mustJSON(t, `{"sys": {"device": {"name": "porch"}}, "switch:0": {"in_mode": "follow"}}`))
		if err != nil {
			t.Fatal(err)
		}
		if !r.Equal() || r.Changed != nil || r.Unverified != nil {
			t.Errorf("got %+v, want equal", r)
		}
	})

	t.Run("drift, unverified and changed components", func(t *testing.T) {
		r, err := Compare(running, mustJSON(t, `{
			"wifi": {"ap": {"enable": false}, "sta": {"ssid": "home", "pass": "secret", "enable": true}},
			"switch:0": {"name": "Porch light"},
			"sys": {"device": {"name": "porch"}},
			"auth": {"enable": true}
		}`))
		if err != nil {
			t.Fatal(err)
		}
		wantDrift := []string{
			`~ auth.enable: false -> true`,
			`~ switch:0.name: null -> "Porch light"`,
			`~ wifi.ap.enable: true -> false`,
			`~ wifi.sta.enable: false -> true`,
			`~ wifi.sta.ssid: "ShellyAP" -> "home"`,
		}
		if !reflect.DeepEqual(r.Drift, wantDrift) {
			t.Errorf("drift\ngot  %q\nwant %q", r.Drift, wantDrift)
		}
		if want := []string{"? wifi.sta.pass"}; !reflect.DeepEqual(r.Unverified, want) {
			t.Errorf("unverified got %q, want %q", r.Unverified, want)
		}
		if want := []string{"auth", "switch:0", "wifi"}; !reflect.DeepEqual(r.Changed, want) {
			t.Errorf("changed got %q, want %q", r.Changed, want)
		}
	})

	t.Run("unverified alone is not drift", func(t *testing.T) {
		r, err := Compare(running, mustJSON(t, `{"wifi": {"sta": {"pass": "secret"}}}`))
		if err != nil {
			t.Fatal(err)
		}
		if !r.Equal() || r.Changed != nil {
			t.Errorf("got %+v, want equal", r)
		}
	})

	t.Run("component missing on device is an error", func(t *testing.T) {
		if _, err := Compare(running, mustJSON(t, `{"cover:0": {"name": "blind"}}`)); err == nil {
			t.Error("want error")
		}
	})
}

func TestYAMLComparesEqualToJSON(t *testing.T) {
	var desired Config
	err := yaml.Unmarshal([]byte("sys:\n  sntp:\n    server: time.google.com\nswitch:0:\n  id: 0\n  auto_off_delay: 60\n"), &desired)
	if err != nil {
		t.Fatal(err)
	}
	running := mustJSON(t, `{"sys": {"sntp": {"server": "time.google.com"}}, "switch:0": {"id": 0, "auto_off_delay": 60.0}}`)
	r, err := Compare(running, desired)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Equal() {
		t.Errorf("drift %q", r.Drift)
	}
}

func TestMerge(t *testing.T) {
	common := mustJSON(t, `{"sys": {"sntp": {"server": "time.google.com"}, "device": {"eco_mode": true}}, "cloud": {"enable": false}}`)
	device := mustJSON(t, `{"sys": {"device": {"name": "porch"}}, "cloud": {"enable": true}}`)

	got := Merge(common, device)
	want := mustJSON(t, `{"sys": {"sntp": {"server": "time.google.com"}, "device": {"eco_mode": true, "name": "porch"}}, "cloud": {"enable": true}}`)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v\nwant %v", got, want)
	}
	if _, ok := common["sys"].(map[string]any)["device"].(map[string]any)["name"]; ok {
		t.Error("Merge modified its base")
	}
}

func TestApplyOrder(t *testing.T) {
	got := ApplyOrder([]string{"wifi", "auth", "switch:0", "cloud", "sys"})
	want := []string{"cloud", "switch:0", "sys", "wifi"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetConfigCall(t *testing.T) {
	methods := []string{"Shelly.GetConfig", "WiFi.SetConfig", "Switch.SetConfig", "Sys.SetConfig"}

	method, params, err := SetConfigCall("wifi", map[string]any{"ap": map[string]any{"enable": false}}, methods)
	if err != nil || method != "WiFi.SetConfig" || params["id"] != nil {
		t.Errorf("wifi: got %q %v %v", method, params, err)
	}

	method, params, err = SetConfigCall("switch:1", map[string]any{"name": "x"}, methods)
	if err != nil || method != "Switch.SetConfig" || params["id"] != 1 {
		t.Errorf("switch:1: got %q %v %v", method, params, err)
	}

	if _, _, err := SetConfigCall("cover:0", nil, methods); err == nil {
		t.Error("cover:0: want error, device has no Cover.SetConfig")
	}
}
