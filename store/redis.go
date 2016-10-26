package store

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v4"
)

const (
	MaxQuera = 100
)

type RedisCli struct {
	*redis.Client
	Prefix string
}

func NewRedisCli(addr, auth string) *RedisCli {
	return &RedisCli{
		redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: auth,
		}),
		"bat",
	}
}

func (r *RedisCli) LogHistory(kr Keyer) error {
	bs, err := Serialize(kr)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%slog/%s", r.Prefix, kr.Key())

	return r.Set(key, bs, 0).Err()
}

func (r *RedisCli) GetHistory(key string) ([]byte, error) {
	val := r.Get(fmt.Sprintf("%slog/%s", r.Prefix, key))
	if err := val.Err(); err != nil {
		return nil, err
	}

	return val.Bytes()
}

func (r *RedisCli) ListHistory(key string, opt ...int) ([]byte, error) {
	offset, limit := 0, MaxQuera
	switch len(opt) {
	case 0:
	case 1:
		offset = opt[0]
	case 2:
		offset, limit = opt[0], opt[1]
	default:
		return nil, fmt.Errorf("wrong length for list history logs.")
	}

	keysRet := r.Keys(fmt.Sprintf("%slog/%s*", r.Prefix, key))
	if err := keysRet.Err(); err != nil {
		return nil, nil
	}

	keys := keysRet.Val()
	if len(keys) <= offset {
		return nil, fmt.Errorf("offset is out range")
	}
	if limit <= 0 {
		limit = MaxQuera
	}
	log.Debugf("get keys(%s): %+v", key, keys)
	log.Debugf("offset: %v, limit: %v", offset, limit)
	if len(keys) > offset+limit {
		keys = keys[offset : offset+limit]
	} else {
		keys = keys[offset:]
	}

	vals := make([][]byte, 0, len(keys))
	for _, k := range keys {
		ret := r.Get(k)
		if ret.Err() != nil {
			return nil, ret.Err()
		}

		bs, err := ret.Bytes()
		if err != nil {
			return nil, err
		}

		vals = append(vals, bs)
	}

	return append([]byte("["), append(bytes.Join(vals, []byte(",")), []byte("]")...)...), nil
}

func Serialize(v interface{}) (interface{}, error) {
	return json.Marshal(v)
	// if _, ok := v.(encoding.BinaryMarshaler); ok {
	// 	// redis default.
	// 	return v, nil
	// }

	// switch reflect.TypeOf(v).Kind() {
	// case reflect.Struct:
	// case reflect.Map:
	// case reflect.Array:
	// case reflect.Slice:
	// default:
	// 	return "", fmt.Errorf("not support type %T of %+v", v, v)
	// }
	// return "", nil
}

func Deserializes(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
	// if bu, ok := v.(encoding.BinaryUnmarshaler); ok {
	// 	return bu.UnmarshalBinary(data)
	// }
	// return nil
}
