package ne

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalNoEscapeHTML(t *testing.T) {
	data := map[string]string{"url": "https://github.com/noble-gase?id=996&name=og"}

	b, err := MarshalNoEscapeHTML(data)
	assert.Nil(t, err)
	assert.Equal(t, string(b), `{"url":"https://github.com/noble-gase?id=996&name=og"}`)
}

func TestVersionCompare(t *testing.T) {
	ok, err := VersionCompare("1.0.0", "1.0.0")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("1.0.0", "1.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare("=1.0.0", "1.0.0")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("=1.0.0", "1.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare("!=4.0.4", "4.0.0")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("!=4.0.4", "4.0.4")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare(">2.0.0", "2.0.1")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare(">2.0.0", "1.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare(">=1.0.0&<2.0.0", "1.0.2")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare(">=1.0.0&<2.0.0", "2.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)

	ok, err = VersionCompare("<2.0.0|>3.0.0", "1.0.2")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("<2.0.0|>3.0.0", "3.0.1")
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = VersionCompare("<2.0.0|>3.0.0", "2.0.1")
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestIsUniqueDuplicateError(t *testing.T) {
	errMySQL := errors.New("Duplicate entry 'value' for key 'key_name'")
	assert.True(t, IsUniqueDuplicateError(errMySQL))

	errPgSQL := errors.New(`duplicate key value violates unique constraint "constraint_name"`)
	assert.True(t, IsUniqueDuplicateError(errPgSQL))

	errSQLite := errors.New("UNIQUE constraint failed: table_name.column_name")
	assert.True(t, IsUniqueDuplicateError(errSQLite))
}

func TestExcelColumnIndex(t *testing.T) {
	assert.Equal(t, 0, ExcelColumnIndex("A"))
	assert.Equal(t, 1, ExcelColumnIndex("B"))
	assert.Equal(t, 26, ExcelColumnIndex("AA"))
	assert.Equal(t, 27, ExcelColumnIndex("AB"))
}
