package playground

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/immanuelhume/gohm"
	"github.com/immanuelhume/gohm/playground/stubs"
)

type User struct {
	data   gohm.EntityData
	client *redis.Client
}

func (u *User) Create(ctx context.Context, name string, age int) stubs.User {
	_new := stubs.User{name, age}
	u.client.HSet(ctx, "user:1")
	return _new
}
