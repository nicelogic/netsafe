package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	sw "github.com/nicelogic/netsafe/sensitive_words"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	st, _ := sw.New(
		buildWordsCall,
		sw.WithMode(sw.ModePinyin, sw.ModeStats),
		sw.WithMaskWord('*'),
	)
	ctx := context.Background()

	// 判断敏感词是否命中
	isHit, hitWord, err := st.Hit(ctx, "xijinping")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, hit_word: 丑八怪
	fmt.Printf("is_hit: %t, hit_word: %s\n", isHit, hitWord)

	// 敏感词替换
	isHit, newText, err := st.MatchReplace(ctx, "你这个丑逼")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, new_text: 你这个**
	fmt.Printf("is_hit: %t, new_text: %s\n", isHit, newText)

	// 组合词匹配
	isHit, hitWord, err = st.Hit(ctx, "听说司马南在美国买房子")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, hit_word: 司马南|美国
	fmt.Printf("is_hit: %t, hit_word: %s\n", isHit, hitWord)

	// 组合词替换
	isHit, newText, err = st.MatchReplace(ctx, "听说司马南在美国买房子")
	if err != nil {
		panic(err)
	}
	// 输出：is_hit: true, new_text: 听说***在**买房子
	fmt.Printf("is_hit: %t, new_text: %s\n", isHit, newText)

	// debug info
	// infos := st.DebugInfos(ctx)
	// for _, info := range infos {
	// 	// 输出：word: 丑八怪, hit_count: 1
	// 	// 输出：word: 丑逼, hit_count: 1
	// 	// 输出：word: 司马南|美国, hit_count: 2
	// 	// 输出：word: 方舟子|死了, hit_count: 0
	// 	// 输出：word: choubaguai, hit_count: 0
	// 	// 输出：word: choubi, hit_count: 0
	// 	// 输出：word: simanan|meiguo, hit_count: 0
	// 	// 输出：word: fangzhouzi|sile, hit_count: 0
	// 	fmt.Printf("word: %s, hit_count: %d\n", info.Word, info.HitCount)
	// }

	time.Sleep(time.Minute)
}

func buildWordsCall(ctx context.Context) (words []string, err error) {

	lines, err := ReadLines("../../sensitive_word_dict.txt")
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func ReadLines(path string) ([]string, error) {
	log.Printf("read lines begin(%v)\n", time.Now())
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	log.Printf("read lines end(%v)\n", time.Now())
	return lines, scanner.Err()
}
