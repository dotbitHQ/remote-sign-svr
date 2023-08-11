package tables

import (
	"time"
)

type TableAddressInfo struct {
	Id           uint64       `json:"id" gorm:"column:id; primaryKey; type:bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '';"`
	AddrChain    AddrChain    `json:"addr_chain" gorm:"column:addr_chain; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Address      string       `json:"address" gorm:"column:address; uniqueIndex:uk_addr; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Private      string       `json:"private" gorm:"column:private; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	AddrStatus   AddrStatus   `json:"addr_status" gorm:"column:addr_status; type:smallint(6) NOT NULL DEFAULT '0' COMMENT '';"`
	Remark       string       `json:"remark" gorm:"column:remark; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	CompressType CompressType `json:"compress_type" gorm:"column:compress_type; type:smallint(6) NOT NULL DEFAULT '0' COMMENT '';"`
	CreatedAt    time.Time    `json:"created_at" gorm:"column:created_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '';"`
	UpdatedAt    time.Time    `json:"updated_at" gorm:"column:updated_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '';"`
}

const (
	TableNameAddressInfo = "t_address_info"
)

func (t *TableAddressInfo) TableName() string {
	return TableNameAddressInfo
}

type AddrChain string

const (
	AddrChainEVM  AddrChain = "EVM"
	AddrChainTRON AddrChain = "TRON"
	AddrChainDOGE AddrChain = "DOGE"
	AddrChainCKB  AddrChain = "CKB"
)

type AddrStatus int

const (
	AddrStatusDefault AddrStatus = 0
	AddrStatusDisable AddrStatus = 1
)

type CompressType int

const (
	CompressTypeFalse CompressType = 0
	CompressTypeTrue  CompressType = 1
)

func (c CompressType) Bool() bool {
	return c == CompressTypeTrue
}
