package models

import (
	github_com_snail_tools_strcase "github.com/snail-tools/strcase"
)

func (Account) TableName() string {
	return "t_" + github_com_snail_tools_strcase.SnakeCase("Account")
}
func (Account) FieldKeyUserID() string {
	return "UserID"
}

func (Account) FieldKeyName() string {
	return "Name"
}

func (Account) FieldKeyEmail() string {
	return "Email"
}

func (Account) FieldKeyPassword() string {
	return "Password"
}

func (Account) Primary() []string {
	return []string{
		"ID",
	}
}

func (Account) Indexes() map[string][]string {
	return map[string][]string{
		"i_org_id": []string{
			"OrgID",
		},
		"i_user_id": []string{
			"UserID",
		},
	}
}
