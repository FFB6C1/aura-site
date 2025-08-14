package markdown

import (
	"fmt"
	"testing"
)

func TestBold(t *testing.T) {
	text := "hello **world**"
	expected := "hello <strong>world</strong>"

	testText := applyBold(text)

	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBoldWithBadFormatting(t *testing.T) {
	text := "hello **world****"
	expected := "hello <strong>world</strong>**"

	testText := applyBold(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestItalic(t *testing.T) {
	text := "hello *world*"
	expected := "hello <em>world</em>"

	testText := applyItalic(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBoldAndItalic(t *testing.T) {
	text := "*hello* **world**"
	expected := "<em>hello</em> <strong>world</strong>"

	testText := applyInline(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestItalicWithinBold(t *testing.T) {
	text := "***hello world***"
	expected := "<strong><em>hello world</strong></em>"

	testText := applyInline(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestLinks(t *testing.T) {
	text := "[hello](world)"
	expected := "<a href=\"world\">hello</a>"

	testText := applyLinks(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestLinksMultiple(t *testing.T) {
	text := "[hello](world) this is a [web](site)"
	expected := "<a href=\"world\">hello</a> this is a <a href=\"site\">web</a>"

	testText := applyLinks(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestImages(t *testing.T) {
	text := "![hello](world)"
	expected := "<img src=\"world\" alt=\"hello\" />"

	testText := applyImages(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestImagesMultiple(t *testing.T) {
	text := "![hello](world) this is a ![web](site)"
	expected := "<img src=\"world\" alt=\"hello\" /> this is a <img src=\"site\" alt=\"web\" />"

	testText := applyImages(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToParagraph(t *testing.T) {
	text := "hello world"
	expected := "<p>hello world</p>"

	testText := blockToParagraph(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToQuote(t *testing.T) {
	text := "> This is a quote\n>\n> spanning three lines"
	expected := "<blockquote>This is a quote\n\nspanning three lines</blockquote>"

	testText := blockToQuote(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToNestedQuote(t *testing.T) {
	text := "> This is a quote\n>\n>> containing a quote\n>\n> ###and a heading\n>\n> spanning four lines"
	expected := "<blockquote>This is a quote\n\n<blockquote>containing a quote</blockquote>\n\n<h3>and a heading</h3>\n\nspanning four lines</blockquote>"

	testText := blockToQuote(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToHeading(t *testing.T) {
	text := "##hello #world"
	expected := "<h2>hello #world</h2>"

	testText := blockToHeading(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToHeadingWithTooManyPound(t *testing.T) {
	text := "#############hello #world"
	expected := "<h6>hello #world</h6>"

	testText := blockToHeading(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToOrderedListNoRecursion(t *testing.T) {
	text := "1. hello\n2. world\n3.nospace"
	expected := "<ol>\n<li>hello</li>\n<li>world</li>\n<li>nospace</li>\n</ol>"

	testText := blockToOrderedList(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToOrderedListRecursion(t *testing.T) {
	text := "1. hello\n    2. world\n3.nospace"
	expected := "<ol>\n<li>hello</li>\n<ol>\n<li>world</li>\n</ol>\n<li>nospace</li>\n</ol>"

	testText := blockToOrderedList(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToUnorderedListNoRecursion(t *testing.T) {
	text := "- hello\n- world"
	expected := "<ul>\n<li>hello</li>\n<li>world</li>\n</ul>"

	testText := blockToUnorderedList(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockToUnorderedListRecursion(t *testing.T) {
	text := "- hello\n    - cruel\n    horrible\n- world"
	expected := "<ul>\n<li>hello</li>\n<ul>\n<li>cruel</li>\n</ul>\n<p>    horrible</p>\n<li>world</li>\n</ul>"

	testText := blockToUnorderedList(text)
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}

func TestBlockChecker(t *testing.T) {
	texts := []string{
		"paragraph",
		"1. ordered\n2. list",
		"- unordered\n- list",
		"---",
		"> quote",
		"## heading2",
	}
	expected := []string{
		"<p>paragraph</p>",
		"<ol>\n<li>ordered</li>\n<li>list</li>\n</ol>",
		"<ul>\n<li>unordered</li>\n<li>list</li>\n</ul>",
		"<hr>",
		"<blockquote>quote</blockquote>",
		"<h2>heading2</h2>",
	}

	for i := range len(texts) {
		testText := blockChecker(texts[i], "test")

		if testText != expected[i] {
			t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected[i], testText))
		}
	}
}

func TestApplyBlocks(t *testing.T) {
	text := "#hello\n\nworld\n\n1. hello\n2. world\n\n---\n\n>hello world"
	expected := "<h1>hello</h1>\n\n<p>world</p>\n\n<ol>\n<li>hello</li>\n<li>world</li>\n</ol>\n\n<hr>\n\n<blockquote>hello world</blockquote>"

	testText := applyBlocks(text, "test")
	if testText != expected {
		t.Fatal(fmt.Sprintf("Expected: %s, Result: %s", expected, testText))
	}
}
