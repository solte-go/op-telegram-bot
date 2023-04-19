package cache

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"sync"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage/storageWrapper/postgresql"
	"time"
)

type dbSync interface {
	GetAllUsers() (users []models.User, err error)
}

type userCache struct {
	user     models.User
	syncTime time.Time
}

type Container struct {
	userCache map[string]*userCache
	logger    *zap.Logger
	sync      dbSync
	mx        sync.RWMutex
}

var ErrorNoUserForUpdate = errors.New("no user for update")

func New(ctx context.Context, st *postgresql.Storage) *Container {

	c := &Container{
		userCache: make(map[string]*userCache),
		sync:      st,
		logger:    zap.L().Named("cache"),
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				c.syncUsersCache()
			case <-ctx.Done():
				ticker.Stop()
				c.logger.Warn("received context done, stopping storageWrapper sync")
				return
			}
		}
	}()

	return c
}

func (c *Container) AddUser(user *models.User) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.userCache[user.Name] = &userCache{
		user:     *user,
		syncTime: time.Now(),
	}
}

func (c *Container) GetUser(name string) (models.User, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	cachedUser, ok := c.userCache[name]
	if !ok {
		return models.User{}, false
	}

	return cachedUser.user, true
}

func (c *Container) UpdateUser(user *models.User) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	cachedUser, ok := c.userCache[user.Name]
	if ok {
		cachedUser.user = *user
		return nil
	}

	return ErrorNoUserForUpdate
}

func (c *Container) syncUsersCache() {
	c.mx.Lock()

	users, err := c.sync.GetAllUsers()
	if err != nil {
		c.logger.Error("can't sync users storageWrapper", zap.Error(err))
		return
	}

	for _, user := range users {
		cachedUser, ok := c.userCache[user.Name]
		if ok {
			cachedUser.user = user
			cachedUser.syncTime = time.Now()
		}
	}

	c.mx.Unlock()

	c.clearUp()

	c.logger.Debug("sync users cache completed")
}

func (c *Container) clearUp() {
	c.mx.Lock()
	defer c.mx.Unlock()

	for user := range c.userCache {
		if time.Now().Sub(c.userCache[user].syncTime) > time.Hour {
			delete(c.userCache, user)
		}
	}
}
