package test

import (
	"testing"

	"daonao.com/quote/lib"
)

func TestHello(t *testing.T) {
	want := "你好，世界。"
	if got := lib.Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
