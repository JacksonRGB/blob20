package dao

import (
	dbmodel "blob-index/model/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *Dao) CreateMint(mint *dbmodel.Mint) (idx int, err error) {
	err = d.db.Clauses(clause.OnConflict{DoNothing: true}).Create(mint).Error
	if err != nil {
		return 0, err
	}
	return int(mint.ID), nil
}

// GetStorageHeight 获取上次缓存的高度
func (d *Dao) GetStorageHeight(key string) (value int, err error) {
	storage := new(dbmodel.Height)
	err = d.db.Model(storage).Where("`key` = ?", key).First(storage).Error
	if err == gorm.ErrRecordNotFound {
		return 19526000, nil
	}
	return storage.IntValue, err
}

// SetStorageHeight 设置上次缓存的高度
func (d *Dao) SetStorageHeight(key string, intValue int) (err error) {
	ret := d.db.Model(&dbmodel.Height{}).Where("`key` = ?", key).Update("int_value", intValue)
	if ret.Error != nil {
		return
	}
	if ret.RowsAffected == 0 {
		err = d.db.Create(&dbmodel.Height{
			Key:      key,
			IntValue: intValue,
		}).Error
	}
	return
}

func (d *Dao) GetMintsByOwner(owner string, page, pageSize int) (mints []*dbmodel.Mint, err error) {
	if owner == "" {
		err = d.db.Model(&dbmodel.Mint{}).Order("`id` desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&mints).Error
		return
	}

	err = d.db.Model(&dbmodel.Mint{}).Where("owner = ?", owner).Order("`id` desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&mints).Error
	return
}

func (d *Dao) GetMintCountByOwner(owner string) (count int, err error) {
	var ct int64
	err = d.db.Model(&dbmodel.Mint{}).Where("owner = ?", owner).Count(&ct).Error
	return int(ct), err
}

func (d *Dao) GetMintCount() (count int, err error) {
	var ct int64
	err = d.db.Model(&dbmodel.Mint{}).Count(&ct).Error
	return int(ct), err
}

type UserAmount struct {
	Owner string `json:"owner"`
	Ct    int    `json:"count"`
}

func (d *Dao) GetTopUsers() (uas []*UserAmount, err error) {
	err = d.db.Model(&dbmodel.Mint{}).Select("owner, count(*) as ct").
		Group("owner").Order("ct desc").Limit(50).Find(&uas).Error
	return
}
