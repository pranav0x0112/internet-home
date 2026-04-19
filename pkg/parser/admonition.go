package parser

import (
	"regexp"
	"strings"
)

// TransformAdmonitions converts GitHub-style admonitions in HTML to semantic admonition divs
// This should be called after markdown -> HTML conversion
func TransformAdmonitions(htmlContent string) string {
	// Pattern to match blockquotes that start with [!TYPE]
	// Uses (?s) flag for DOTALL mode to match across newlines
	pattern := regexp.MustCompile(`(?s)<blockquote>\s*<p>\[!([^\]]+)\]\s*([^<]*)</p>(.*?)</blockquote>`)

	result := htmlContent
	matches := pattern.FindAllStringSubmatchIndex(result, -1)

	// Process matches in reverse order to maintain string indices when replacing
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		startIdx := match[0]
		endIdx := match[1]

		// Extract components from groups
		typeStart := match[2]
		typeEnd := match[3]
		admonType := strings.ToLower(htmlContent[typeStart:typeEnd])

		// Group 2: first line content (after [!TYPE] and before </p>)
		firstLineStart := match[4]
		firstLineEnd := match[5]
		firstLineContent := strings.TrimSpace(htmlContent[firstLineStart:firstLineEnd])

		// Group 3: remaining content
		restStart := match[6]
		restEnd := match[7]
		restContent := htmlContent[restStart:restEnd]

		// Build complete content
		var contentBuilder strings.Builder
		if firstLineContent != "" {
			contentBuilder.WriteString("<p>")
			contentBuilder.WriteString(firstLineContent)
			contentBuilder.WriteString("</p>")
		}
		contentBuilder.WriteString(strings.TrimSpace(restContent))

		content := contentBuilder.String()

		// Build the semantic admonition HTML
		admonition := buildAdmonitionHTML(admonType, content)

		// Replace in result
		result = result[:startIdx] + admonition + result[endIdx:]
	}

	return result
}

// buildAdmonitionHTML creates semantic HTML for an admonition
func buildAdmonitionHTML(admonType string, content string) string {
	normalizedType := normalizeType(admonType)

	// Construct the admonition div with proper structure
	html := `<div class="admonition ` + admonType + `">
  <div class="admonition-title">` + normalizedType + `</div>
  <div class="admonition-content">
    ` + content + `
  </div>
</div>`

	return html
}

// normalizeType converts admonition type to title case
// "note" -> "Note", "fun fact" -> "Fun Fact"
func normalizeType(typeStr string) string {
	parts := strings.Fields(typeStr)
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, " ")
}
