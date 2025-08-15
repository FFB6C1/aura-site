package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

func Convert(text string) string {
	return applyBlocks(text, "converter")
}

func applyBlocks(text string, caller string) string {
	blocks := strings.Split(text, "\n\n")
	for i := range len(blocks) {
		blocks[i] = blockChecker(blocks[i], caller)
		blocks[i] = applyInline(blocks[i])
	}
	return strings.Join(blocks, "\n\n")
}

func applyInline(text string) string {
	text = applyImages(text)
	text = applyLinks(text)
	text = applyBold(text)
	text = applyItalic(text)
	text = applyCode(text)
	return text
}

// -- Helper functions to identify and apply block-level markdown

func blockChecker(text string, caller string) string {
	stripped := strings.TrimLeft(text, " ")
	if stripped == "---" {
		return blockToHorizontalLine(stripped)
	}
	if stripped[0] == '>' {
		if caller == "list" {
			return blockToQuote(text)
		}
		return blockToQuote(stripped)
	}
	if stripped[0] == '#' {
		if caller == "list" {
			return blockToHeading(text)
		}
		return blockToHeading(stripped)
	}
	if stripped[1] == '.' {
		return blockToOrderedList(stripped)
	}
	if stripped[0] == '-' || stripped[0] == '*' || stripped[0] == '+' {
		return blockToUnorderedList(stripped)
	}
	if caller == "quote" {
		return stripped
	}
	if caller == "list" {
		return blockToParagraph(text)
	}
	return blockToParagraph(stripped)
}

// -- Helper functions for block-level markdown

func blockToParagraph(text string) string {
	return fmt.Sprintf("<p>%s</p>", text)
}

func blockToQuote(text string) string {
	//CHESCA rewrite this later to add recursion, quote blocks should be capable of containing other blocks
	lines := strings.Split(text, "\n")
	for i := range len(lines) {
		lines[i] = lines[i][1:]
	}
	newText := applyBlocks(strings.Join(lines, "\n"), "quote")

	return fmt.Sprintf("<blockquote>%s</blockquote>", newText)
}

func blockToHeading(text string) string {
	count := 0
	spaces := ""
	for text[0] == ' ' {
		text = text[1:]
		spaces += " "
	}
	for text[0] == '#' {
		text = text[1:]
		count += 1
	}
	if text[0] == ' ' {
		text = text[1:]
	}
	if count > 6 {
		count = 6
	}
	return fmt.Sprintf("<h%d>%s%s</h%d>", count, spaces, text, count)
}

func blockToOrderedList(text string) string {
	//CHESCA rewrite this later to add recursion, lists should be capable of containing nested block elements
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line[0] == ' ' {
			lines[i] = applyBlocks(line, "list")
		} else if line[1] == '.' {
			if line[2] == ' ' {
				lines[i] = "<li>" + line[3:] + "</li>"
			} else {
				lines[i] = "<li>" + line[2:] + "</li>"
			}
		}
	}
	formattedList := strings.Join(lines, "\n")
	return fmt.Sprintf("<ol>\n%s\n</ol>", formattedList)
}

func blockToUnorderedList(text string) string {
	//CHESCA rewrite this later to add recursion, lists should be capable of containing nested block elements
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line[0] == ' ' {
			lines[i] = applyBlocks(line, "list")
		} else if line[1] == ' ' {
			lines[i] = "<li>" + line[2:] + "</li>"
		}
	}
	formattedList := strings.Join(lines, "\n")
	return fmt.Sprintf("<ul>\n%s\n</ul>", formattedList)
}

func blockToHorizontalLine(text string) string {
	if text == "---" {
		return "<hr>"
	} else {
		return text
	}
}

// -- Helper functions for inline markdown
func applyBold(text string) string {
	count := strings.Count(text, "**")
	if count/2 != 0 {
		count -= 1
	}
	for i := 0; i < count; i += 2 {
		text = strings.Replace(text, "**", "<strong>", 1)
		text = strings.Replace(text, "**", "</strong>", 1)
	}
	return text
}

func applyItalic(text string) string {
	count := strings.Count(text, "*")
	if count/2 != 0 {
		count -= 1
	}
	for i := 0; i < count; i += 2 {
		text = strings.Replace(text, "*", "<em>", 1)
		text = strings.Replace(text, "*", "</em>", 1)
	}
	return text
}

func applyCode(text string) string {
	count := strings.Count(text, "`")
	if count/2 != 0 {
		count -= 1
	}
	for i := 0; i < count; i += 2 {
		text = strings.Replace(text, "`", "<code>", 1)
		text = strings.Replace(text, "`", "</code>", 1)
	}
	return text
}

func applyLinks(text string) string {
	linkRegex := regexp.MustCompile("\\[[\\s\\S][^\\]]+\\]\\([\\s\\S][^\\)]+\\)")
	links := linkRegex.FindAllString(text, -1)
	if links != nil {
		for _, link := range links {
			linkParts := strings.Split(link, "](")
			newLink := fmt.Sprintf("<a href=\"%s\">%s</a>", linkParts[1][:len(linkParts[1])-1], linkParts[0][1:])
			text = strings.Replace(text, link, newLink, 1)
		}
	}
	return text
}

func applyImages(text string) string {
	imgRegex := regexp.MustCompile("!\\[[\\s\\S][^\\]]+\\]\\([\\s\\S][^\\)]+\\)")
	images := imgRegex.FindAllString(text, -1)
	if images != nil {
		for _, img := range images {
			linkParts := strings.Split(img, "](")
			newImg := fmt.Sprintf("<img src=\"%s\" alt=\"%s\" />", linkParts[1][:len(linkParts[1])-1], linkParts[0][2:])
			text = strings.Replace(text, img, newImg, 1)
		}
	}
	return text
}
