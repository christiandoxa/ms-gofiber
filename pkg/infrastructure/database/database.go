package database

import (
	"sort"
	"sync"
	"time"
)

type TodoRecord struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DB struct {
	mutex sync.RWMutex
	todos map[string]TodoRecord
}

func Connect() *DB {
	return &DB{
		todos: map[string]TodoRecord{},
	}
}

func (db *DB) CreateTodo(record TodoRecord) TodoRecord {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.todos[record.ID] = record
	return record
}

func (db *DB) GetTodo(id string) (TodoRecord, bool) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	record, ok := db.todos[id]
	return record, ok
}

func (db *DB) ListTodos() []TodoRecord {
	db.mutex.RLock()
	defer db.mutex.RUnlock()
	records := make([]TodoRecord, 0, len(db.todos))
	for _, record := range db.todos {
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.Before(records[j].CreatedAt)
	})
	return records
}

func (db *DB) UpdateTodo(record TodoRecord) (TodoRecord, bool) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if _, ok := db.todos[record.ID]; !ok {
		return TodoRecord{}, false
	}
	db.todos[record.ID] = record
	return record, true
}

func (db *DB) DeleteTodo(id string) bool {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	if _, ok := db.todos[id]; !ok {
		return false
	}
	delete(db.todos, id)
	return true
}
