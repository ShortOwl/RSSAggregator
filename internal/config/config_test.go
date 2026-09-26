package config

import "testing"

func TestParseRetentionDays(t *testing.T) {
	for _, tt := range []struct {
		value   string
		want    int32
		invalid bool
	}{
		{"", 50, false}, {"30", 30, false}, {"1", 1, false},
		{"0", 0, true}, {"-1", 0, true}, {"90.5", 0, true},
		{"invalid", 0, true}, {"2147483648", 0, true},
	} {
		t.Run(tt.value, func(t *testing.T) {
			got, err := parseRetentionDays(tt.value)
			if (err != nil) != tt.invalid || got != tt.want {
				t.Fatalf("got (%d, %v), want (%d, invalid=%v)", got, err, tt.want, tt.invalid)
			}
		})
	}
}

func TestLoadRetentionDays(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DB_URL", "postgres://unused")
	t.Setenv("JWT_SECRET", "test-secret-at-least-32-characters-long")
	t.Setenv("POST_RETENTION_DAYS", "30")
	if got := Load().PostRetentionDays; got != 30 {
		t.Fatalf("PostRetentionDays = %d, want 30 from POST_RETENTION_DAYS", got)
	}
}
