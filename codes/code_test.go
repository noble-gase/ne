package codes

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	err := errors.New("oh no")
	assert.True(t, Is(nil, nil))
	assert.False(t, Is(nil, OK))
	assert.False(t, Is(err, nil))
	assert.False(t, Is(err, OK))
	assert.True(t, Is(New(0, "success"), OK))
	assert.False(t, Is(New(1, "failed"), OK))
}
