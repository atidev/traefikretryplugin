package traefikretryplugin

import (
	"fmt"
	. "github.com/niki-timofe/traefikretryplugin/pkg/ranges"
	. "github.com/niki-timofe/traefikretryplugin/pkg/structuredheaders"
)

type RetryPolicy struct {
	codes    Range
	attempts int
}

func (p *RetryPolicy) Applicable(status int) bool {
	return p.codes.Includes(status)
}

func (p *RetryPolicy) CanRetry(attempt int) bool {
	return attempt < p.attempts
}

func (p *RetryPolicy) String() string {
	return fmt.Sprintf("Policy: codes: %s, attempts: %d", p.codes.String(), p.attempts)
}

func ParsePolicy(hp map[string]ListItem) (*RetryPolicy, error) {
	codes, err := parseCodes(hp)
	if err != nil {
		return nil, fmt.Errorf("traefikretryplugin.ParsePolicy: can't parse codes: %w", err)
	}

	attempts, err := parseAttempts(hp)
	if err != nil {
		return nil, fmt.Errorf("traefikretryplugin.ParsePolicy: can't parse attempts: %w", err)
	}

	return &RetryPolicy{
		codes:    codes,
		attempts: attempts,
	}, nil
}

func parseAttempts(hp map[string]ListItem) (int, error) {
	a, err := hp["attempts"].Item()
	if err != nil {
		return 0, fmt.Errorf("traefikretryplugin.parseAttempts: can't parse item: %w", err)
	}

	as, err := a.Number()
	if err != nil {
		return 0, fmt.Errorf("traefikretryplugin.parseAttempts: can't parse number: %w", err)
	}

	attempts, err := as.Integer()
	if err != nil {
		return 0, fmt.Errorf("traefikretryplugin.parseAttempts: can't parse integer %w", err)
	}

	return attempts, nil
}

func parseCodes(hp map[string]ListItem) (Range, error) {
	c, err := hp["codes"].Item()
	if err != nil {
		return nil, fmt.Errorf("traefikretryplugin.parseCodes: can't parse item: %w", err)
	}

	cs, err := c.Str()
	if err != nil {
		return nil, fmt.Errorf("traefikretryplugin.parseCodes can't parse string: %w", err)
	}

	codes, err := NewRange(cs)
	if err != nil {
		return nil, fmt.Errorf("traefikretryplugin.parseCodes: can't parse range: %w", err)
	}

	return codes, nil
}
