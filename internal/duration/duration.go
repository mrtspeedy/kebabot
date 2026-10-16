package duration

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const Day = 24 * time.Hour

// ParseDuration parses a duration composed of integer values and s, m, h, or d units.
func ParseDuration(input string) (time.Duration, error) {
	if input == "" {
		return 0, fmt.Errorf("duration cannot be empty")
	}

	const maxDuration = uint64(1<<63 - 1)
	var total uint64

	for i := 0; i < len(input); {
		start := i
		for i < len(input) && input[i] >= '0' && input[i] <= '9' {
			i++
		}
		if start == i || i == len(input) {
			return 0, fmt.Errorf("invalid duration %q", input)
		}

		value, err := strconv.ParseUint(input[start:i], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", input, err)
		}

		var unit time.Duration
		switch input[i] {
		case 's':
			unit = time.Second
		case 'm':
			unit = time.Minute
		case 'h':
			unit = time.Hour
		case 'd':
			unit = Day
		default:
			return 0, fmt.Errorf("unsupported duration unit %q", input[i])
		}
		i++

		unitValue := uint64(unit)
		if value > maxDuration/unitValue {
			return 0, fmt.Errorf("duration %q overflows time.Duration", input)
		}
		part := value * unitValue
		if total > maxDuration-part {
			return 0, fmt.Errorf("duration %q overflows time.Duration", input)
		}
		total += part
	}

	return time.Duration(total), nil
}

// Format converts a positive duration to a compact string using d, h, m, and s.
func Format(d time.Duration) string {
	if d <= 0 {
		return "0s"
	}

	units := []struct {
		size time.Duration
		name string
	}{
		{Day, "d"},
		{time.Hour, "h"},
		{time.Minute, "m"},
		{time.Second, "s"},
	}

	var parts []string
	for _, unit := range units {
		if d >= unit.size {
			parts = append(parts, fmt.Sprintf("%d%s", d/unit.size, unit.name))
			d %= unit.size
		}
	}

	return strings.Join(parts, " ")
}
