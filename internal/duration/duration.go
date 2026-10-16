package duration

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Define day
const Day = 24 * time.Hour

// Parses a duration composed of integer values and s, m, h, or d units.
func ParseDuration(input string) (time.Duration, error) {
	// Check that the duration exists
	if input == "" {
		return 0, fmt.Errorf("duration cannot be empty")
	}

	// Define max duration
	const maxDuration = uint64(1<<63 - 1)

	// Define a variable to store the final total duration
	var total uint64

	for i := 0; i < len(input); {
		start := i
		// Check for mentions of s, m, h or d, otherwise skip iteration until found
		for i < len(input) && input[i] >= '0' && input[i] <= '9' {
			i++
		}
		// If no mentions found then the input is invalid
		if start == i || i == len(input) {
			return 0, fmt.Errorf("invalid duration %q", input)
		}

		// Parse the integer values in the string as unsigned integers
		value, err := strconv.ParseUint(input[start:i], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration %q: %w", input, err)
		}

		var unit time.Duration
		// Check that mentions of a non-integer value correspond with a time
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

		// Check whether duration is within max duaration value
		unitValue := uint64(unit)
		if value > maxDuration/unitValue {
			return 0, fmt.Errorf("duration %q overflows time.Duration", input)
		}
		part := value * unitValue
		if total > maxDuration-part {
			return 0, fmt.Errorf("duration %q overflows time.Duration", input)
		}

		// Add to total duration
		total += part
	}

	return time.Duration(total), nil
}

// Converts given time into a more human readable format using s, m, h and d
func Format(d time.Duration) string {
	// If the duration is 0, just return 0s
	if d <= 0 {
		return "0s"
	}

	// Define s, m, h and d to their corresponding units of time
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
	// For each unit of time, check if the duration given can fit into it
	// If so, put the duration into that unit as many times as possible before moving onto the next
	// Else move onto the next unit
	for _, unit := range units {
		if d >= unit.size {
			parts = append(parts, fmt.Sprintf("%d%s", d/unit.size, unit.name))
			d %= unit.size
		}
	}

	return strings.Join(parts, " ")
}
