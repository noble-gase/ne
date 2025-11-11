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

func TestWithMsg(t *testing.T) {
	assert.ErrorIs(t, OK.WithMsg("success"), New(0, "success"))
	assert.ErrorIs(t, Err.WithMsgF("user(id=%d) not found", 1), New(-1, "user(id=1) not found"))
}

func TestFromError(t *testing.T) {
	assert.ErrorIs(t, FromError(errors.New("something wrong")), New(-1, "something wrong"))
}
