package dao

import "remote-sign-svr/tables"

func (d *DbDao) GetAddressInfo(addrChain tables.AddrChain, addr string) (info tables.TableAddressInfo, err error) {
	err = d.db.Where("addr_chain=? AND address=?",
		addrChain, addr).Find(&info).Error
	return
}
