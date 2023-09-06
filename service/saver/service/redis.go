package service

import "github.com/go-redis/redis"

//redis存储文件名-哈希键值对

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := Redis.Ping().Err()
	if err != nil {
		panic(err)
	}
}

// 存储文件名-哈希键值对
func storeFileNameAndHash(fileName, hash string) error {
	err := Redis.HSet("file_name", fileName, hash).Err()
	if err != nil {
		return err
	}
	return nil
}
