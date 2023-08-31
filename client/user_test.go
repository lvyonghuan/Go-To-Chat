package client

import (
	"log"
	"testing"
)

func TestUserRegister(t *testing.T) {
	err, uid := registerService("test", "114")
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
	log.Println(uid)
}
