package netsafe

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	sw "github.com/nicelogic/netsafe/sensitive_words"
)

type NetSafe struct {
	sensitiveWordChecker sw.SensitiveWordChecker
}

func (client *NetSafe) Init(ctx context.Context) error {
	var err error
	client.sensitiveWordChecker, err = sw.New(
		buildWordsCall,
		sw.WithMode(sw.ModePinyin, sw.ModeStats),
		sw.WithMaskWord('*'),
	)
	if err != nil {
		return err
	}
	return nil
}

func (client *NetSafe) CheckSensitiveWords(ctx context.Context, text string) (bool, string, error) {
	return client.sensitiveWordChecker.Hit(ctx, text)
}

func buildWordsCall(ctx context.Context) (words []string, err error) {
	lines, err := readLines("./sensitive_words.txt")
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func readLines(path string) ([]string, error) {
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
