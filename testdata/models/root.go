package models

type PrimaryID struct {
	// 自增 ID
	ID uint64 `db:"f_id,autoincrement" json:"-"`
}

type OperationTimes struct {
	// 创建时间
	CreatedAt uint64 `db:"f_created_at,default='0'" json:"createdAt" `
	// 更新时间
	UpdatedAt uint64 `db:"f_updated_at,default='0'" json:"updatedAt"`
}

type OperationTimesWithDeletedAt struct {
	OperationTimes
	// 删除时间
	DeletedAt uint64 `db:"f_deleted_at,default='0'" json:"-"`
}
