package tasks

import (
	"database/sql"
)

type PGStore struct {
	DB *sql.DB
}

func (ps *PGStore) Insert(newtask *NewTask) (*Task, error) {
	t := newtask.ToTask()
	tx, err := ps.DB.Begin() // Transactions are important
	if err != nil {
		return nil, err
	}

	// Parameter markers $1$2 automatically encode parameters to be sanatized
	sql := `insert into tasks 
	(title, createdAt, modifiedAt, complete)
	values ($1,$2,$3,$4) returning id`
	// Query row gets the just inserted row back
	row := tx.QueryRow(sql, t.Title, t.CreatedAt, t.ModifiedAt, t.Complete) // Specifying the values for parameter markers
	err = row.Scan(&t.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	sql = `insert into tags(taskID, tag
	values($1,$2))`
	for _, tag := range t.Tags {
		// Exec doesn't get the inserted row back
		_, err := tx.Exec(sql, t.ID, tag)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return t, nil

	// Would be better to use prepared statements
}

func Get(ID interface{}) (*Task, error) {
	// Can use a function in postgress to recieve multiple result sets

	return nil, nil
}
func GetAll() ([]*Task, error) {
	return nil, nil
}
func Update(task *Task) error {
	return nil
}
