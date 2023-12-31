package dao

import (
	"fmt"
	"gorm.io/gorm/clause"
	"remote-sign-svr/tables"
)

func (d *DbDao) GetAddressInfo(addr string) (info tables.TableAddressInfo, err error) {
	err = d.db.Where("address=?", addr).Find(&info).Error
	return
}

func (d *DbDao) GetAddressListGroupByAddrChain() (list []tables.TableAddressInfo, err error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE id IN(SELECT MAX(id) FROM %s GROUP BY addr_chain)",
		tables.TableNameAddressInfo, tables.TableNameAddressInfo)
	err = d.db.Raw(sql).Find(&list).Error
	return
}

func (d *DbDao) CreateAddressInfo(addrInfo tables.TableAddressInfo) error {
	return d.db.Clauses(clause.Insert{
		Modifier: "IGNORE",
	}).Create(&addrInfo).Error
}

func (d *DbDao) UpdateAddressStatus(addr string, addrStatus tables.AddrStatus) error {
	return d.db.Model(tables.TableAddressInfo{}).
		Where("address=?", addr).
		Updates(map[string]interface{}{
			"addr_status": addrStatus,
		}).Error
}
