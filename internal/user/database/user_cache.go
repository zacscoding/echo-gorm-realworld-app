package database

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"time"
)

func NewUserCacheDB(conf *config.Config, cli redis.UniversalClient, delegate UserDB) UserDB {
	return &userCache{
		conf:     conf,
		prefix:   conf.CacheConfig.Prefix,
		ttl:      conf.CacheConfig.TTL,
		cli:      cli,
		cache:    cache.New(&cache.Options{Redis: cli}),
		delegate: delegate,
	}
}

type userCache struct {
	conf     *config.Config
	prefix   string
	ttl      time.Duration
	cli      redis.UniversalClient
	cache    *cache.Cache
	delegate UserDB
}

func (uc *userCache) Save(ctx context.Context, u *userModel.User) error {
	if err := uc.delegate.Save(ctx, u); err != nil {
		return err
	}
	_ = uc.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   uc.getUserCacheKey(u.ID),
		Value: u,
		TTL:   uc.ttl,
	})
	return nil
}

func (uc *userCache) Update(ctx context.Context, u *userModel.User) error {
	if err := uc.delegate.Update(ctx, u); err != nil {
		return err
	}
	_ = uc.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   uc.getUserCacheKey(u.ID),
		Value: u,
		TTL:   uc.ttl,
		SetXX: true,
	})
	return nil
}

func (uc *userCache) FindByID(ctx context.Context, userID uint) (*userModel.User, error) {
	var find userModel.User
	err := uc.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   uc.getUserCacheKey(userID),
		Value: &find,
		TTL:   uc.ttl,
		Do: func(item *cache.Item) (interface{}, error) {
			return uc.delegate.FindByID(ctx, userID)
		},
	})
	if err != nil {
		return nil, err
	}
	return &find, nil
}

func (uc *userCache) FindByName(ctx context.Context, username string) (*userModel.User, error) {
	return uc.delegate.FindByName(ctx, username)
}

func (uc *userCache) FindByEmail(ctx context.Context, email string) (*userModel.User, error) {
	return uc.delegate.FindByEmail(ctx, email)
}

func (uc *userCache) Follow(ctx context.Context, userID, followerID uint) error {
	return uc.delegate.Follow(ctx, userID, followerID)
}

func (uc *userCache) IsFollow(ctx context.Context, userID, followerID uint) (bool, error) {
	return uc.delegate.IsFollow(ctx, userID, followerID)
}

func (uc *userCache) IsFollows(ctx context.Context, userID uint, followerIDs []uint) (map[uint]bool, error) {
	return uc.delegate.IsFollows(ctx, userID, followerIDs)
}

func (uc *userCache) UnFollow(ctx context.Context, userID, followerID uint) error {
	return uc.delegate.UnFollow(ctx, userID, followerID)
}

func (uc *userCache) FindFollowerIDs(ctx context.Context, userID uint) ([]uint, error) {
	return uc.delegate.FindFollowerIDs(ctx, userID)
}

func (uc *userCache) getUserCacheKey(id uint) string {
	return fmt.Sprintf("%susers.%d", uc.prefix, id)
}
