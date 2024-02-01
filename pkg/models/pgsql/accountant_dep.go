package mysql

import (
	"DAKExDUCK/assignment_1/pkg/models"
	"database/sql"
	"errors"
)

type AccountDep struct {
	DB *sql.DB
}

func (model *AccountDep) GetAll() ([]models.AccountDep, error) {
	query := `
		SELECT id, name, sname, age
		FROM accountant_dep
	`

	rows, err := model.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []models.AccountDep{}
	for rows.Next() {
		account := models.AccountDep{}
		err := rows.Scan(&account.ID, &account.Name, &account.SecondName, &account.Age)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (model *AccountDep) Get(id int) ([]models.AccountDep, error) {
	query := `
	SELECT
		id, name, sname, age
	FROM
		public.accountant_dep
	WHERE
		id = $1;`
	row := model.DB.QueryRow(query, id)
	accounts := []models.AccountDep{}
	account := models.AccountDep{}

	err := row.Scan(&account.ID, &account.Name, &account.SecondName, &account.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	accounts = append(accounts, account)
	return accounts, nil
}

func (model *AccountDep) Insert(name string, sname string, age int) (int, error) {
	query := `
	INSERT INTO
		accountant_dep (name, sname, age)
	VALUES($1, $2, $3)
	RETURNING id
	`
	var id int
	err := model.DB.QueryRow(query, name, sname, age).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
