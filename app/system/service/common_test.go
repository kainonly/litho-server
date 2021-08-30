package service

import (
	"os"
	"testing"
)

var s *Tests

type Tests struct {
	Index *Index
}

func TestMain(m *testing.M) {
	os.Chdir(`../../../`)
}
