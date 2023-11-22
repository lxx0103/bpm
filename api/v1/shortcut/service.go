package shortcut

import (
	"bpm/core/database"
	"errors"
)

type shortcutService struct {
}

func NewShortcutService() *shortcutService {
	return &shortcutService{}
}

func (s *shortcutService) GetShortcutByID(id int64, organizationID int64) (*ShortcutResponse, error) {
	db := database.InitMySQL()
	query := NewShortcutQuery(db)
	shortcut, err := query.GetShortcutByID(id, organizationID)
	return shortcut, err
}

func (s *shortcutService) NewShortcut(info ShortcutNew, organizationID int64) error {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "组织ID错误"
		return errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewShortcutRepository(tx)
	currentShortcutType := info.ShortcutType
	shortcutTypeName := ""
	for {
		shortcutType, err := repo.GetShortcutTypeByID(currentShortcutType)
		if err != nil {
			msg := "模版类别不存在"
			return errors.New(msg)
		}
		if shortcutTypeName == "" {
			shortcutTypeName = shortcutType.Name
		} else {
			shortcutTypeName = shortcutType.Name + " - " + shortcutTypeName
		}
		currentShortcutType = shortcutType.ParentID
		if shortcutType.ParentID == 0 {
			break
		}
	}
	info.ShortcutTypeName = shortcutTypeName
	_, err = repo.CreateShortcut(info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *shortcutService) GetShortcutList(filter ShortcutFilter, organizationID int64) (int, *[]ShortcutResponse, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewShortcutQuery(db)
	count, err := query.GetShortcutCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetShortcutList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, err
}

func (s *shortcutService) UpdateShortcut(shortcutID int64, info ShortcutUpdate, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewShortcutRepository(tx)
	oldShortcut, err := repo.GetShortcutByID(shortcutID)
	if err != nil {
		msg := "快捷模版记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldShortcut.OrganizationID && organizationID != 0 {
		msg := "快捷模版记录不存在"
		return errors.New(msg)
	}
	currentShortcutType := oldShortcut.ShortcutType
	shortcutTypeName := ""
	for {
		shortcutType, err := repo.GetShortcutTypeByID(currentShortcutType)
		if err != nil {
			msg := "模版类别不存在"
			return errors.New(msg)
		}
		if shortcutTypeName == "" {
			shortcutTypeName = shortcutType.Name
		} else {
			shortcutTypeName = shortcutType.Name + " - " + shortcutTypeName
		}
		currentShortcutType = shortcutType.ParentID
		if shortcutType.ParentID == 0 {
			break
		}
	}
	info.ShortcutTypeName = shortcutTypeName
	err = repo.UpdateShortcut(shortcutID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *shortcutService) DeleteShortcut(shortcutID, organizationID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewShortcutRepository(tx)
	oldShortcut, err := repo.GetShortcutByID(shortcutID)
	if err != nil {
		msg := "快捷模版记录不存在"
		return errors.New(msg)
	}
	if organizationID != oldShortcut.OrganizationID && organizationID != 0 {
		msg := "快捷模版记录不存在"
		return errors.New(msg)
	}
	err = repo.DeleteShortcut(shortcutID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *shortcutService) GetShortcutTypeList(filter ShortcutTypeFilter, organizationID int64) (*[]ShortcutTypeResponse, error) {
	if organizationID == 0 && filter.OrganizationID == 0 {
		msg := "必须指定组织"
		return nil, errors.New(msg)
	}
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewShortcutQuery(db)
	list, err := query.GetShortcutTypeList(filter)
	if err != nil {
		return nil, err
	}
	return list, err
}

func (s *shortcutService) NewShortcutType(info ShortcutTypeNew, organizationID int64) error {
	if organizationID == 0 && info.OrganizationID == 0 {
		msg := "必须指定组织"
		return errors.New(msg)
	}
	if organizationID != 0 {
		info.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewShortcutRepository(tx)
	if info.ParentID != 0 {
		_, err := repo.GetShortcutTypeByID(info.ParentID)
		if err != nil {
			msg := "父级类别不存在"
			return errors.New(msg)
		}
	}
	err = repo.CreateShortcutType(info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *shortcutService) GetShortcutTypeByID(id int64, organizationID int64) (*ShortcutTypeResponse, error) {
	db := database.InitMySQL()
	query := NewShortcutQuery(db)
	shortcutType, err := query.GetShortcutTypeByID(id, organizationID)
	return shortcutType, err
}

func (s *shortcutService) UpdateShortcutType(shortcutTypeID int64, info ShortcutTypeUpdate, organizationID int64) error {
	if info.ParentID == shortcutTypeID {
		msg := "不能更新父级为自己"
		return errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewShortcutRepository(tx)
	oldShortcutType, err := repo.GetShortcutTypeByID(shortcutTypeID)
	if err != nil {
		msg := "快捷模版类别不存在"
		return errors.New(msg)
	}
	if organizationID != oldShortcutType.OrganizationID && organizationID != 0 {
		msg := "快捷模版类别不存在"
		return errors.New(msg)
	}
	err = repo.UpdateShortcutType(shortcutTypeID, info)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *shortcutService) DeleteShortcutType(shortcutTypeID, organizationID int64, byUser string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewShortcutRepository(tx)
	oldShortcut, err := repo.GetShortcutTypeByID(shortcutTypeID)
	if err != nil {
		msg := "快捷模版类别不存在"
		return errors.New(msg)
	}
	if organizationID != oldShortcut.OrganizationID && organizationID != 0 {
		msg := "快捷模版类别不存在"
		return errors.New(msg)
	}
	count, err := repo.GetChildCount(shortcutTypeID)
	if err != nil {
		msg := "获取子级类别失败"
		return errors.New(msg)
	}
	if count > 0 {
		msg := "不能删除有子级类别的类别"
		return errors.New(msg)
	}
	err = repo.DeleteShortcutType(shortcutTypeID, byUser)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
