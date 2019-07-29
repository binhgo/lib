package lib

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"github.com/spf13/viper"
	"lib/Err"
	"lib/errors"
	"strconv"
	"sync"
	"time"
)

var (
	redisClientOnce sync.Once
	redisClient     *RedisClient
)

const (
	RKEY_DISPATCH_TASKLIST = "RKEY_DISPATCH_TASKLIST"
	RKEY_PRICE_DAILY       = "RKEY_PRICE_DAILY"
)

type RedisClient struct {
	client   *redis.Client
	Address  string
	Password string
}

func NewRedisClient(DB int) (*RedisClient) {

	redisClientOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis_host"),
			Password: viper.GetString("redis_pwd"),
			DB:       DB, // use default DB
		})

		_, err := client.Ping().Result()

		if err != nil {
			panic(err)
		}

		redisClient = &RedisClient{client, viper.GetString("redis_host"), viper.GetString("redis_pwd")}
	})

	return redisClient
}

func (this *RedisClient) Set(key string, val string, args ...int) error {
	expire := 48600 * 1000 //expire 单位秒
	for i, p := range args {
		switch i {
		case 0:
			expire = p
		}
	}

	expireDuration, _ := time.ParseDuration(strconv.Itoa(expire) + "s")
	err := this.client.Set(key, val, expireDuration).Err()
	if err != nil {
		panic(err)
	}

	return err
}

func (this *RedisClient) Get(key string) (string, error) {

	strCmd := this.client.Get(key)
	if (strCmd.Err() != nil && strCmd.Err().Error() == "redis: nil") {
		return "", strCmd.Err()
	}

	val, err := strCmd.Result()
	if err != nil {
		panic(err)
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, key+" does not exist")
	} else if err != nil {
		panic(err)
	} else {
		//fmt.Println(key, val)
	}

	return val, err
}

func (this *RedisClient) Del(key string) (error) {

	strCmd := this.client.Del(key)
	if (strCmd.Err() != nil && strCmd.Err().Error() == "redis: nil") {
		return strCmd.Err()
	}

	_, err := strCmd.Result()
	if err != nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, key+" remove error")
	}

	return err
}

func (this *RedisClient) HSet(key string, field string, value interface{}) (bool) {

	valStr, _ := json.Marshal(value)
	_, err := this.client.HSet(key, field, valStr).Result()
	if err != nil {
		panic(err)
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, "key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		//fmt.Println("key2", val)
	}

	return true
}

func (this *RedisClient) HSetStr(key string, field string, value string) (bool) {

	_, err := this.client.HSet(key, field, value).Result()
	if err != nil {
		panic(err)
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, "key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		//fmt.Println("key2", val)
	}

	return true
}

func (this *RedisClient) HMset(key string, fields map[string]interface{}) (bool) {

	_, err := this.client.HMSet(key, fields).Result()
	if err != nil {
		panic(err)
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, "key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		//fmt.Println("key2", val)
	}

	return true
}

func (this *RedisClient) HGetAll(key string) (map[string]string, error) {

	val, err := this.client.HGetAll(key).Result()
	if err != nil {
		panic(err)
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, "key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		//fmt.Println("key2", val)
	}

	return val, nil
}

func (this *RedisClient) HGet(key string, field string) (string, error) {

	val, err := this.client.HGet(key, field).Result()
	if err != nil {
		return val, err
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, "key2 does not exist")
	} else if err != nil {
		return val, err
	} else {
		//fmt.Println("key2", val)
	}

	return val, err
}

func (this *RedisClient) HDel(key string, field string) (error) {

	_, err := this.client.HDel(key, field).Result()

	if err != nil {
		return err
	}

	if err == redis.Nil {
		err = errors.NewUserError(Err.ERR_REDIS_ERR, "key2 does not exist")
	} else if err != nil {
		return err
	} else {
		//fmt.Println("key2", val)
	}

	return err
}

func (this *RedisClient) RPop(key string) (string, error) {

	StringCmd := this.client.RPop(key)
	if StringCmd.Err() != nil {
		return "", StringCmd.Err()
	}

	val, err := StringCmd.Result()
	return val, err
}

func (this *RedisClient) LIndex(key string, index int64) (string, error) {

	StringCmd := this.client.LIndex(key, index)
	if StringCmd.Err() != nil {
		return "", StringCmd.Err()
	}

	val, err := StringCmd.Result()
	return val, err
}

func (this *RedisClient) LRem(key string, count int64, value interface{}) (int64, error) {

	intCmd := this.client.LRem(key, count, value)
	if intCmd.Err() != nil {
		return 0, intCmd.Err()
	}

	val, err := intCmd.Result()
	return val, err
}

func (this *RedisClient) LLen(key string) (int64, error) {

	IntCmd := this.client.LLen(key)
	if IntCmd.Err() != nil {
		return 0, IntCmd.Err()
	}

	return IntCmd.Result()
}

func (this *RedisClient) RPush(key string, val interface{}) error {

	valBytes, err := jsoniter.Marshal(val)
	if err != nil {
		return errors.NewUserError(Err.ERR_REDIS_ERR, err.Error())
	}

	fmt.Println(string(valBytes))
	intCmd := this.client.RPush(key, valBytes)
	//intCmd = this.client.RPush(key, []byte("11111"))
	if intCmd.Err() != nil {
		return intCmd.Err()
	}

	return nil
}

func (this *RedisClient) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return this.client.Expire(key, expiration)
}

func (this *RedisClient) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return this.client.ExpireAt(key, tm)
}

func (this *RedisClient) Destroy() {
	redisClientOnce = sync.Once{}
	this.client.Close();
}
