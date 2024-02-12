package main

import (
	pgsql "DAKExDUCK/assignment_1/pkg/models/pgsql"
	"log"
)

type Application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	newsModel  *pgsql.NewsModel
	accountDep *pgsql.AccountDep
	users      *pgsql.UserDB
}
