package parser

import (
	"regexp"
	"strings"
)

// TransformAdmonitions converts GitHub-style admonitions in HTML to semantic admonition divs
// This should be called after markdown -> HTML conversion
func TransformAdmonitions(htmlContent string) string {
	// Match blockquotes starting with [!TYPE], capturing type and all remaining content
	pattern := regexp.MustCompile(`(?s)<blockquote>\s*<p>\[!([^\]]+)\](.*?)</blockquote>`)

	result := htmlContent
	for {
		match := pattern.FindStringSubmatchIndex(result)
		if match == nil {
			break
		}

		startIdx := match[0]
		endIdx := match[1]

		// Extract type
		typeStart := match[2]
		typeEnd := match[3]
		admonType := strings.ToLower(result[typeStart:typeEnd])

		// Extract raw content (everything after ] until </blockquote>)
		contentStart := match[4]
		contentEnd := match[5]
		rawContent := result[contentStart:contentEnd]

		// Clean up the content
		// Remove leading <br> tags and whitespace
		cleanedContent := regexp.MustCompile(`^(<br\s*/?>|\s)+`).ReplaceAllString(rawContent, "")

		// Find first </p> to properly handle paragraphs
		firstParaEnd := strings.Index(cleanedContent, "</p>")
		var finalContent string

		if firstParaEnd >= 0 {
			firstPara := cleanedContent[:firstParaEnd]
			restAfterPara := cleanedContent[firstParaEnd+4:] // skip "</p>"
			restAfterPara = strings.TrimSpace(restAfterPara)

			if restAfterPara != "" {
				finalContent = "<p>" + strings.TrimSpace(firstPara) + "</p>" + restAfterPara
			} else {
				finalContent = "<p>" + strings.TrimSpace(firstPara) + "</p>"
			}
		} else {
			finalContent = "<p>" + strings.TrimSpace(cleanedContent) + "</p>"
		}

		admonition := buildAdmonitionHTML(admonType, finalContent)
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
