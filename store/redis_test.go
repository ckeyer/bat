package store

import (
	"testing"

	"gopkg.in/redis.v4"
)

var cli *redis.Client

func TestCodec(t *testing.T) {
	var a struct {
		Name    string
		Age     int
		Profile struct {
			Email string
		}
	}
	a.Name = "ckeyer"
	a.Age = 123
	a.Profile.Email = "adfasd@asdf.asdf"

	encode := func(v interface{}) {
		bs, err := Serialize(v)
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Logf("decode: (%T)%+v --> %s", v, v, bs)

		Deserializes(bs.([]byte), v)
		t.Logf("decode: %+v", v)

	}
	encode(a)
	encode(123)
	encode("123")
	encode([]byte("asdf"))
	encode([]string{"asdf", "cef"})
}
