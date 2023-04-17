package postgresql

import (
	"github.com/stretchr/testify/suite"
	"telegram-bot/solte.lab/pkg/config"
	"telegram-bot/solte.lab/pkg/storage"
	"testing"
)

type storageTearDown interface {
	Close() error
	DropTables() error
}

type querySuiteTests struct {
	suite.Suite
	storage storage.Storage
}

func (s *querySuiteTests) SetupSuite() {
	conf := config.NewTestConfig()
	st, err := New(conf.PostgreSQL)
	s.Suite.NoError(err)
	s.storage = st
}

func (s *querySuiteTests) Test_Save() {
	p := &storage.Page{
		UserName: "Pupu Romeo",
		URL:      "https://www.solte.com/labs1",
	}

	s.Suite.Run("Save", func() {
		err := s.storage.Save(p)
		s.Suite.NoError(err)
	})

	s.Suite.Run("Is Exists", func() {
		isExist, err := s.storage.IsExist(p)
		s.Suite.NoError(err)
		s.Suite.True(isExist)
	})

	s.Suite.Run("Save again", func() {
		err := s.storage.Save(p)
		s.Suite.Error(err)
	})

	s.Suite.Run("Pick a Random link", func() {
		randomPage, err := s.storage.PickRandom(p.UserName)
		s.Suite.NoError(err)
		s.Suite.NotNil(p)
		p = randomPage
	})

	s.Suite.Run("Remove", func() {
		err := s.storage.Remove(p)
		s.Suite.NoError(err)
	})
}

func (s *querySuiteTests) TearDownSuite() {
	err := s.storage.(storageTearDown).DropTables()
	s.Suite.NoError(err)

	err = s.storage.(storageTearDown).Close()
	s.Suite.NoError(err)
}

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(querySuiteTests))
}
