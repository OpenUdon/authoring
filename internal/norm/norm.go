package norm

import (
	"strings"
)

// Token returns Authoring's canonical normalized token form.
func Token(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "-", "_")
	return strings.Join(strings.Fields(value), "_")
}

// Trim returns value without surrounding whitespace.
func Trim(value string) string {
	return strings.TrimSpace(value)
}

// FirstNonEmpty returns the first trimmed non-empty value.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

// Metadata returns a copy of in with trimmed non-empty keys and values.
func Metadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		out[key] = value
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// CompareStrings compares paired string values in order.
func CompareStrings(values ...string) int {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] < values[i+1] {
			return -1
		}
		if values[i] > values[i+1] {
			return 1
		}
	}
	return 0
}

// CompareSeverity compares Authoring severity strings in readiness order.
func CompareSeverity(a, b string) int {
	return SeverityRank(a) - SeverityRank(b)
}

// SeverityRank returns the readiness sort rank for severity.
func SeverityRank(severity string) int {
	switch Token(severity) {
	case "", "error", "blocking":
		return 0
	case "warning":
		return 1
	case "advisory":
		return 2
	case "info":
		return 3
	default:
		return 0
	}
}
