package example

import (
	"bpm/api/v1/common"
	"bpm/api/v1/vendors"
	"bpm/core/database"
	"errors"
)

type exampleService struct {
}

func NewExampleService() *exampleService {
	return &exampleService{}
}

func (s *exampleService) GetExampleByID(id int64, organizationID int64) (*Example, error) {
	db := database.InitMySQL()
	query := NewExampleQuery(db)
	example, err := query.GetExampleByID(id, organizationID)
	return example, err
}

func (s *exampleService) NewExample(info ExampleNew, organizationID int64) (*Example, error) {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	// exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, 0)
	// if err != nil {
	// 	return nil, err
	// }
	// if exist != 0 {
	// 	msg := "案例名称重复"
	// 	return nil, errors.New(msg)
	// }
	exampleID, err := repo.CreateExample(info)
	if err != nil {
		return nil, err
	}
	example, err := repo.GetExampleByID(exampleID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return example, err
}

func (s *exampleService) GetExampleList(filter ExampleFilter, organizationID int64) (int, *[]ExampleListResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewExampleQuery(db)
	count, err := query.GetExampleCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetExampleList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *exampleService) UpdateExample(exampleID int64, info ExampleNew, organizationID int64) (*Example, error) {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	// exist, err := repo.CheckNameExist(info.Name, info.OrganizationID, exampleID)
	// if err != nil {
	// 	return nil, err
	// }
	// if exist != 0 {
	// 	msg := "案例名称重复"
	// 	return nil, errors.New(msg)
	// }
	_, err = repo.UpdateExample(exampleID, info)
	if err != nil {
		return nil, err
	}
	example, err := repo.GetExampleByID(exampleID, info.OrganizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return example, err
}

func (s *exampleService) GetExampleMaterialList(exampleID int64) (*[]ExampleMaterialResponse, error) {
	db := database.InitMySQL()
	query := NewExampleQuery(db)
	list, err := query.GetExampleMaterialList(exampleID)
	return list, err
}

func (s *exampleService) GetExampleMaterialByID(id int64, materialID int64) (*ExampleMaterialResponse, error) {
	db := database.InitMySQL()
	query := NewExampleQuery(db)
	example, err := query.GetExampleMaterialByID(id, materialID)
	return example, err
}

func (s *exampleService) NewExampleMaterial(info ExampleMaterialNew, exampleID, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	vendorsRepo := vendors.NewVendorsRepository(tx)
	commonRepo := common.NewCommonRepository(tx)
	_, err = repo.GetExampleByID(exampleID, organizationID)
	if err != nil {
		msg := "案例不存在"
		return errors.New(msg)
	}
	_, err = commonRepo.GetMaterialByID(info.MaterialID)
	if err != nil {
		msg := "材料不存在"
		return errors.New(msg)
	}
	if info.VendorID != 0 {
		_, err = vendorsRepo.GetVendorsByID(info.VendorID)
		if err != nil {
			msg := "供应商不存在"
			return errors.New(msg)
		}
	}
	if info.BrandID != 0 {
		_, err = commonRepo.GetBrandByID(info.BrandID)
		if err != nil {
			msg := "品牌不存在"
			return errors.New(msg)
		}
	}
	err = repo.CreateExampleMaterial(info, exampleID)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (s *exampleService) UpdateExampleMaterial(info ExampleMaterialNew, exampleID, ID, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	vendorsRepo := vendors.NewVendorsRepository(tx)
	commonRepo := common.NewCommonRepository(tx)
	_, err = repo.GetExampleByID(exampleID, organizationID)
	if err != nil {
		msg := "案例不存在"
		return errors.New(msg)
	}
	_, err = repo.GetExampleMaterialByID(ID)
	if err != nil {
		msg := "案例材料不存在"
		return errors.New(msg)
	}
	_, err = commonRepo.GetMaterialByID(info.MaterialID)
	if err != nil {
		msg := "材料不存在"
		return errors.New(msg)
	}
	if info.VendorID != 0 {
		_, err = vendorsRepo.GetVendorsByID(info.VendorID)
		if err != nil {
			msg := "供应商不存在"
			return errors.New(msg)
		}
	}
	if info.BrandID != 0 {
		_, err = commonRepo.GetBrandByID(info.BrandID)
		if err != nil {
			msg := "品牌不存在"
			return errors.New(msg)
		}
	}
	err = repo.UpdateExampleMaterial(info, ID)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (s *exampleService) DeleteExampleMaterial(exampleID, ID, organizationID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewExampleRepository(tx)
	_, err = repo.GetExampleByID(exampleID, organizationID)
	if err != nil {
		msg := "案例不存在"
		return errors.New(msg)
	}
	_, err = repo.GetExampleMaterialByID(ID)
	if err != nil {
		msg := "案例材料不存在"
		return errors.New(msg)
	}
	err = repo.DeleteExampleMaterial(ID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}
