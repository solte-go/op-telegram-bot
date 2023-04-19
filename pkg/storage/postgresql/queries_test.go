package postgresql

import (
	"github.com/stretchr/testify/suite"
	"telegram-bot/solte.lab/pkg/config"
	"telegram-bot/solte.lab/pkg/models"
	"testing"
)

type storageTests interface {
	Close() error
	DropTables() error
	InsertUser(user *models.User) (err error)
}

type querySuiteTests struct {
	suite.Suite
	storage storageTests
}

func (s *querySuiteTests) SetupSuite() {
	conf := config.NewTestConfig()
	st, err := New(conf.PostgreSQL)
	s.Suite.NoError(err)
	s.storage = st
}

func (s *querySuiteTests) Test_Save() {
	u := &models.User{
		Name:     "test",
		Language: "",
		Topic:    "",
	}

	s.Suite.Run("SaveUser", func() {
		err := s.storage.InsertUser(u)
		s.Suite.NoError(err)
	})

	//s.Suite.Run("Is Exists", func() {
	//	isExist, err := s.storage.IsExist(p)
	//	s.Suite.NoError(err)
	//	s.Suite.True(isExist)
	//})
	//
	//s.Suite.Run("Save again", func() {
	//	err := s.storage.Save(p)
	//	s.Suite.Error(err)
	//})
	//
	//s.Suite.Run("Pick a Random link", func() {
	//	randomPage, err := s.storage.PickRandom(p.UserName)
	//	s.Suite.NoError(err)
	//	s.Suite.NotNil(p)
	//	p = randomPage
	//})
	//
	//s.Suite.Run("Remove", func() {
	//	err := s.storage.Remove(p)
	//	s.Suite.NoError(err)
	//})
}

func (s *querySuiteTests) TearDownSuite() {
	err := s.storage.DropTables()
	s.Suite.NoError(err)

	err = s.storage.Close()
	s.Suite.NoError(err)
}

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(querySuiteTests))
}
