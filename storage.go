package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	createStockTable() error
	CreateStock(*DailyStock) error
	DeleteStock(string) error
	GetStocks() ([]*DailyStock, error)
	GetStockByTicker(string) ([]*DailyStock, error)
}
type DailyStock struct {
	AfterHours float64 `json:"afterHours"`
	Close      float64 `json:"close"`
	FromVal    string  `json:"from"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Open       float64 `json:"open"`
	PreMarket  float64 `json:"preMarket"`
	Status     string  `json:"status"`
	Symbol     string  `json:"symbol"`
	Volume     float64 `json:"volume"`
}

func newDailyStock(afterHours float64, close float64, from string, high float64, low float64, open float64, preMarket float64, status string, symbol string, volume float64) *DailyStock {
	d := DailyStock{
		AfterHours: afterHours,
		Close:      close,
		FromVal:    from,
		High:       high,
		Low:        low,
		Open:       open,
		PreMarket:  preMarket,
		Status:     status,
		Symbol:     symbol,
		Volume:     volume,
	}
	return &d
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "postgres://postgres:postgrespw@localhost:32772?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
func (s *PostgresStore) Init() error {
	return s.createStockTable()
}

func (s *PostgresStore) createStockTable() error {
	query := `create table if not exists stocks (
		AfterHours varchar(100),
		Close varchar(100),
		FromVal varchar(100),
		High varchar(100),
		Low varchar(100),
		Open varchar(100),
		PreMarket varchar(100), 
		Status varchar(100),
		Symbol varchar(100),
		Volume varchar(100)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateStock(stock *DailyStock) error {
	query := `insert into stocks 
	(AfterHours, Close, FromVal, High, Low, Open, PreMarket, Status, Symbol, Volume)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := s.db.Query(
		query,
		stock.AfterHours,
		stock.Close,
		stock.FromVal,
		stock.High,
		stock.Low,
		stock.Open,
		stock.PreMarket,
		stock.Status,
		stock.Symbol,
		stock.Volume)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*DailyStock) error {
	return nil
}

func (s *PostgresStore) DeleteStock(ticker string) error {
	_, err := s.db.Query("delete from account where Symbol = $1", ticker)
	return err
}

func (s *PostgresStore) GetStocksByDay(number int) (*DailyStock, error) {
	rows, err := s.db.Query("select * from account where Symbol = [%d]", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanStockIn(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)
}

func (s *PostgresStore) GetStockByTicker(symbol string) ([]*DailyStock, error) {
	rows, err := s.db.Query("select * from stocks where Symbol = [%d]", symbol)
	if err != nil {
		return nil, err
	}
	stockProfile := []*DailyStock{}
	for rows.Next() {
		stock, err := scanStockIn(rows)
		if err != nil {
			return nil, err
		}
		stockProfile = append(stockProfile, stock)
	}
	if len(stockProfile) > 1 {
		return stockProfile, nil
	}

	return nil, fmt.Errorf("account %d not found", symbol)
}

func (s *PostgresStore) GetStocks() ([]*DailyStock, error) {
	rows, err := s.db.Query("select * from stocks")
	if err != nil {
		return nil, err
	}

	stockProfile := []*DailyStock{}
	for rows.Next() {
		stock, err := scanStockIn(rows)
		if err != nil {
			return nil, err
		}
		stockProfile = append(stockProfile, stock)
	}

	return stockProfile, nil
}

func scanStockIn(rows *sql.Rows) (*DailyStock, error) {
	stock := new(DailyStock)
	err := rows.Scan(
		&stock.AfterHours,
		&stock.Close,
		&stock.FromVal,
		&stock.High,
		&stock.Low,
		&stock.Open,
		&stock.PreMarket,
		&stock.Status,
		&stock.Symbol,
		&stock.Volume)

	return stock, err
}
