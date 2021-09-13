package store

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)
// Store struct
type Store struct {
	config              *Config
	db                  *sqlx.DB
	rcurrencyRepository *RCurrencyRepository
}
// New object of store
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open connection to DB
func (s *Store) Open() error {
	db,err := sqlx.Open("sqlserver", s.config.DatabaseUrl)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

// Close connection to DB
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

// store.User().Create()
func (s *Store) RCurrency() *RCurrencyRepository {
	if s.rcurrencyRepository != nil {
		return s.rcurrencyRepository
	}
	s.rcurrencyRepository = &RCurrencyRepository{
		store: s,
	}
	return s.rcurrencyRepository
}
