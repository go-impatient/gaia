package str

import (
	"strings"
	"unicode"

	"github.com/satori/go.uuid"
)

// IsBlank checks if a string is whitespace or empty ("").
// goutils.IsBlank("")        = true
// goutils.IsBlank(" ")       = true
// goutils.IsBlank("bob")     = false
// goutils.IsBlank("  bob  ") = false
func IsBlank(str string) bool {
	strLen := len(str)
	if str == "" || strLen == 0 {
		return true
	}
	for i := 0; i < strLen; i++ {
		if unicode.IsSpace(rune(str[i])) == false {
			return false
		}
	}
	return true
}

// IsNotBlank
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// DefaultIfBlank
func DefaultIfBlank(str, def string) string {
	if IsBlank(str) {
		return def
	} else {
		return str
	}
}

// IsEmpty checks if a string is empty (""). Returns true if empty, and false otherwise.
func IsEmpty(str string) bool {
	return len(str) == 0
}

// IsNotEmpty
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// Substr 截取字符串
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	if start > len(bt) {
		start = start % len(bt)
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt)
	} else {
		end = start + length
	}
	return string(bt[start:end])
}

// UUID
func UUID() string {
	u := uuid.NewV4()
	return strings.ReplaceAll(u.String(), "-", "")
}

// Equals
func Equals(a, b string) bool {
	return a == b
}

// EqualsIgnoreCase
func EqualsIgnoreCase(a, b string) bool {
	return a == b || strings.ToUpper(a) == strings.ToUpper(b)
}

// RuneLen 字符成长度
func RuneLen(s string) int {
	bt := []rune(s)
	return len(bt)
}

// GetSummary 获取summary
func GetSummary(s string, length int) string {
	s = strings.TrimSpace(s)
	summary := Substr(s, 0, length)
	if RuneLen(s) > length {
		summary += "..."
	}
	return summary
}
