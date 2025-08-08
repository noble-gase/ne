package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	now1 := time.Now()
	err1 := Retry(context.Background(), func(ctx context.Context) error {
		fmt.Println("Retry...")
		return nil
	}, 3, time.Second)
	assert.Nil(t, err1)
	assert.Equal(t, 0, int(time.Since(now1).Seconds()))

	now2 := time.Now()
	err2 := Retry(context.Background(), func(ctx context.Context) error {
		fmt.Println("Retry...")
		return errors.New("something wrong")
	}, 3, time.Second)
	assert.NotNil(t, err2)
	assert.Equal(t, 2, int(time.Since(now2).Seconds()))
}
