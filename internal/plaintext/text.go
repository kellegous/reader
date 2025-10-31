package plaintext

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func From(content string) string {
	if !strings.Contains(content, "<html") {
		content = fmt.Sprintf("<html><body>%s</body></html>", content)
	}

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content
	}

	var b strings.Builder
	toTextContent(doc, &b)
	return b.String()
}

func toTextContent(root *html.Node, b *strings.Builder) {
	if root == nil {
		return
	}

	switch root.Type {
	case html.TextNode:
		b.WriteString(root.Data)
	case html.ElementNode:
		// Skip script and style elements
		if root.Data == "script" || root.Data == "style" {
			return
		}

		// Add space before certain block elements
		if root.Data == "p" || root.Data == "div" || root.Data == "br" {
			b.WriteString(" ")
		}

		// Recursively process child nodes
		for child := root.FirstChild; child != nil; child = child.NextSibling {
			toTextContent(child, b)
		}

		// Add space after certain block elements
		if root.Data == "p" || root.Data == "div" {
			b.WriteString(" ")
		}
	case html.DocumentNode:
		for child := root.FirstChild; child != nil; child = child.NextSibling {
			toTextContent(child, b)
		}
	}
}
