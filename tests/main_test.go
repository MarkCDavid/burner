package tests

import (
	"testing"
	pkgHello "github.com/MarkCDavid/burner/internal/hello"
)

func TestHelloWorld(t *testing.T) {
	expected := "Hello, World!"
	result := pkgHello.HelloWorld()
	
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

