package common

import (
	"bpm/core/database"
	"errors"
)

type commonService struct {
}

func NewCommonService() *commonService {
	return &commonService{}
}

func (s *commonService) GetBrandByID(id int64) (*BrandResponse, error) {
	db := database.InitMySQL()
	query := NewCommonQuery(db)
	brand, err := query.GetBrandByID(id)
	return brand, err
}

func (s *commonService) NewBrand(info BrandNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	exist, err := repo.CheckBrandNameExist(info.Name, 0)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "品牌名称重复"
		return errors.New(msg)
	}
	err = repo.CreateBrand(info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *commonService) GetBrandList(filter BrandFilter) (int, *[]BrandResponse, error) {
	db := database.InitMySQL()
	query := NewCommonQuery(db)
	count, err := query.GetBrandCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetBrandList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *commonService) UpdateBrand(brandID int64, info BrandNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	_, err = repo.GetBrandByID(brandID)
	if err != nil {
		msg := "品牌不存在"
		return errors.New(msg)
	}
	exist, err := repo.CheckBrandNameExist(info.Name, brandID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "品牌名称重复"
		return errors.New(msg)
	}
	err = repo.UpdateBrand(brandID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *commonService) DeleteBrand(brandID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	_, err = repo.GetBrandByID(brandID)
	if err != nil {
		msg := "品牌不存在"
		return errors.New(msg)
	}
	exist, err := repo.CheckBrandActive(brandID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "品牌正在使用"
		return errors.New(msg)
	}
	err = repo.DeleteBrand(brandID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *commonService) GetMaterialByID(id int64) (*MaterialResponse, error) {
	db := database.InitMySQL()
	query := NewCommonQuery(db)
	brand, err := query.GetMaterialByID(id)
	return brand, err
}

func (s *commonService) NewMaterial(info MaterialNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	exist, err := repo.CheckMaterialNameExist(info.Name, 0)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "材料名称重复"
		return errors.New(msg)
	}
	err = repo.CreateMaterial(info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *commonService) GetMaterialList(filter MaterialFilter) (int, *[]MaterialResponse, error) {
	db := database.InitMySQL()
	query := NewCommonQuery(db)
	count, err := query.GetMaterialCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetMaterialList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *commonService) UpdateMaterial(brandID int64, info MaterialNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	_, err = repo.GetMaterialByID(brandID)
	if err != nil {
		msg := "材料不存在"
		return errors.New(msg)
	}
	exist, err := repo.CheckMaterialNameExist(info.Name, brandID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "材料名称重复"
		return errors.New(msg)
	}
	err = repo.UpdateMaterial(brandID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *commonService) DeleteMaterial(brandID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCommonRepository(tx)
	_, err = repo.GetMaterialByID(brandID)
	if err != nil {
		msg := "品牌不存在"
		return errors.New(msg)
	}
	exist, err := repo.CheckMaterialActive(brandID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "品牌正在使用"
		return errors.New(msg)
	}
	err = repo.DeleteMaterial(brandID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
