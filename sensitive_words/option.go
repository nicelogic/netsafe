package sensitive_words

import (
)

type options struct {
	// 掩码字符，默认使用 *
	maskWord rune
	// 查找/替换模式，默认开启拼音
	mode Mode
	// 过滤特殊字符，默认过滤除中英文数字之外的所有字符
	filterChars []rune
	// 创建敏感词回调方法
	buildWordsCall BuildWordsFn
}

type Option func(*options)

func WithMaskWord(word rune) Option {
	return func(o *options) {
		o.maskWord = word
	}
}

func WithMode(modes ...Mode) Option {
	return func(o *options) {
		for _, m := range modes {
			o.mode |= m
		}
	}
}

func WithFilterChars(filterChars ...rune) Option {
	return func(o *options) {
		o.filterChars = filterChars
	}
}

