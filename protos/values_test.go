package protos

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryToMessage(t *testing.T) {
	query := "user[name]=foo&user.tags=GO&user.tags=RUST&user.attrs.age=20&user.attrs.height=180&friends[0].name=bar&friends[0].tags=PHP&friends[0].tags=JAVA&friends[0].attrs.age=18&friends[0].attrs.height=175"
	values, err := url.ParseQuery(query)
	assert.Nil(t, err)

	msg := new(DemoRequest)
	err = ValuesToMessage(msg, values)
	assert.Nil(t, err)
	b, _ := json.Marshal(msg)
	fmt.Println(string(b))
}

func TestMessageToQuery(t *testing.T) {
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
	fmt.Println("[query]", query)

	// 验证values
	msg2 := new(DemoRequest)
	err = ValuesToMessage(msg2, values)
	assert.Nil(t, err)
	b, _ := json.Marshal(msg2)
	fmt.Println("[message]", string(b))
}
