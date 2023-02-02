package vendor

import (
	"bpm/core/database"
	"errors"
)

type vendorService struct {
}

func NewVendorService() *vendorService {
	return &vendorService{}
}

func (s *vendorService) GetVendorByID(id int64) (*VendorResponse, error) {
	db := database.InitMySQL()
	query := NewVendorQuery(db)
	vendor, err := query.GetVendorByID(id)
	if err != nil {
		return nil, err
	}
	materials, err := query.GetVendorMaterial(id)
	if err != nil {
		return nil, err
	}
	brands, err := query.GetVendorBrand(id)
	if err != nil {
		return nil, err
	}
	pictures, err := query.GetVendorPicture(id)
	if err != nil {
		return nil, err
	}
	qrcodes, err := query.GetVendorQrcode(id)
	if err != nil {
		return nil, err
	}
	vendor.Material = *materials
	vendor.Brand = *brands
	vendor.Picture = *pictures
	vendor.Qrcode = *qrcodes
	return vendor, err
}

func (s *vendorService) NewVendor(info VendorNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewVendorRepository(tx)
	exist, err := repo.CheckVendorNameExist(info.Name, 0)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "商家名称重复"
		return errors.New(msg)
	}
	vendorID, err := repo.CreateVendor(info)
	if err != nil {
		return err
	}
	if len(info.Material) > 0 {
		for _, material := range info.Material {
			materialExist, err := repo.CheckMaterialExist(material)
			if err != nil {
				return err
			}
			if materialExist != 1 {
				msg := "材料不存在"
				return errors.New(msg)
			}
			err = repo.CreateVendorMaterial(vendorID, material, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Brand) > 0 {
		for _, brand := range info.Brand {
			brandExist, err := repo.CheckBrandExist(brand)
			if err != nil {
				return err
			}
			if brandExist != 1 {
				msg := "品牌不存在"
				return errors.New(msg)
			}
			err = repo.CreateVendorBrand(vendorID, brand, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Picture) > 0 {
		for _, picture := range info.Picture {
			err = repo.CreateVendorPicture(vendorID, picture, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Qrcode) > 0 {
		for _, qrcode := range info.Qrcode {
			err = repo.CreateVendorQrcode(vendorID, qrcode, info.User)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *vendorService) GetVendorList(filter VendorFilter) (int, *[]VendorResponse, error) {
	db := database.InitMySQL()
	query := NewVendorQuery(db)
	count, err := query.GetVendorCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetVendorList(filter)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range *list {
		materials, err := query.GetVendorMaterial(v.ID)
		if err != nil {
			return 0, nil, err
		}
		brands, err := query.GetVendorBrand(v.ID)
		if err != nil {
			return 0, nil, err
		}
		pictures, err := query.GetVendorPicture(v.ID)
		if err != nil {
			return 0, nil, err
		}
		qrcodes, err := query.GetVendorQrcode(v.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[k].Material = *materials
		(*list)[k].Brand = *brands
		(*list)[k].Picture = *pictures
		(*list)[k].Qrcode = *qrcodes
	}
	return count, list, err
}

func (s *vendorService) UpdateVendor(vendorID int64, info VendorNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewVendorRepository(tx)
	_, err = repo.GetVendorByID(vendorID)
	if err != nil {
		msg := "商家不存在"
		return errors.New(msg)
	}
	exist, err := repo.CheckVendorNameExist(info.Name, vendorID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "商家名称重复"
		return errors.New(msg)
	}
	err = repo.DeleteVendorMaterial(vendorID, info.User)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorBrand(vendorID, info.User)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorPicture(vendorID, info.User)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorQrcode(vendorID, info.User)
	if err != nil {
		return err
	}
	err = repo.UpdateVendor(vendorID, info)
	if err != nil {
		return err
	}

	if len(info.Material) > 0 {
		for _, material := range info.Material {
			materialExist, err := repo.CheckMaterialExist(material)
			if err != nil {
				return err
			}
			if materialExist != 1 {
				msg := "材料不存在"
				return errors.New(msg)
			}
			err = repo.CreateVendorMaterial(vendorID, material, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Brand) > 0 {
		for _, brand := range info.Brand {
			brandExist, err := repo.CheckBrandExist(brand)
			if err != nil {
				return err
			}
			if brandExist != 1 {
				msg := "品牌不存在"
				return errors.New(msg)
			}
			err = repo.CreateVendorBrand(vendorID, brand, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Picture) > 0 {
		for _, picture := range info.Picture {
			err = repo.CreateVendorPicture(vendorID, picture, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Qrcode) > 0 {
		for _, qrcode := range info.Qrcode {
			err = repo.CreateVendorQrcode(vendorID, qrcode, info.User)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *vendorService) DeleteVendor(vendorID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewVendorRepository(tx)
	_, err = repo.GetVendorByID(vendorID)
	if err != nil {
		msg := "商家不存在"
		return errors.New(msg)
	}
	err = repo.DeleteVendor(vendorID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorMaterial(vendorID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorBrand(vendorID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorPicture(vendorID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorQrcode(vendorID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
