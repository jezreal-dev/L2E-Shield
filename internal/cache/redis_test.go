package cache

import (
	"testing"
)

func TestNewCacheInvalidURL(t *testing.T) {
	_, err := New("invalid-url-format")
	if err == nil {
		t.Error("expected error for invalid redis url, got nil")
	}
}
