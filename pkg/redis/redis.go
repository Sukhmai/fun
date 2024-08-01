package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{
		rdb: rdb,
	}
}

func (c *RedisClient) Close() error {
	err := c.rdb.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) SaveAnswer(userId, answer string, questionNum int) error {
	var ctx = context.Background()

	// err := c.rdb.Set(ctx, userKey, answer, 0).Err()
	err := c.rdb.HSet(ctx, "user:"+userId, fmt.Sprintf("question:%d", questionNum), answer).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) GetAnswer(userId string, questionNum int) (string, error) {
	var ctx = context.Background()

	val, err := c.rdb.HGet(ctx, getUserKey(userId), getQuestionKey(questionNum)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *RedisClient) GetAnswers(userId string) (map[string]string, error) {
	var ctx = context.Background()

	val, err := c.rdb.HGetAll(ctx, getUserKey(userId)).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *RedisClient) SaveReply(userId1, userId2 string, questionNum, replyNum int, reply string) error {
	var ctx = context.Background()

	userKey := getUserReplyKey(userId1, userId2, questionNum, replyNum)
	err := c.rdb.Set(ctx, userKey, reply, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) GetReply(userId1, userId2 string, questionNum, replyNum int) (string, error) {
	var ctx = context.Background()

	userKey := getUserReplyKey(userId1, userId2, questionNum, replyNum)
	val, err := c.rdb.Get(ctx, userKey).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func getUserKey(userId string) string {
	return fmt.Sprintf("user:"+userId, userId)
}

func getQuestionKey(questionNum int) string {
	return fmt.Sprintf("question:%d", questionNum)
}

func getUserReplyKey(userId1, userId2 string, questionNum, replyNum int) string {
	return fmt.Sprintf("users:%s:%s:answer%d:reply%d", userId1, userId2, questionNum, replyNum)
}
