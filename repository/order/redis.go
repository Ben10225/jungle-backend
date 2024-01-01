package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"jungle-proj/structs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx *gin.Context, reserveData structs.ReserveData) error {
	data, err := json.Marshal(reserveData)
	if err != nil {
		return fmt.Errorf("fail to encoder order: %w", err)
	}

	key := orderIDKey(11)

	res := r.Client.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		return fmt.Errorf("fail to set: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepo) FindByID(ctx *gin.Context, id uint64) (structs.ReserveData, error) {
	key := orderIDKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return structs.ReserveData{}, ErrNotExist
	} else if err != nil {
		return structs.ReserveData{}, fmt.Errorf("get order: %w", err)
	}

	var reserveData structs.ReserveData
	err = json.Unmarshal([]byte(value), &reserveData)
	if err != nil {
		return structs.ReserveData{}, fmt.Errorf("failed to decode reserveData json: %w", err)
	}

	return reserveData, nil
}
