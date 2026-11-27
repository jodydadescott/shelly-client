package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestToAuthResponse(t *testing.T) {
	req := &AuthRequest{
		AuthType: "digest", Nonce: 140976, NonceCount: 1, Realm: "shellypluswdus-3c8a1fec368c",
		Algorithm: "SHA-256", Username: "admin", Password: "secret",
	}
	auth, err := req.ToAuthResponse()
	if err != nil {
		t.Fatal(err)
	}

	ha1 := getSHA256("admin:shellypluswdus-3c8a1fec368c:secret")
	ha2 := getSHA256("dummy_method:dummy_uri")
	want := getSHA256(fmt.Sprintf("%s:140976:1:%d:auth:%s", ha1, auth.Cnonce, ha2))
	if auth.Response != want {
		t.Errorf("response = %s, want %s", auth.Response, want)
	}

	// Shelly's auth object takes nonce and cnonce as JSON numbers.
	b, err := json.Marshal(auth)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(b, &fields); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"nonce", "cnonce"} {
		if _, ok := fields[k].(float64); !ok {
			t.Errorf("%s is %T in %s, want a number", k, fields[k], b)
		}
	}
}

func TestRequestAuthKey(t *testing.T) {
	id, method := 1, "Shelly.GetConfig"
	b, err := json.Marshal(&Request{Auth: &AuthResponse{Realm: "r"}, ID: &id, Method: &method})
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(b, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["auth"]; !ok {
		t.Errorf("no lower-case auth key in %s", b)
	}
	b, _ = json.Marshal(&Request{ID: &id, Method: &method})
	if string(b) != `{"id":1,"method":"Shelly.GetConfig"}` {
		t.Errorf("request without auth = %s", b)
	}
}
