package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv_Get(t *testing.T) {
	assert.NotEqual(t, "blubb", Get("PATH", "blubb"))
	assert.Equal(t, "blabb", Get("blibb", "blabb"))
}
