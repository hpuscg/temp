/*
#Time      :  2021/1/22 2:00 下午
#Author    :  chuangangshen@deepglint.com
#File      :  license.go
#Software  :  GoLand
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"temp/GoModelTest/sqlTest/util/time"
)

type License struct {
	TemplateID   uint       `json:"template_id"` // 模板ID
	Template     Template   `json:"template" gorm:"ForeignKey:TemplateID"`
	UserID       uint       `json:"user_id"` // 使用者ID
	User         LicUser    `json:"user" gorm:"ForeignKey:UserID"`
	UUID         string     `json:"uuid" gorm:"INDEX"` // 设备uuid
	ProductID    uint       `json:"product_id"`        // 产品ID
	Product      Product    `json:"product" gorm:"ForeignKey:ProductID"`
	LicenseType  int        `json:"license_type"` // 证书类型，0是app启用，1是重置密码
	Description  string     `json:"description"`
	License      []byte     `json:"license" gorm:"type:LONGBLOB"`
	CreateUserID uint       `json:"create_user_id"` // 创建证书人员ID
	CreateUser   LicUser    `json:"create_user" gorm:"ForeignKey:CreateUserID"`
	GoodID       uint       `json:"good_id"` // 批量签发ID
	Good         Good       `json:"good" gorm:"ForeignKey:GoodID"`
	DateNums     int        `json:"date_nums"` // 授权天数
	LicenseData  Entries    `json:"license_data" gorm:"type:VARCHAR(10000)"`
	BatchIssueID uint       `json:"batch_issue_id"`
	BatchIssue   BatchIssue `json:"batch_issue" gorm:"ForeignKey:BatchIssueID"`
	BaseModel
}

type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.LicTime
	UpdatedAt time.LicTime
	DeletedAt *time.LicTime `sql:"index"`
}

type BatchIssue struct {
	BaseModel
	Name        string    `json:"name" gorm:"unique"` // 批量签发名称
	TemplateID  uint      `json:"template_id"`        // 模板ID
	Template    Template  `json:"template" gorm:"ForeignKey:TemplateID"`
	UserID      uint      `json:"user_id"` // 使用者ID
	User        LicUser   `json:"user" gorm:"ForeignKey:UserID"`
	UUID        BatchUUID `json:"uuid" gorm:"type:LONGBLOB"` // 设备uuid
	ProductID   uint      `json:"product_id"`                // 产品ID
	Product     Product   `json:"product" gorm:"ForeignKey:ProductID"`
	LicenseType int       `json:"license_type"` // 证书类型，0是app启用，1是重置密码
	Description string    `json:"description"`
	DateNums    int       `json:"date_nums"` // 授权天数
	BatchData   Entries   `json:"batch_data" gorm:"type:VARCHAR(10000)"`
}

type Good struct {
	GoodData    Entries  `json:"good_data" gorm:"type:VARCHAR(10000)"` // 扩展数据
	DateNums    int      `json:"date_nums"`                            // 授权天数
	ProductID   uint     `json:"product_id"`                           // 产品ID
	Product     Product  `json:"product" gorm:"ForeignKey:ProductID"`
	UserID      uint     `json:"user_id"` // 使用者ID
	User        LicUser  `json:"user" gorm:"ForeignKey:UserID"`
	TemplateID  uint     `json:"template_id"` // 模板ID
	Template    Template `json:"template" gorm:"ForeignKey:TemplateID"`
	Nums        uint     `json:"nums"`                    // 批量数量
	Status      uint     `json:"status"`                  // 状态，1表示正常
	Description string   `json:"description"`             // 描述
	LicenseType int      `json:"license_type"`            // 证书类型
	UsedNum     uint     `json:"used_num"`                // 已使用数量
	GoodName    string   `json:"good_name" gorm:"unique"` // 商品名称
	BaseModel
}

type Product struct {
	Name        string `json:"name" gorm:"unique"`         // 产品名称
	DeviceModel string `json:"device_model" gorm:"unique"` // 产品型号
	Description string `json:"description"`
	LicenseNum  int    `json:"license_num" gorm:"-"`  // 该产品签发的证书数量
	TemplateNum int    `json:"template_num" gorm:"-"` // 该产品生成的模板数量
	BaseModel
}

type LicUser struct {
	UserName    string  `json:"user_name" gorm:"INDEX;UNIQUE;NOT NULL"` // 用户名
	Password    string  `json:"-" gorm:"NOT NULL"`                      // 密码
	RealName    string  `json:"real_name" grom:"UNIQUE;NOT NULL"`       // 昵称
	RoleID      uint    `json:"role_id" gorm:"NOT NULL"`                // 用户类型
	Role        LicRole `json:"role" gorm:"ForeignKey:RoleID"`
	InUse       uint    `json:"in_use"` // 0关闭用户，1启用用户
	Email       string  `json:"email"`  // 电子邮件
	Phone       string  `json:"phone"`  // 电话号码
	Creator     string  `json:"creator"`
	Description string  `json:"description"`
	BaseModel
}

type LicRole struct {
	BaseModel
	Name        string `json:"name" gorm:"INDEX;UNIQUE"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
}

type Template struct {
	ProductID    uint    `json:"product_id"` // 产品ID
	Product      Product `json:"product" gorm:"ForeignKey:ProductID"`
	Name         string  `json:"name" gorm:"INDEX;UNIQUE"` // 模版名称
	DateNums     int     `json:"date_nums"`                // 授权天数
	LicenseType  int     `json:"license_type"`             // 0是app启用，1是重置密码
	TemplateData Entries `json:"template_data" gorm:"type:VARCHAR(10000)"`
	Description  string  `json:"description"`
	BaseModel
}

type Entry struct {
	Type        string `json:"type"`  // 指定数据类型
	Value       string `json:"value"` // 数据
	Fixed       bool   `json:"fixed"` // 是否可编辑
	Description string `json:"description"`
	ToolTip     string `json:"tool_tip"` // 提示内容
}

type Entries map[string]Entry

func (this Entries) Value() (driver.Value, error) {
	b, e := json.Marshal(this)
	return string(b), e
}

func (this *Entries) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		return json.Unmarshal([]byte(value.(string)), this)
	default:
		return json.Unmarshal(value.([]byte), this)
	}
}

type BatchUUID []string

func (this BatchUUID) Value() (driver.Value, error) {
	b, e := json.Marshal(this)
	return string(b), e
}

func (this *BatchUUID) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		return json.Unmarshal([]byte(value.(string)), this)
	default:
		return json.Unmarshal(value.([]byte), this)
	}
}
