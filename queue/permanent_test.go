package queue

import (
	"errors"
	"fmt"
	"testing"
)

// =============================================================================
// Permanent Error Tests
// =============================================================================

func TestPermanent_NilReturnsNil(t *testing.T) {
	if Permanent(nil) != nil {
		t.Error("Permanent(nil) should return nil")
	}
}

func TestPermanent_MatchesErrPermanent(t *testing.T) {
	err := Permanent(errors.New("bad payload"))

	if !errors.Is(err, ErrPermanent) {
		t.Error("Permanent() error should match ErrPermanent via errors.Is")
	}
}

func TestPermanent_WrappedMatchesErrPermanent(t *testing.T) {
	err := fmt.Errorf("handler: %w", Permanent(errors.New("bad payload")))

	if !errors.Is(err, ErrPermanent) {
		t.Error("wrapped Permanent() error should match ErrPermanent via errors.Is")
	}
}

func TestPermanent_SentinelWrapMatches(t *testing.T) {
	err := fmt.Errorf("unknown type: %w", ErrPermanent)

	if !errors.Is(err, ErrPermanent) {
		t.Error("error wrapping ErrPermanent should match via errors.Is")
	}
}

func TestPermanent_PreservesMessage(t *testing.T) {
	inner := errors.New("bad payload")
	err := Permanent(inner)

	if err.Error() != "bad payload" {
		t.Errorf("Permanent().Error() = %q, want %q", err.Error(), "bad payload")
	}
}

func TestPermanent_Unwrap(t *testing.T) {
	inner := errors.New("bad payload")
	err := Permanent(inner)

	if !errors.Is(err, inner) {
		t.Error("Permanent() should unwrap to the inner error")
	}
}

func TestPermanent_OrdinaryErrorDoesNotMatch(t *testing.T) {
	err := errors.New("transient failure")

	if errors.Is(err, ErrPermanent) {
		t.Error("ordinary error should not match ErrPermanent")
	}
}
