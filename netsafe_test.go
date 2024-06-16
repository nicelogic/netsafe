package netsafe

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestUserFromJwt(t *testing.T) {
	netSafe := NetSafe{}
	ctx := context.Background()
	err := (&netSafe).Init(ctx, "./sensitive_words.txt")
	if err != nil {
		t.Errorf("init err")
		return;
	}
	isHit, hitWord, err := netSafe.CheckSensitiveWords(ctx, "caoniaye")
	if err != nil {
		t.Errorf("init err")
		return;
	}
	log.Printf("isHit(%v).hitWord(%v)\n", isHit, hitWord)
	time.Sleep(time.Duration(time.Second * 10))
}
