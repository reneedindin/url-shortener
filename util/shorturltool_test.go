package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenUrlID(t *testing.T) {
	url := "https://www.google.com"
	expireAt := "2021-07-18T16:58:30+08:00"
	urlID := GenUrlID(url, expireAt)
	expected := "aHR0cHM6Ly93d3cuZ29vZ2xlLmNvbTIwMjEtMDctMThUMTY6NTg6MzArMDg6MDA="
	assert.Equal(t, expected, urlID)
}
