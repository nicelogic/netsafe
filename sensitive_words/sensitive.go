package sensitive_words

import (
	"context"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mozillazg/go-pinyin"
	"github.com/nicelogic/netsafe/sensitive_words/dfa"
)

type SensitiveWordChecker interface {
	// Hit 判断是否命中敏感词，且返回命中的敏感词
	Hit(ctx context.Context, text string) (isHit bool, hitWord string, err error)
	// HitMust 严格模式，最少命中几个敏感词
	HitMust(ctx context.Context, text string, times int) (isHit bool, hitWords []string, err error)
	// MatchReplace 敏感词替换
	MatchReplace(ctx context.Context, text string) (isHit bool, lastText string, err error)
	// DebugInfos 输出当前所有敏感词
	DebugInfos(ctx context.Context) (results []*dfa.Stats)
}

var _ SensitiveWordChecker = (*sensitiveWord)(nil)

func New(buildWords BuildWordsFn, path string, opts ...Option) (SensitiveWordChecker, error) {
	o := options{
		maskWord:       '*',
		buildWordsCall: buildWords,
		mode:           ModePinyin,
	}
	for _, fn := range opts {
		fn(&o)
	}

	st := &sensitiveWord{
		options: o,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_ = cancel
	if err := st.buildWords(ctx, path); err != nil {
		log.Printf("build words failed.err(%v)\n", err)
		return nil, err
	}

	log.Printf("init success\n")
	return st, nil
}

type sensitiveWord struct {
	options
	trieTree atomic.Value
}

func (st *sensitiveWord) buildWords(ctx context.Context, path string) error {
	log.Printf("rebuild words, start_time(%v)\n", time.Now().Format("2006-01-02 15:04:05"))
	words, err := st.buildWordsCall(ctx, path)
	if err != nil {
		return err
	}

	tree := dfa.NewTrieTree()
	tree.WithFilterChars(st.filterChars)

	if err = st.mode.Range(func(value Mode) error {
		switch value {
		case ModePinyin: // 开启拼音模式
			for _, word := range words {
				if !pinyinWordReg.MatchString(word) {
					continue
				}
				var pinyinWords []string
				for _, segWord := range strings.Split(word, "|") {
					pinyinWords = append(pinyinWords, strings.Join(pinyin.LazyConvert(segWord, nil), ""))
				}
				words = append(words, strings.Join(pinyinWords, "|"))
			}
		case ModeStats: // 开启命中敏感词统计
			tree.WithStats()
		}
		return nil
	}); err != nil {
		return err
	}

	tree.AddWords(words...)
	st.trieTree.Store(tree)
	log.Printf("rebuild words success, end_time(%v)\n", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func (st *sensitiveWord) Hit(ctx context.Context, text string) (isHit bool, hitWord string, err error) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	isHit, hitWords := tree.Detect(text, 1)
	if isHit {
		return true, hitWords[0], nil
	}
	return false, "", nil
}

func (st *sensitiveWord) HitMust(ctx context.Context, text string, times int) (isHit bool, hitWords []string, err error) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	isHit, hitWords = tree.Detect(text, times)
	return isHit, hitWords, nil
}

func (st *sensitiveWord) MatchReplace(ctx context.Context, text string) (isHit bool, lastText string, err error) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	isHit, lastText = tree.Replace(text, st.maskWord)
	return isHit, lastText, nil
}

func (st *sensitiveWord) DebugInfos(ctx context.Context) (results []*dfa.Stats) {
	tree := st.trieTree.Load().(*dfa.TrieTree)
	return tree.DebugInfos()
}
