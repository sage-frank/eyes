package domain

import (
	"database/sql"
)

type DSUser struct {
	ID              string       `json:"id"`                         // 用户唯一标识符
	Nickname        string       `json:"nickname,omitempty"`         // 用户昵称
	Username        string       `json:"username,omitempty"`         // 用户名
	Password        string       `json:"password,omitempty"`         // 密码
	Salt            string       `json:"salt,omitempty"`             // 盐
	Phone           string       `json:"phone,omitempty"`            // 用户电话号码
	Email           string       `json:"email,omitempty"`            // 用户电子邮箱
	Avatar          string       `json:"avatar,omitempty"`           // 用户头像URL
	Description     string       `json:"description,omitempty"`      // 用户描述或简介
	IsHot           int          `json:"is_hot,omitempty"`           // 是否是热门用户
	IsPro           int          `json:"is_pro,omitempty"`           // 假设这是专业用户标识
	CreatedIP       string       `json:"created_ip,omitempty"`       // 创建时的IP地址
	CreatedLocation string       `json:"created_location,omitempty"` // 创建时的位置
	CreatedAt       sql.NullTime `json:"created_at"`                 // 创建时间
	UpdatedAt       sql.NullTime `json:"updated_at"`                 // 最后更新时间
	Status          int          `json:"status,omitempty"`           // 用户状态
}
