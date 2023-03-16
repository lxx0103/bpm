package vendors

import (
	"bpm/core/database"
	"errors"
)

type vendorsService struct {
}

func NewVendorsService() *vendorsService {
	return &vendorsService{}
}

func (s *vendorsService) GetVendorsByID(id int64) (*VendorsResponse, error) {
	db := database.InitMySQL()
	query := NewVendorsQuery(db)
	vendors, err := query.GetVendorsByID(id)
	if err != nil {
		return nil, err
	}
	materials, err := query.GetVendorsMaterial(id)
	if err != nil {
		return nil, err
	}
	brands, err := query.GetVendorsBrand(id)
	if err != nil {
		return nil, err
	}
	pictures, err := query.GetVendorsPicture(id)
	if err != nil {
		return nil, err
	}
	qrcodes, err := query.GetVendorsQrcode(id)
	if err != nil {
		return nil, err
	}
	vendors.Material = *materials
	vendors.Brand = *brands
	vendors.Picture = *pictures
	vendors.Qrcode = *qrcodes
	return vendors, err
}

func (s *vendorsService) NewVendors(info VendorsNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewVendorsRepository(tx)
	exist, err := repo.CheckVendorsNameExist(info.Name, 0)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "商家名称重复"
		return errors.New(msg)
	}
	vendorsID, err := repo.CreateVendors(info)
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
			err = repo.CreateVendorsMaterial(vendorsID, material, info.User)
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
			err = repo.CreateVendorsBrand(vendorsID, brand, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Picture) > 0 {
		for _, picture := range info.Picture {
			err = repo.CreateVendorsPicture(vendorsID, picture, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Qrcode) > 0 {
		for _, qrcode := range info.Qrcode {
			err = repo.CreateVendorsQrcode(vendorsID, qrcode, info.User)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *vendorsService) GetVendorsList(filter VendorsFilter) (int, *[]VendorsResponse, error) {
	db := database.InitMySQL()
	query := NewVendorsQuery(db)
	count, err := query.GetVendorsCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetVendorsList(filter)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range *list {
		materials, err := query.GetVendorsMaterial(v.ID)
		if err != nil {
			return 0, nil, err
		}
		brands, err := query.GetVendorsBrand(v.ID)
		if err != nil {
			return 0, nil, err
		}
		pictures, err := query.GetVendorsPicture(v.ID)
		if err != nil {
			return 0, nil, err
		}
		qrcodes, err := query.GetVendorsQrcode(v.ID)
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

func (s *vendorsService) UpdateVendors(vendorsID int64, info VendorsNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewVendorsRepository(tx)
	_, err = repo.GetVendorsByID(vendorsID)
	if err != nil {
		msg := "商家不存在"
		return errors.New(msg)
	}
	exist, err := repo.CheckVendorsNameExist(info.Name, vendorsID)
	if err != nil {
		return err
	}
	if exist != 0 {
		msg := "商家名称重复"
		return errors.New(msg)
	}
	err = repo.DeleteVendorsMaterial(vendorsID, info.User)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsBrand(vendorsID, info.User)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsPicture(vendorsID, info.User)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsQrcode(vendorsID, info.User)
	if err != nil {
		return err
	}
	err = repo.UpdateVendors(vendorsID, info)
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
			err = repo.CreateVendorsMaterial(vendorsID, material, info.User)
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
			err = repo.CreateVendorsBrand(vendorsID, brand, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Picture) > 0 {
		for _, picture := range info.Picture {
			err = repo.CreateVendorsPicture(vendorsID, picture, info.User)
			if err != nil {
				return err
			}
		}
	}
	if len(info.Qrcode) > 0 {
		for _, qrcode := range info.Qrcode {
			err = repo.CreateVendorsQrcode(vendorsID, qrcode, info.User)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *vendorsService) DeleteVendors(vendorsID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewVendorsRepository(tx)
	_, err = repo.GetVendorsByID(vendorsID)
	if err != nil {
		msg := "商家不存在"
		return errors.New(msg)
	}
	err = repo.DeleteVendors(vendorsID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsMaterial(vendorsID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsBrand(vendorsID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsPicture(vendorsID, byUser)
	if err != nil {
		return err
	}
	err = repo.DeleteVendorsQrcode(vendorsID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
