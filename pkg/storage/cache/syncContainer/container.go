package syncContainer

import (
	"errors"
	"sync"
	"telegram-bot/solte.lab/pkg/models"
	"time"
)

type userCache struct {
	user     models.User
	syncTime time.Time
}

type Container struct {
	userCache map[string]*userCache
	mx        sync.RWMutex
}

var ErrorNoUserForUpdate = errors.New("no user for update")

func New() *Container {

	c := &Container{userCache: make(map[string]*userCache)}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ticker.C:
				c.clearUp()
			}
			//TODO ctx
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

func (c *Container) UpdateUser(user models.User) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	cachedUser, ok := c.userCache[user.Name]
	if ok {
		cachedUser.user = user
		return nil
	}

	return ErrorNoUserForUpdate
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
