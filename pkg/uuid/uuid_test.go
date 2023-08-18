package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	id := Random()
	assert.NotEmpty(t, id)
}
