package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type User struct {
	ID           int
	Mail         string
	PasswordHash string
}

type Permission struct {
	key       string
	permitted bool
}

type AccountDep struct {
	ID         int
	Name       string
	SecondName string
	Age        int
}

type News struct {
	ID          int
	MoodleID    int
	Title       string
	Body        string
	Attachments []uint8
	Created     time.Time
	Tags        []Tag
}

type Tag struct {
	ID     int
	NameEN string
	NameRU string
}
