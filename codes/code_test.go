package codes

import (
	"errors"
	"fmt"
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
	assert.True(t, Is(fmt.Errorf("oh yeah: %w", New(0, "success")), OK))
	assert.False(t, Is(New(1, "failed"), OK))
}

func TestNew(t *testing.T) {
	t.Log(OK.New("success"))
	t.Log(Unknown.New("user(id=%d) not found", 1))
}
