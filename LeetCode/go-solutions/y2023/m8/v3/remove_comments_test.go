package main

import (
	"fmt"
	"testing"
)

func TestRemoveComments(t *testing.T) {

	t.Run("1", func(t *testing.T) {

		resource := []string{
			"/*Test program */", "int main()",
			"{ ",
			"  // variable declaration ",
			"int a, b, c;",
			"/* This is a test",
			"   multiline  ",
			"   comment for ",
			"   testing */",
			"a = b + c;",
			"}",
		}

		got := removeComments(resource)
		fmt.Printf("got: %v", got)
	})

	t.Run("2", func(t *testing.T) {

		resource := []string{
			"a/*comment",
			"line",
			"more_comment*/b",
		}

		got := removeComments(resource)
		fmt.Printf("got: %v", got)
	})

	t.Run("3", func(t *testing.T) {

		resource := []string{
			"a",
		}

		got := removeComments(resource)
		fmt.Printf("got: %v", got)
	})

	t.Run("4", func(t *testing.T) {

		resource := []string{
			"main() {",
			"   int x = 1; // Its comments here",
			"   x++;",
			"   cout << x;",
			"   //return x;",
			"   x--;",
			"}",
		}

		got := removeComments(resource)
		fmt.Printf("got: %v", got)
	})
}
