package sqls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinify(t *testing.T) {
	s := `SELECT
	id,
	name
FROM
	demo
WHERE
	name LIKE '%hello%';`

	assert.Equal(t, "SELECT id, name FROM demo WHERE name LIKE '%hello%';", Minify(s))
}
