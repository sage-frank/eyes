package domain

import "time"

type Azure struct {
	ID      int64
	Author  int64
	Title   string
	Content string
}

type UserInfo struct {
	ID         int64
	RealName   string
	Email      string
	openID     string
	valid      int64
	CreateDate time.Time
}
