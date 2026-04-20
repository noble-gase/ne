package protokit

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValuesToMessage(t *testing.T) {
	query := "user[name]=foo&user.tags=GO&user.tags=RUST&user.attrs.age=20&user.attrs.height=180&friends[0].name=bar&friends[0].tags=PHP&friends[0].tags=JAVA&friends[0].attrs.age=18&friends[0].attrs.height=175"
	values, err := url.ParseQuery(query)
	assert.Nil(t, err)

	msg := new(DemoRequest)
	err = ValuesToMessage(msg, values)
	assert.Nil(t, err)
	b, _ := json.Marshal(msg)
	t.Log(string(b))

	// --- User ---
	assert.Equal(t, "foo", msg.User.Name)
	assert.Equal(t, []Tag{Tag_GO, Tag_RUST}, msg.User.Tags)
	assert.Equal(t, "20", msg.User.Attrs["age"])
	assert.Equal(t, "180", msg.User.Attrs["height"])

	// --- Friends ---
	assert.Len(t, msg.Friends, 1)

	assert.Equal(t, "bar", msg.Friends[0].Name)
	assert.Equal(t, []Tag{Tag_PHP, Tag_JAVA}, msg.Friends[0].Tags)
	assert.Equal(t, "18", msg.Friends[0].Attrs["age"])
	assert.Equal(t, "175", msg.Friends[0].Attrs["height"])
}

func TestMessageToValues(t *testing.T) {
	msg1 := &DemoRequest{
		User: &User{
			Name: "foo",
			Tags: []Tag{Tag_GO, Tag_RUST},
			Attrs: map[string]string{
				"age":    "20",
				"height": "180",
			},
		},
		Friends: []*User{
			{
				Name: "bar",
				Tags: []Tag{Tag_PHP, Tag_JAVA},
				Attrs: map[string]string{
					"age":    "18",
					"height": "175",
				},
			},
		},
	}
	values := MessageToValues(msg1)
	query, err := url.QueryUnescape(values.Encode())
	assert.Nil(t, err)
	t.Log("[query]", query)

	// 验证values
	msg2 := new(DemoRequest)
	err = ValuesToMessage(msg2, values)
	assert.Nil(t, err)
	b, _ := json.Marshal(msg2)
	t.Log("[message]", string(b))
}

func TestEnumNumber(t *testing.T) {
	values, _ := url.ParseQuery("user.tags=0&user.tags=1&user.tags=2")
	msg := new(DemoRequest)
	assert.NoError(t, ValuesToMessage(msg, values))
	assert.Equal(t, []Tag{Tag_GO, Tag_RUST, Tag_PHP}, msg.User.Tags)
}

func TestScalarListWithIndex(t *testing.T) {
	values, _ := url.ParseQuery("user.tags[0]=GO&user.tags[2]=RUST&user.tags[-1]=PHP")
	msg := new(DemoRequest)
	assert.NoError(t, ValuesToMessage(msg, values))
	// 下标只决定顺序；-1, 0, 2 升序后为 PHP, GO, RUST
	assert.Equal(t, []Tag{Tag_PHP, Tag_GO, Tag_RUST}, msg.User.Tags)
}

func TestPanicResistance(t *testing.T) {
	values, _ := url.ParseQuery("a=1&a.b=2") // 类型冲突
	msg := new(DemoRequest)
	assert.NotPanics(t, func() { _ = ValuesToMessage(msg, values) })
}

func TestOOMResistance(t *testing.T) {
	values, _ := url.ParseQuery("friends[2147483647].name=x")
	msg := new(DemoRequest)
	// 不应 OOM（应在 ms 级完成且内存不飙）
	assert.NoError(t, ValuesToMessage(msg, values))
}

func TestFlattenKeyEdgeCases(t *testing.T) {
	// 这些过去会 panic，修复后应当 silent skip（没有匹配字段）
	cases := []string{
		"=x",
		"[]=x",
		"[][]=x",
		"[=x",
		"]=x",
	}
	for _, c := range cases {
		values, _ := url.ParseQuery(c)
		msg := new(DemoRequest)
		assert.NotPanics(t, func() { _ = ValuesToMessage(msg, values) }, "input=%q", c)
	}
}
