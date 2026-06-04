package security

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	appauth "github.com/viictormotorhead/financial-app/internal/application/auth"
)

func TestJWTTokenService_IssueIncludesUsername(t *testing.T) {
	svc := NewJWTTokenService("test-secret", time.Hour)

	token, err := svc.Issue(context.Background(), appauth.TokenClaims{
		UserID:   "RMWH96py12OJCUL4TRYjk",
		Username: "victor",
	})
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 JWT parts, got %d", len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode payload: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if decoded["sub"] != "RMWH96py12OJCUL4TRYjk" {
		t.Fatalf("sub: got %v", decoded["sub"])
	}
	if decoded["username"] != "victor" {
		t.Fatalf("username missing or wrong: got %v", decoded)
	}

	parsed, err := svc.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.UserID != "RMWH96py12OJCUL4TRYjk" || parsed.Username != "victor" {
		t.Fatalf("parsed claims: %+v", parsed)
	}
}
