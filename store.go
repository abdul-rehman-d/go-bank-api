package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	UpdateACcount(*Account) (*Account, error)
	DeleteAccount(int) (*Account, error)
	GetAccount(int) (*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func (s *PostgresStorage) Init() error {
	err := s.createAccountTable()
	return err
}

func (s *PostgresStorage) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id					serial			primary key,
		firstName 	varchar(50),
		lastName		varchar(50),
		number			serial,
		balance			numeric,
		created_at	timestamp
 	);`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "postgres://root:postgres@localhost:5433/go-chat?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db,
	}, nil
}

func (s *PostgresStorage) CreateAccount(acc *Account) (*Account, error) {
	insertQuery := `
	insert into account (firstName, lastName, number, balance, created_at)
	values($1, $2, $3, $4, $5)
	RETURNING id, firstName, lastName, number, balance, created_at`

	createdAccount := &Account{}

	err := s.db.QueryRow(
		insertQuery,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	).Scan(
		&createdAccount.ID,
		&createdAccount.FirstName,
		&createdAccount.LastName,
		&createdAccount.Number,
		&createdAccount.Balance,
		&createdAccount.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return createdAccount, nil
}
func (*PostgresStorage) UpdateACcount(acc *Account) (*Account, error) {
	return nil, nil
}
func (*PostgresStorage) DeleteAccount(id int) (*Account, error) {
	return nil, nil
}
func (*PostgresStorage) GetAccount(id int) (*Account, error) {
	return nil, nil
}
