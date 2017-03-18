package redis

import (
	"os"

	"github.com/garyburd/redigo/redis"
)

type RedisCli struct {
	conn redis.Conn
}

var instanceRedisCli *RedisCli = nil

func Connect() (conn *RedisCli) {
	if instanceRedisCli == nil {
		instanceRedisCli = new(RedisCli)
		var err error
		redisAddr := os.Getenv("YOUFIE_REDIS_ADDR")
		if len(redisAddr) == 0 {
			if os.Getenv("GET_HOSTS_FROM") == "env" {
				redisAddr = os.Getenv("REDIS_MASTER_SERVICE_HOST")
			} else if os.Getenv("GET_HOSTS_FROM") == "dns" {
				redisAddr = "redis-master"
			}
		}

		instanceRedisCli.conn, err = redis.Dial("tcp", redisAddr+":6379")

		if err != nil {
			panic(err)
		}

		//if _, err := instanceRedisCli.conn.Do("AUTH", "trevorprater"); err != nil {
		//		instanceRedisCli.conn.Close()
		//		panic(err)
		//	}
	}

	return instanceRedisCli
}

func (redisCli *RedisCli) SetValue(key string, value string, expiration ...interface{}) error {
	_, err := redisCli.conn.Do("SET", key, value)

	if err == nil && expiration != nil {
		redisCli.conn.Do("EXPIRE", key, expiration[0])
	}

	return err
}

func (redisCli *RedisCli) GetValue(key string) (interface{}, error) {
	return redisCli.conn.Do("GET", key)
}
