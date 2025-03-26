package ne

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekAround(t *testing.T) {
	tz, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)
	now := time.Unix(1562909685, 0).In(tz)
	monday, sunday := WeekAround(time.DateOnly, now)
	assert.Equal(t, "2019-07-08", monday)
	assert.Equal(t, "2019-07-14", sunday)
}
