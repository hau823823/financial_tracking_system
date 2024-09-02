package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisSet(t *testing.T) {
	db, mock := redismock.NewClientMock() // 使用 redismock 模擬 Redis 客戶端
	redisCache := &Redis{client: db}

	ctx := context.Background()
	key := "test_key"
	value := "test_value"
	expiration := 10 * time.Second

	// 設置模擬行為，當調用 SET 時返回 nil 表示成功
	mock.ExpectSet(key, value, expiration).SetVal("OK")

	err := redisCache.Set(ctx, key, value, expiration)
	assert.NoError(t, err)
	mock.ExpectationsWereMet() // 確保所有預期行為都已發生
}

func TestRedisGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisCache := &Redis{client: db}

	ctx := context.Background()
	key := "test_key"
	expectedValue := "test_value"

	// 設置模擬行為，當調用 GET 時返回預期的值
	mock.ExpectGet(key).SetVal(expectedValue)

	value, err := redisCache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
	mock.ExpectationsWereMet()
}

func TestRedisDelete(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisCache := &Redis{client: db}

	ctx := context.Background()
	key := "test_key"

	// 設置模擬行為，當調用 DEL 時返回成功
	mock.ExpectDel(key).SetVal(1)

	err := redisCache.Delete(ctx, key)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestRedisClose(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisCache := &Redis{client: db}

	// 設置模擬行為，當調用 Close 時返回成功
	//mock.ExpectClose().SetVal(nil)

	err := redisCache.Close()
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}
