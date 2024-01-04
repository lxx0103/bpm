package costControl

import (
	"bpm/api/v1/auth"
	"bpm/api/v1/position"
	"bpm/api/v1/project"
	"bpm/core/database"
	"errors"
	"fmt"
	"strconv"
)

type costControlService struct {
}

func NewCostControlService() *costControlService {
	return &costControlService{}
}

func (s *costControlService) GetBudgetList(filter ReqBudgetFilter) (int, *[]RespBudget, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	count, err := query.GetBudgetCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetBudgetList(filter)
	if err != nil {
		return 0, nil, err
	}
	for key, budget := range *list {
		pictures, err := query.GetBudgetPictureList(budget.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[key].Picture = *pictures
	}
	return count, list, nil
}

func (s *costControlService) NewBudget(info ReqBudgetNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	projectRepo := project.NewProjectRepository(tx)
	_, err = projectRepo.GetProjectByID(info.ProjectID, info.OrganizationID)
	if err != nil {
		msg := "项目不存在"
		return errors.New(msg)
	}
	BudgetID, err := repo.CreateBudget(info)
	if err != nil {
		return err
	}
	for _, picture := range info.Picture {
		var pictureInfo ReqBudgetPictureNew
		pictureInfo.BudgetID = BudgetID
		pictureInfo.Picture = picture
		pictureInfo.User = info.User
		err = repo.CreateBudgetPicture(pictureInfo)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *costControlService) UpdateBudget(info ReqBudgetUpdate, id, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldBudget, err := repo.GetBudgetByID(id)
	if err != nil {
		return err
	}
	if oldBudget.OrganizationID != organizationID && organizationID != 0 {
		msg := "预算记录不存在或无权限"
		return errors.New(msg)
	}
	info.Used = oldBudget.Used
	info.Balance = oldBudget.Balance - (oldBudget.Budget - info.Budget)
	err = repo.DeleteBudgetPicture(id)
	if err != nil {
		return err
	}
	err = repo.UpdateBudget(info, id)
	if err != nil {
		return err
	}
	for _, picture := range info.Picture {
		var pictureInfo ReqBudgetPictureNew
		pictureInfo.BudgetID = id
		pictureInfo.Picture = picture
		pictureInfo.User = info.User
		err = repo.CreateBudgetPicture(pictureInfo)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *costControlService) GetBudgetByID(id, organizationID int64) (*RespBudget, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	budget, err := query.GetBudgetByID(id)
	if err != nil {
		return nil, err
	}
	if budget.OrganizationID != organizationID && organizationID != 0 {
		msg := "预算记录不存在或无权限"
		return nil, errors.New(msg)
	}
	pictures, err := query.GetBudgetPictureList(id)
	if err != nil {
		return nil, err
	}
	budget.Picture = *pictures
	return budget, nil
}

func (s *costControlService) DeleteBudget(id, organizationID int64, user string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldBudget, err := repo.GetBudgetByID(id)
	if err != nil {
		return err
	}
	if oldBudget.OrganizationID != organizationID && organizationID != 0 {
		msg := "预算记录不存在或无权限"
		return errors.New(msg)
	}
	err = repo.DeleteBudget(id, user)
	if err != nil {
		return err
	}
	err = repo.DeleteBudgetPicture(id)
	if err != nil {
		msg := "删除预算图片失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *costControlService) NewPaymentRequest(info ReqPaymentRequestNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	auditInfo, err := repo.GetPaymentRequestTypeAudit(info.OrganizationID, info.PaymentRequestType)
	if err != nil {
		msg := "检查审核设置失败"
		return errors.New(msg)
	}
	if len(*auditInfo) == 0 {
		msg := "必须先设置审核人员才能新建请款"
		return errors.New(msg)
	}
	if info.ProjectID != 0 {
		projectRepo := project.NewProjectRepository(tx)
		_, err := projectRepo.GetProjectByID(info.ProjectID, info.OrganizationID)
		if err != nil {
			msg := "项目不存在"
			return errors.New(msg)
		}
		if info.BudgetID != 0 {
			budget, err := repo.GetBudgetByID(info.BudgetID)
			if err != nil {
				msg := "预算记录不存在"
				return errors.New(msg)
			}
			if budget.ProjectID != info.ProjectID {
				msg := "预算记录不存在或无权限"
				return errors.New(msg)
			}
			if budget.BudgetType != info.PaymentRequestType {
				msg := "预算类型与请款类型不一致"
				return errors.New(msg)
			}
		}
	}
	id, err := repo.CreatePaymentRequest(info)
	if err != nil {
		msg := "创建请款记录失败"
		return errors.New(msg)
	}
	for _, picture := range info.Picture {
		var pictureInfo ReqPaymentRequestPictureNew
		pictureInfo.PaymentRequestID = id
		pictureInfo.Picture = picture
		pictureInfo.User = info.User
		err = repo.CreatePaymentRequestPicture(pictureInfo)
		if err != nil {
			msg := "创建请款记录图片失败"
			return errors.New(msg)
		}
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = id
	history.OrganizationID = info.OrganizationID
	history.User = info.User
	history.Action = "新建"
	history.Remark = "请款已新建，当前状态为待审核"
	err = repo.CreatePaymentRequestHistory(history)
	if err != nil {
		msg := "创建请款历史失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil

}

func (s *costControlService) UpdatePaymentRequest(info ReqPaymentRequestUpdate, id, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	auditInfo, err := repo.GetPaymentRequestTypeAudit(organizationID, info.PaymentRequestType)
	if err != nil {
		msg := "检查审核设置失败"
		return errors.New(msg)
	}
	if len(*auditInfo) == 0 {
		msg := "必须先设置审核人员才能新建请款"
		return errors.New(msg)
	}
	oldPaymentRequest, err := repo.GetPaymentRequestByID(id)
	if err != nil {
		msg := "获取请款记录失败"
		return errors.New(msg)
	}
	if oldPaymentRequest.OrganizationID != organizationID && organizationID != 0 {
		msg := "请款记录不存在或无权限"
		return errors.New(msg)
	}
	if oldPaymentRequest.UserID != info.UserID {
		msg := "仅能修改自己创建的请款记录"
		return errors.New(msg)
	}
	if oldPaymentRequest.Status != 1 && oldPaymentRequest.Status != 3 {
		msg := "请款记录状态错误"
		return errors.New(msg)
	}
	if info.ProjectID != 0 {
		projectRepo := project.NewProjectRepository(tx)
		_, err := projectRepo.GetProjectByID(info.ProjectID, organizationID)
		if err != nil {
			msg := "项目不存在"
			return errors.New(msg)
		}
		if info.BudgetID != 0 {
			budget, err := repo.GetBudgetByID(info.BudgetID)
			if err != nil {
				msg := "预算记录不存在"
				return errors.New(msg)
			}
			if budget.ProjectID != info.ProjectID {
				msg := "预算记录不存在或无权限"
				return errors.New(msg)
			}
			if budget.BudgetType != info.PaymentRequestType {
				msg := "预算类型与请款类型不一致"
				return errors.New(msg)
			}
		}
	}
	info.Status = 1
	err = repo.DeletePaymentRequestPicture(id)
	if err != nil {
		return err
	}
	err = repo.UpdatePaymentRequest(info, id)
	if err != nil {
		return err
	}
	for _, picture := range info.Picture {
		var pictureInfo ReqPaymentRequestPictureNew
		pictureInfo.PaymentRequestID = id
		pictureInfo.Picture = picture
		pictureInfo.User = info.User
		repo.CreatePaymentRequestPicture(pictureInfo)
		if err != nil {
			return err
		}
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = id
	history.OrganizationID = oldPaymentRequest.OrganizationID
	history.User = info.User
	history.Action = "更新"
	history.Remark = "请款已更新，当前状态为待审核"
	err = repo.CreatePaymentRequestHistory(history)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s *costControlService) GetPaymentRequestList(filter ReqPaymentRequestFilter) (int, *[]RespPaymentRequest, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	count, err := query.GetPaymentRequestCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPaymentRequestList(filter)
	if err != nil {
		return 0, nil, err
	}
	for key, budget := range *list {
		pictures, err := query.GetPaymentRequestPictureList(budget.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[key].Picture = *pictures
	}
	return count, list, nil
}

func (s *costControlService) GetPaymentRequestByID(id, organizationID int64) (*RespPaymentRequest, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	budget, err := query.GetPaymentRequestByID(id)
	if err != nil {
		return nil, err
	}
	if budget.OrganizationID != organizationID && organizationID != 0 {
		msg := "预算记录不存在或无权限"
		return nil, errors.New(msg)
	}
	pictures, err := query.GetPaymentRequestPictureList(id)
	if err != nil {
		return nil, err
	}
	budget.Picture = *pictures
	return budget, nil
}

func (s *costControlService) DeletePaymentRequest(id, organizationID int64, user string, userID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldPaymentRequest, err := repo.GetPaymentRequestByID(id)
	if err != nil {
		msg := "获取请款记录失败"
		return errors.New(msg)
	}
	if oldPaymentRequest.OrganizationID != organizationID && organizationID != 0 {
		msg := "请款记录不存在或无权限"
		return errors.New(msg)
	}
	if oldPaymentRequest.UserID != userID {
		msg := "只能删除自己创建的请款"
		return errors.New(msg)
	}
	err = repo.DeletePaymentRequest(id, user)
	if err != nil {
		msg := "删除请款失败"
		return errors.New(msg)
	}
	err = repo.DeletePaymentRequestPicture(id)
	if err != nil {
		msg := "删除请款图片失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *costControlService) UpdatePaymentRequestType(info ReqPaymentRequestTypeUpdate, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	positionRepo := position.NewPositionRepository(tx)
	userRepo := auth.NewAuthRepository(tx)
	err = repo.DeletePaymentRequestTypeAudit(info.ReqPaymentRequestType, organizationID, info.User)
	if err != nil {
		msg := "更新审核设置失败"
		return errors.New(msg)
	}
	for _, audit := range info.AuditInfo {
		for _, auditTo := range audit.AuditTo {
			var auditInfo ReqPaymentRequestTypeAudit
			auditInfo.PaymentRequestType = info.ReqPaymentRequestType
			auditInfo.OrganizationID = organizationID
			auditInfo.AuditLevel = audit.AuditLevel
			auditInfo.AuditType = audit.AuditType
			auditInfo.AuditTo = auditTo
			if auditInfo.AuditType == 1 {
				_, err := positionRepo.GetPositionByID(auditTo, organizationID)
				if err != nil {
					msg := "审核层次" + strconv.Itoa(int(auditInfo.AuditLevel)) + "职位不存在"
					return errors.New(msg)
				}
			} else {
				userInfo, err := userRepo.GetUserByID(auditTo)
				if err != nil {
					msg := "审核层次" + strconv.Itoa(int(auditInfo.AuditLevel)) + "用户不存在"
					return errors.New(msg)
				}
				if userInfo.OrganizationID != organizationID {
					msg := "审核层次" + strconv.Itoa(int(auditInfo.AuditLevel)) + "用户不存在"
					return errors.New(msg)
				}
			}
			auditInfo.User = info.User
			err = repo.CreatePaymentRequestTypeAudit(auditInfo)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil

}

func (s *costControlService) GetPaymentRequestTypeList(filter ReqPaymentRequestTypeFilter) (*[]RespPaymentRequestType, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	var type1 RespPaymentRequestType
	type1.PaymentRequestType = 1
	type1.PaymentRequestTypeName = "采购类"
	res1, err := query.GetPaymentRequestTypeList(filter.OrganizationID, 1)
	if err != nil {
		msg := "获取审核设置失败"
		fmt.Println(err)
		return nil, errors.New(msg)
	}
	type1.Audit = *res1
	var type2 RespPaymentRequestType
	type2.PaymentRequestType = 2
	type2.PaymentRequestTypeName = "采购类"
	res2, err := query.GetPaymentRequestTypeList(filter.OrganizationID, 2)
	if err != nil {
		msg := "获取审核设置失败"
		fmt.Println(err)
		return nil, errors.New(msg)
	}
	type2.Audit = *res2
	res := &[]RespPaymentRequestType{type1, type2}
	return res, err
}
