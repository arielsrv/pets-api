package test_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	//	teardown()
	os.Exit(code)
}
