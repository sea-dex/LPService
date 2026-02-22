package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrDataNotExists = errors.New("data does not exist")
	ErrCacheCMD      = errors.New("cache command execute failed")
)

type Store interface {
	GetVal(ctx context.Context, key string, value interface{}) error
	SetVal(ctx context.Context, key string, value interface{}, expires time.Duration) error
	DeleteVal(ctx context.Context, key string) error
}

type RedisClient struct {
	*redis.Client
}

func InitCache(addr, password string, db int) (*RedisClient, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       db,
		Password: password,
	})

	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisClient{
		cli,
	}, nil
}

func (cli *RedisClient) GetVal(ctx context.Context, key string, value interface{}) error {
	val, err := cli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return ErrDataNotExists
	} else if err != nil {
		return fmt.Errorf("get cache: %w: %w", ErrCacheCMD, err)
	}

	return deserialize([]byte(val), value)
}

func (cli *RedisClient) SetVal(ctx context.Context, key string, value interface{}, expires time.Duration) error {
	bs, err := serialize(value)
	if err != nil {
		return err
	}

	err = cli.Set(ctx, key, bs, expires).Err()
	if err != nil {
		return fmt.Errorf("set cache: %w: %w", ErrCacheCMD, err)
	}

	return nil
}

func (cli *RedisClient) DeleteVal(ctx context.Context, key ...string) error {
	err := cli.Del(ctx, key...).Err()
	if err != nil {
		return fmt.Errorf("delete cache: %w: %w", ErrCacheCMD, err)
	}

	return nil
}

func serialize(value interface{}) ([]byte, error) {
	if bs, ok := value.([]byte); ok {
		return bs, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	default:
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)

	if err := encoder.Encode(value); err != nil { // 编码
		return nil, fmt.Errorf("encode: %w", err)
	}

	return b.Bytes(), nil
}

func deserialize(byt []byte, ptr interface{}) error {
	if bs, ok := ptr.(*[]byte); ok {
		*bs = byt

		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() { //nolint:exhaustive
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err := strconv.ParseInt(string(byt), 10, 64)

			if err != nil {
				return fmt.Errorf("parse int: %w", err)
			} else {
				p.SetInt(i)
			}

			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err := strconv.ParseUint(string(byt), 10, 64)

			if err != nil {
				return fmt.Errorf("parse uint: %w", err)
			} else {
				p.SetUint(i)
			}

			return nil
		default: // Other types decode by gob
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)

	if err := decoder.Decode(ptr); err != nil { // 解码
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
