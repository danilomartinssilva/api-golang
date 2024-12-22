package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil

}

func (s *PostgresStorage) CreateAccount(a *Account) error {
	insertSql := "INSERT INTO accounts (first_name, last_name, number, balance) VALUES ($1, $2, $3, $4)"
	resp, err := s.db.Exec(insertSql,
		a.FirstName, a.LastName, a.Number, a.Balance)

	if err != nil {
		return err
	}

	fmt.Print("%+v\n", resp)

	return nil
}

func (s *PostgresStorage) DeleteAccount(id int) error { return nil }

func (s *PostgresStorage) UpdateAccount(a *Account) error { return nil }

func (s *PostgresStorage) GetAccountById(id int) (*Account, error) { return nil, nil }

func (s *PostgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		number SERIAL,
		balance INT,
		created_at TIMESTAMP DEFAULT NOW()	 
		
	)`
	_, err := s.db.Exec(query)
	return err

}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {

	selectSql := "SELECT * FROM accounts"
	rows, err := s.db.Query(selectSql)

	if err != nil {
		return nil, err
	}

	//defer rows.Close()

	accounts := []*Account{}

	for rows.Next() {
		a := &Account{}
		if err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Number, &a.Balance, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, nil

}
