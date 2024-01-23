package redisfold

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	inst *redis.Client
	ctx  context.Context
}

func GetRedisInstance() *RedisInstance {
	return &RedisInstance{redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}), context.Background(),
	}
}

func (r *RedisInstance) BindUserToRoom(wsadress string, room string) error {
	err := r.inst.Set(r.ctx, wsadress, room, 0).Err()
	if err != nil {
		return fmt.Errorf("Value to redis cannot be set")
	}
	return nil
}

func (r *RedisInstance) GetUsersRoom(wsadress string) (string, error) {
	value, err := r.inst.Get(r.ctx, wsadress).Result()
	return value, err
}

func (r *RedisInstance) UnbindUserToRoom(wsadress string) error {
	err := r.inst.Del(r.ctx, wsadress).Err()
	if err != nil {
		return fmt.Errorf("Value to redis cannot be set")
	}
	return nil
}

// select * from room1 => ["wsAdr1":self_id1,"wsAdr2":self_id2]
func (r *RedisInstance) SetRoomClient(roomID string, wsAdr string, self_id string) {
	userSession := r.getRoomClients(roomID)
	userSession[wsAdr] = self_id
	for k, v := range userSession {
		r.inst.HSet(r.ctx, roomID, k, v)
	}
}

func (r *RedisInstance) DeleteRoomsClient(roomID string, self_id string) {
	userSession := r.getRoomClients(roomID)
	for k, v := range userSession {
		if v == self_id {
			continue
		}
		r.inst.HSet(r.ctx, roomID, k, v)
	}
}

// solid declaration. Non-exportable
func (r *RedisInstance) getRoomClients(roomID string) map[string]string {
	return r.inst.HGetAll(r.ctx, roomID).Val()
}

// user address is gotten rid from room for broadcating
func (r *RedisInstance) PropagationList(roomID string, wsAdr string) map[string]string {
	allClients := r.getRoomClients(roomID)
	delete(allClients, wsAdr)
	return allClients
}
