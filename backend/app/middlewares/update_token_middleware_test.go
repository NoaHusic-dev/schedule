package middlewares

import "testing"

func TestBearerTokenMatches(t *testing.T) {
	tests := []struct {
		name   string
		header string
		token  string
		want   bool
	}{
		{"matching token", "Bearer s3cret", "s3cret", true},
		{"wrong token", "Bearer guess", "s3cret", false},
		{"no header", "", "s3cret", false},
		{"token without the Bearer scheme", "s3cret", "s3cret", false},
		{"a prefix of the token", "Bearer s3c", "s3cret", false},
		{"no token configured rejects everything", "Bearer ", "", false},
		{"no token configured rejects a missing header", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bearerTokenMatches(tt.header, []byte(tt.token)); got != tt.want {
				t.Errorf("bearerTokenMatches(%q, %q) = %v, want %v", tt.header, tt.token, got, tt.want)
			}
		})
	}
}
