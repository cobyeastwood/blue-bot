package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilityStringUpper(t *testing.T) {
	t.Run("Takes a string and returns a string uppercase", func(t *testing.T) {
		var s string

		s = "tEst"

		assert.Equal(t, "TEST", StringUpper(s))
	})
}
