package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"orders-api/models"
)

type RedisRepo struct {
	Client *redis.Client
}

var ErrNotExists = errors.New("order not exists")

func orderIdKey(id int64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to encoder order %w", err)
	}

	key := orderIdKey(order.OrderID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)

	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to insert order %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add orders to set %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute transaction %w", err)
	}
	return nil
}

func (r *RedisRepo) FindByID(ctx context.Context, id int64) (models.Order, error) {
	key := orderIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return models.Order{}, ErrNotExists
	} else if err != nil {
		return models.Order{}, fmt.Errorf("failed to find order by id  %w", err)
	}

	var order models.Order

	if err := json.Unmarshal([]byte(value), order); err != nil {
		return models.Order{}, fmt.Errorf("failed to decode order %w", err)
	}
	return order, nil
}

func (r *RedisRepo) DeleteById(ctx context.Context, id int64) error {
	key := orderIdKey(id)

	txn := r.Client.TxPipeline()
	err := txn.Del(ctx, key).Err()

	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExists
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order by id %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute transaction %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order models.Order) error {
	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to encoder order %w", err)
	}

	key := orderIdKey(order.OrderID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()

	if errors.Is(err, redis.Nil) {
		return ErrNotExists
	} else if err != nil {
		return fmt.Errorf("failed to update order %w,", err)
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Orders []models.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()

	if len(keys) == 0 {
		return FindResult{
			Orders: make([]models.Order, 0),
		}, nil
	}
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids %w", err)
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()

	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders %w", err)
	}

	orders := make([]models.Order, len(xs))

	for i, x := range xs {
		x := x.(string)

		var order models.Order
		err := json.Unmarshal([]byte(x), &order)

		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json :%w", err)
		}
		orders[i] = order
	}
	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
