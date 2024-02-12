package mysql

import (
	"DAKExDUCK/assignment_1/pkg/models"
	"database/sql"
	"errors"
)

type UserDB struct {
	DB *sql.DB
}

func (model *UserDB) Get(mail string) (map[string]models.User, error) {
	query := `
	SELECT
		id, mail, password
	FROM
		users
	WHERE
		mail = $1;`
	users := map[string]models.User{}
	row := model.DB.QueryRow(query, mail)
	user := models.User{}

	err := row.Scan(&user.ID, &user.Mail, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, models.ErrNoRecord
		} else {
			return users, err
		}
	}

	users[user.Mail] = user
	return users, nil
}

func (model *UserDB) GetUserPermissions(email string) ([]string, error) {
	query := `
	SELECT
		permission
	FROM
		user_permissions
	WHERE
		user_email = $1;`

	rows, err := model.DB.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (model *UserDB) SetPermissions(email string, permissions []string) error {
	tx, err := model.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("DELETE FROM user_permissions WHERE user_email = $1", email)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO user_permissions (user_email, permission) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, permission := range permissions {
		_, err = stmt.Exec(email, permission)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (model *UserDB) Insert(mail string, password string) (int, error) {
	query := `
	INSERT INTO
		users (mail, password)
	VALUES($1, $2)
	RETURNING id
	`
	var id int
	err := model.DB.QueryRow(query, mail, password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
