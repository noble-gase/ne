package validates

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func NullStringRequired(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) != 0
}

func NullIntGte(fl validator.FieldLevel) bool {
	i, err := strconv.ParseInt(fl.Param(), 0, 64)
	if err != nil {
		return false
	}
	return fl.Field().Int() >= i
}

type ParamsValidate struct {
	ID   sql.NullInt64  `valid:"nullint_gte=10"`
	Desc sql.NullString `valid:"nullstring_required"`
}

func TestDefaultValidator(t *testing.T) {
	params1 := new(ParamsValidate)
	params1.ID = sql.NullInt64{
		Int64: 9,
		Valid: true,
	}

	opts := []Option{
		WithValuerType(sql.NullString{}, sql.NullInt64{}),
		WithValidateFunc("nullint_gte", NullIntGte),
		WithTranslation("nullint_gte", "{0}必须大于或等于{1}", true),
		WithValidateFunc("nullstring_required", NullStringRequired),
		WithTranslation("nullstring_required", "{0}为必填字段", true),
	}

	err := ValidateStruct(params1, opts...)
	assert.NotNil(t, err)
	t.Log("err validate params:", err.Error())

	params2 := &ParamsValidate{
		ID: sql.NullInt64{
			Int64: 13,
			Valid: true,
		},
		Desc: sql.NullString{
			String: "og",
			Valid:  true,
		},
	}
	err = ValidateStruct(params2, opts...)
	assert.Nil(t, err)
}

func TestNewValidator(t *testing.T) {
	testV := New(
		WithValuerType(sql.NullString{}, sql.NullInt64{}),
		WithValidateFunc("nullint_gte", NullIntGte),
		WithTranslation("nullint_gte", "{0}必须大于或等于{1}", true),
		WithValidateFunc("nullstring_required", NullStringRequired),
		WithTranslation("nullstring_required", "{0}为必填字段", true),
	)

	params1 := new(ParamsValidate)
	params1.ID = sql.NullInt64{
		Int64: 9,
		Valid: true,
	}
	err := testV.ValidateStruct(params1)
	assert.NotNil(t, err)
	t.Log("err validate params:", err.Error())

	params2 := &ParamsValidate{
		ID: sql.NullInt64{
			Int64: 13,
			Valid: true,
		},
		Desc: sql.NullString{
			String: "og",
			Valid:  true,
		},
	}
	err = testV.ValidateStruct(params2)
	assert.Nil(t, err)
}
