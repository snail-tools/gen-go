package models

import "github.com/snail-tools/gen-go/testdata/models/tmp"

//go:generate tools db -n wxf -t Account

// Account 账户
//@def primary f_id ID
//@def unique_index i_userID UserID
//@def index i_name Name
type Account struct {
	PrimaryID
	RefAccountID
	// 姓名
	Name string `db:"f_name,size=50,default=''"`
	// 密码
	Password string `db:"f_password"`
	// 用户ID
	Nickname string `db:"f_nick_name,size=90,default=''"`
	OperationTimes
	tmp.TmpDeep
}

type RefAccountID struct {
	// 账户ID
	AccountID uint64 `db:"f_account_id"`
}
