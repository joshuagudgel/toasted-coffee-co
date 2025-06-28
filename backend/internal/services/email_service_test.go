package services

import (
	"testing"
)

func TestSanitizeInput(t *testing.T) {
	emailService := NewEmailService()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Plain text passes through unchanged",
			input:    "Hello, this is a normal message.",
			expected: "Hello, this is a normal message.",
		},
		{
			name:     "Basic HTML tags are removed",
			input:    "<b>Bold text</b> and <i>italic text</i>",
			expected: "Bold text and italic text",
		},
		{
			name:     "Script tags are removed completely",
			input:    "Before <script>alert('XSS');</script> After",
			expected: "Before  After",
		},
		{
			name:     "Malicious attributes are removed",
			input:    "<div onmouseover=\"alert('XSS')\">Hover me</div>",
			expected: "Hover me",
		},
		{
			name:     "URL with javascript protocol is sanitized",
			input:    "<a href=\"javascript:alert('XSS')\">Click me</a>",
			expected: "Click me",
		},
		{
			name:     "Complex nested payload is sanitized",
			input:    "<div><script>document.write('<img src=\"x\" onerror=\"alert(1)\">')</script></div>",
			expected: "",
		},
		{
			name:     "Handles HTML entities",
			input:    "&lt;script&gt;alert('XSS');&lt;/script&gt;",
			expected: "&lt;script&gt;alert(&#39;XSS&#39;);&lt;/script&gt;",
		},
		{
			name:     "Handles single quotes",
			input:    "Text with 'single' quotes",
			expected: "Text with &#39;single&#39; quotes",
		},
		{
			name:     "Handles double quotes",
			input:    "Text with \"double\" quotes",
			expected: "Text with &#34;double&#34; quotes",
		},
		{
			name:     "SVG based XSS vector",
			input:    "<svg><g/onload=alert(2)//<p>",
			expected: "",
		},
		{
			name:     "Style attribute with expressions",
			input:    "<div style=\"background-image: url(javascript:alert('XSS'))\">Styled div</div>",
			expected: "Styled div",
		},
	}

	for _, tc := range tests {
		t.Logf("Running test case: %s", tc.name)
		t.Run(tc.name, func(t *testing.T) {
			result := emailService.sanitizeInput(tc.input)
			if result != tc.expected {
				t.Errorf("Expected: %q\nGot: %q", tc.expected, result)
			}
		})
	}
}
