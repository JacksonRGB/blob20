package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type Mint struct {
	ID           uint   `gorm:"primaryKey"`
	TxHash       string `gorm:"type:varchar(255);uniqueIndex;not null;comment:所在交易"`
	Owner        string `gorm:"type:varchar(255);index;not null;comment:owner"`
	Block        int    `gorm:"type:int;not null;comment:block"`
	BlockTime    int    `gorm:"type:int;not null;comment:block time"`
	BlobGas      uint   `gorm:"type:bigint;not null;comment:blob gas"` // blob消耗的eth
	BlobGasPrice uint   `gorm:"type:bigint;not null;comment:blob gas price"`
	TxGas        uint   `gorm:"type:bigint;not null;comment:tx gas"` // tx消耗的eth
	TotalGas     uint   `gorm:"type:bigint;not null;comment:total gas"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Height struct {
	Key      string `gorm:"primaryKey"`
	IntValue int    `gorm:"type:int;not null"` // 配置value
}
