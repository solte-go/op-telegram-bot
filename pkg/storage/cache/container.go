package cache

import (
	"context"
	"errors"
	"sync"
	"telegram-bot/solte.lab/pkg/storage/emsql"
	"time"

	"go.uber.org/zap"
	"telegram-bot/solte.lab/pkg/models"
)

type dbSync interface {
	GetAllUsers() (users []models.User, err error)
}

type userCache struct {
	user     *models.User
	syncTime time.Time
}

type Container struct {
	userCache map[string]*userCache
	logger    *zap.Logger
	sync      dbSync
	mx        sync.RWMutex
}

var ErrorNoUserForUpdate = errors.New("no user for update")
var ErrorUserEmpty = errors.New("user should not be is empty")

func New(ctx context.Context, st emsql.OPContract) *Container {
	c := &Container{
		userCache: make(map[string]*userCache),
		sync:      st.User(),
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
		user:     user,
		syncTime: time.Now(),
	}
}

func (c *Container) GetUser(name string) (*models.User, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	cachedUser, ok := c.userCache[name]
	if !ok {
		return nil, false
	}

	return cachedUser.user, true
}

func (c *Container) UpdateUserWithUpset(user *models.User) error {
	if user == nil {
		return ErrorUserEmpty
	}
	err := c.UpdateUser(user)
	if err != nil {
		c.AddUser(user)
	}
	return nil
}

func (c *Container) UpdateUser(user *models.User) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	cachedUser, ok := c.userCache[user.Name]
	if ok {
		cachedUser.user = user
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
			cachedUser.user.Topic = user.Topic
			cachedUser.user.Language = user.Language
			cachedUser.user.Offset = user.Offset
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
