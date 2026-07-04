// Package security provides utilities for inspecting and filtering malicious payloads.
package security

import (
	"fmt"
	"strings"
)

var injectionPatterns = []string{
	"ignore previous directions",
	"ignore previous instructions",
	"ignore system directions",
	"ignore system instructions",
	"ignore the above instructions",
	"system override",
	"you are now a",
}

// ValidatePrompt inspects the provided text against known jailbreak patterns
// and oversized payload thresholds, returning an error if the payload is malicious.
func ValidatePrompt(prompt string) error {
	if len(prompt) > 10000 {
		return fmt.Errorf("prompt payload exceeds size limit of 10000 characters")
	}

	lowerPrompt := strings.ToLower(prompt)
	for _, pattern := range injectionPatterns {
		if strings.Contains(lowerPrompt, pattern) {
			return fmt.Errorf("prompt injection attempt detected: contains blocked pattern %q", pattern)
		}
	}

	return nil
}
