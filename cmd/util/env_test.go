package util

import "testing"

func TestExpandEnv(t *testing.T) {
	t.Setenv("SHELLY_TEST_WIFI_PASS", "s3cr$t")

	got, err := ExpandEnv([]byte(`pass: "${SHELLY_TEST_WIFI_PASS}", literal: "a$b $X"`))
	if err != nil {
		t.Fatal(err)
	}
	if want := `pass: "s3cr$t", literal: "a$b $X"`; string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}

	if _, err := ExpandEnv([]byte(`${SHELLY_TEST_UNSET_B} ${SHELLY_TEST_UNSET_A}`)); err == nil ||
		err.Error() != "config references unset environment variables: SHELLY_TEST_UNSET_A, SHELLY_TEST_UNSET_B" {
		t.Errorf("got %v", err)
	}
}
