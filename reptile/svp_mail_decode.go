package reptile

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

/*
这个 JavaScript 脚本是 Cloudflare 用来保护网页上的电子邮件地址不被爬虫轻易抓取的功能，被称为 "Email Address Obfuscation"。
它的核心原理是：
将原始电子邮件地址进行 XOR 加密，并编码为十六进制字符串。
在前端页面上，用一个特殊的 <a> 标签（带有 href="/cdn-cgi/l/email-protection#..."）或者一个带有 class="__cf_email__" 和 data-cfemail="..." 属性的 <span> 标签来占位。
当浏览器加载页面时，此 JavaScript 脚本会执行，查找这些特殊标签。
脚本读取十六进制字符串，使用第一个字节作为密钥，对后续所有字节进行 XOR 解密，还原出原始的电子邮件地址。
最后，用解密后的邮件地址替换掉页面上的占位标签。
*/
const (
	// Cloudflare's email protection path prefix
	cfEmailProtectionPath = "/cdn-cgi/l/email-protection#"
	// The class name for email-protected elements
	cfEmailClassName = "__cf_email__"
	// The data attribute holding the encoded email
	cfDataAttribute = "data-cfemail"
)

// decodeCfEmail decodes a Cloudflare-encoded hexadecimal string into an email address.
// This is the core logic ported from the JS `n` function.
func decodeCfEmail(encodedString string) (string, error) {
	// Must have at least 2 chars for the key and some for the content
	if len(encodedString) < 2 {
		return "", fmt.Errorf("encoded string is too short: %s", encodedString)
	}

	// The first two hex characters are the XOR key
	keyHex := encodedString[0:2]
	key, err := strconv.ParseInt(keyHex, 16, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse key '%s': %w", keyHex, err)
	}

	var decodedBytes []byte
	// Loop through the rest of the string, two characters at a time
	for i := 2; i < len(encodedString); i += 2 {
		// Get the next two hex characters
		charHex := encodedString[i : i+2]
		// Parse them as a hex number
		charCode, err := strconv.ParseInt(charHex, 16, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse char '%s': %w", charHex, err)
		}
		// XOR with the key to get the original character code
		decodedBytes = append(decodedBytes, byte(charCode^key))
	}

	return string(decodedBytes), nil
}

// processNode recursively traverses the HTML tree, finds protected emails, and decodes them.
// This function replaces the JS functions `c`, `o`, `a`, and `i`.
func processNode(n *html.Node) {
	if n.Type == html.ElementNode {
		// Case 1: Handle <a> tags with protected href
		// Corresponds to JS function `c`
		if n.Data == "a" {
			for i, attr := range n.Attr {
				if attr.Key == "href" && strings.HasPrefix(attr.Val, cfEmailProtectionPath) {
					encodedString := strings.TrimPrefix(attr.Val, cfEmailProtectionPath)
					if decodedEmail, err := decodeCfEmail(encodedString); err == nil {
						// Replace the attribute value
						n.Attr[i].Val = "mailto:" + decodedEmail
					}
				}
			}
		}

		// Case 2: Handle elements with data-cfemail attribute
		// Corresponds to JS function `o`
		isCfEmailElement := false
		var encodedString string
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, cfEmailClassName) {
				isCfEmailElement = true
			}
			if attr.Key == cfDataAttribute {
				encodedString = attr.Val
			}
		}

		if isCfEmailElement && encodedString != "" {
			if decodedEmail, err := decodeCfEmail(encodedString); err == nil {
				// Create a new text node with the decoded email
				textNode := &html.Node{
					Type: html.TextNode,
					Data: decodedEmail,
				}
				// Replace the protected element with the new text node
				if n.Parent != nil {
					n.Parent.InsertBefore(textNode, n)
					n.Parent.RemoveChild(n)
				}
			}
		}
	}

	// Recurse for all children
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processNode(c)
	}
}

// DecodeCloudflareEmails takes an HTML string as input and returns a version
// with all Cloudflare-protected emails decoded.
func DecodeCloudflareEmails(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	processNode(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to render HTML: %w", err)
	}

	// html.Render adds <html><head></head><body>...</body></html> wrappers,
	// which we might not want if the input was a fragment.
	// This is a simple way to strip them for cleaner output.
	output := buf.String()
	output = strings.TrimPrefix(output, "<html><head></head><body>")
	output = strings.TrimSuffix(output, "</body></html>")

	return output, nil
}

// processSvpLinks finds all Cloudflare-protected links in a block of text
// and replaces them with the decoded content.
// This version is much simpler as it uses regex for this specific text format.
func processSvpLinks(text string) (string, error) {
	// This regex finds the entire <a> tag but specifically captures
	// the content of the `data-cfemail` attribute.
	// <a ... data-cfemail="<CAPTURE THIS>" ...> ... </a>
	re := regexp.MustCompile(`<a [^>]*data-cfemail="([a-f0-9]+)"[^>]*>.*?</a>`)

	// ReplaceAllStringFunc finds all matches and calls a function to get the replacement string.
	// This is very efficient for this kind of task.
	processedText := re.ReplaceAllStringFunc(text, func(match string) string {
		// `FindStringSubmatch` returns the full match and all captured groups.
		// submatches[0] is the full match (the whole <a> tag)
		// submatches[1] is the first captured group (the hex string)
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			// This should not happen if the regex matched, but it's a safe check.
			return match
		}

		encodedString := submatches[1]
		decoded, err := decodeCfEmail(encodedString)
		if err != nil {
			log.Printf("Warning: could not decode '%s', leaving original. Error: %v", encodedString, err)
			return match // On error, leave the original tag
		}

		// Return the decoded string as the replacement for the entire <a> tag.
		return decoded
	})

	return processedText, nil
}

// ParseSvpHtml 解析HTML并解析SVP邮件后获取URL
func ParseSvpHtml(htmlContent []byte) string {
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		panic(err.Error())
	}
	// 找到最后一个pre
	pre := doc.Find(`pre`).Last()
	if pre.Length() > 0 {
		ret, err := pre.Html()
		if err != nil {
			log.Printf("An error occurred: %v\n", err)
		}
		decodedHTML, err := processSvpLinks(ret)
		if err != nil {
			log.Fatalf("An error occurred: %v", err)
		} else {
			// 解析HTML
			doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(decodedHTML)))
			if err != nil {
				panic(err.Error())
			}
			return doc.Text()
		}
	}
	return ""
}
