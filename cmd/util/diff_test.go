package util

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	type sta struct {
		SSID   *string `json:"ssid,omitempty"`
		Enable *bool   `json:"enable,omitempty"`
	}
	type wifi struct {
		Sta *sta `json:"sta,omitempty"`
	}
	type sw struct {
		ID   *int    `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}
	type config struct {
		Wifi   *wifi `json:"wifi,omitempty"`
		Switch []*sw `json:"switch,omitempty"`
	}
	str := func(s string) *string { return &s }
	b := func(v bool) *bool { return &v }
	n := func(v int) *int { return &v }

	tests := []struct {
		name             string
		running, desired any
		want             []string
	}{
		{
			name:    "equal",
			running: &config{Wifi: &wifi{Sta: &sta{SSID: str("home"), Enable: b(true)}}},
			desired: &config{Wifi: &wifi{Sta: &sta{SSID: str("home"), Enable: b(true)}}},
			want:    nil,
		},
		{
			name:    "changed, added and removed fields",
			running: &config{Wifi: &wifi{Sta: &sta{SSID: str("ShellyAP"), Enable: b(false)}}},
			desired: &config{Wifi: &wifi{Sta: &sta{SSID: str("home")}}, Switch: []*sw{{ID: n(0), Name: str("porch")}}},
			want: []string{
				`+ switch: [{"id":0,"name":"porch"}]`,
				`- wifi.sta.enable: false`,
				`~ wifi.sta.ssid: "ShellyAP" -> "home"`,
			},
		},
		{
			name:    "arrays compared by index",
			running: &config{Switch: []*sw{{ID: n(0), Name: str("a")}, {ID: n(1), Name: str("b")}}},
			desired: &config{Switch: []*sw{{ID: n(0), Name: str("x")}}},
			want: []string{
				`- switch[1]: {"id":1,"name":"b"}`,
				`~ switch[0].name: "a" -> "x"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Diff(tt.running, tt.desired)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}
