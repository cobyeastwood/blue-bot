package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchLong(t *testing.T) {
	t.Run("Returns a set of data for a stock over a particular time period", func(t *testing.T) {

		var d DataSet

		fc := NewFetchConfig(MINUTE, "10", 200)

		assert.IsType(t, d, FetchLong("Hepa", fc))
	})
}
func TestFetchShort(t *testing.T) {
	t.Run("Returns a set of data for a stock over a particular time period", func(t *testing.T) {

		var testStr [][]string

		mc := NewMainConfig()

		assert.IsType(t, testStr, FetchShort("Hepa", mc.c))
	})
}
