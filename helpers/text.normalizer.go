package helpers

import (
	"regexp"
	"strings"
)

var (
	// 删除标点符号前的空格
	re1, re2, re3 *regexp.Regexp
)

// TextNormalizer 文本处理
type TextNormalizer struct {
	text            string // 待处理的文本
	trimSpace       bool   // 删除文本两端多余的空格
	cleanExtraSpace bool   // 清理多余空格（多个单词之间、标点符号前后多余的空格）
	halfWidth       bool   // 转为半角
}

func init() {
	re1 = regexp.MustCompile(`\s+`)                // 删除多余空格
	re2 = regexp.MustCompile(`\s([,.!?;:()\[\]])`) // 删除标点符号前的空格
	re3 = regexp.MustCompile(`([,.!?;:])(\S)`)     // 补充标点符号后的空格，但不包括引号和括号
}

func NewTextNormalizer() *TextNormalizer {
	return &TextNormalizer{}
}

func (n *TextNormalizer) SetText(text string) *TextNormalizer {
	n.text = text
	return n
}

func (n *TextNormalizer) TrimSpace(v bool) *TextNormalizer {
	n.trimSpace = v
	return n
}

func (n *TextNormalizer) CleanExtraSpace(v bool) *TextNormalizer {
	n.cleanExtraSpace = v
	return n
}

func (n *TextNormalizer) HalfWidth(v bool) *TextNormalizer {
	n.halfWidth = v
	return n
}

func (n *TextNormalizer) String() string {
	text := n.text
	if text == "" {
		return ""
	}

	if n.trimSpace {
		text = strings.TrimSpace(text)
	}
	if n.halfWidth {
		text = strings.NewReplacer(
			"，", ",",
			"（", "(",
			"）", ")",
			"！", "!",
			"？", "?",
			"：", ":",
			"；", ";",
			"【", "[",
			"】", "]",
			"。", ".",
			"“", "\"",
			"”", "\"",
		).Replace(text)
	}
	if n.cleanExtraSpace {
		return re3.ReplaceAllString(re2.ReplaceAllString(re1.ReplaceAllString(text, " "), "$1"), "$1 $2")
	}
	return text
}
