package dao

import (
	"fmt"
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
