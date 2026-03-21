package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//начало решения

//SQLMap представляет карту, которая хранится в SQL-базе данных
// SQLMap представляет карту, которая хранится в SQL-базе данных
type SQLMap struct{
	Map map[string]any
	db *sql.DB
	GetStmt *sql.Stmt
	SetStmt *sql.Stmt
	DelStmt *sql.Stmt
	Timeout time.Duration
}

// NewSQLMap создает новую SQL-карту в указанной базе
func NewSQLMap(db *sql.DB) (*SQLMap, error) {
	query := `create table if not exists map(key text primary key, val blob)`
	_, err := db.Exec(query)
	if err != nil {
		return nil, err
	}
	get, err := db.Prepare(`select val from map where key = ?`) 
	if err != nil {
		return nil, err
	}
	set, err := db.Prepare(`insert into map(key, val) values (?, ?)
on conflict (key) do update set val = excluded.val`) 
	if err != nil {
		return nil, err
	}
	del, err := db.Prepare(`delete from map where key = ?`) 
	if err != nil {
		return nil, err
	}
	return &SQLMap{
		Map: map[string]any{},
		db: db,
		GetStmt: get,
		SetStmt: set,
		DelStmt: del,
		Timeout: 60 * time.Second,
	}, nil
}

// SetTimeout устанавливает максимальное время выполнения
// отдельного метода карты.
func (m *SQLMap) SetTimeout(d time.Duration) {
	m.Timeout = d
}

// Get возвращает значение для указанного ключа.
// Если такого ключа нет - возвращает ошибку sql.ErrNoRows.
func (m *SQLMap) Get(key string) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()

	row := m.GetStmt.QueryRowContext(ctx, key)
	var val any
	err := row.Scan(&val)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

// Set устанавливает значение для указанного ключа.
// Если такой ключ уже есть - затирает старое значение (это не считается ошибкой).
func (m *SQLMap) Set(key string, val any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	_, err := m.SetStmt.ExecContext(ctx, key, val)
	if err != nil {
		return err
	}
	return nil
}

// Delete удаляет запись карты с указанным ключом.
// Если такого ключа нет - ничего не делает (это не считается ошибкой).
func (m *SQLMap) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	res, err := m.DelStmt.ExecContext(ctx, key)
	if err != nil {
		return err
	}
	if count, _ := res.RowsAffected(); count == 0 {
		return nil
	}
	return nil
}

// SetItems устанавливает значения указанных ключей.
func (m *SQLMap) SetItems(items map[string]any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	defer cancel()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
        return err
    }
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, m.SetStmt)
	for key, val := range items {
		_, err := txStmt.ExecContext(ctx, key, val)
	if err != nil {
		return err
	}
	}
	return tx.Commit()
}

// Close освобождает ресурсы, занятые картой в базе.
func (m *SQLMap) Close() error {
	err := m.GetStmt.Close()
	if err != nil {
		return err
	}
	err = m.SetStmt.Close()
	if err != nil {
		return err
	}
	err = m.DelStmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// конец решения

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	m, err := NewSQLMap(db)
	if err != nil {
		panic(err)
	}
	defer m.Close()

	m.SetTimeout(10 * time.Millisecond)

	m.Set("name", "Alice")
	m.Get("name")
}
