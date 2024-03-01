package costControl

import (
	"bpm/api/v1/auth"
	"bpm/api/v1/position"
	"bpm/api/v1/project"
	"bpm/core/database"
	"bpm/core/queue"
	"encoding/json"
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
	for _, audit := range *auditInfo {
		var auditNew ReqPaymentRequestAuditNew
		auditNew.PaymentRequestID = id
		auditNew.AuditLevel = audit.AuditLevel
		auditNew.AuditType = audit.AuditType
		auditNew.AuditTo = audit.AuditTo
		auditNew.User = info.User
		err = repo.CreatePaymentRequestAudit(auditNew)
		if err != nil {
			msg := "创建请款记录审核失败"
			return errors.New(msg)
		}
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = id
	history.OrganizationID = info.OrganizationID
	history.User = info.User
	history.Action = "新建"
	history.Content = ""
	history.Remark = "请款已新建，当前状态为待审核"
	historyID, err := repo.CreatePaymentRequestHistory(history)
	if err != nil {
		msg := "创建请款历史失败"
		return errors.New(msg)
	}
	for _, picture := range info.Picture {
		var pictureInfo ReqPaymentRequestHistoryPictureNew
		pictureInfo.PaymentRequestHistoryID = historyID
		pictureInfo.Picture = picture
		pictureInfo.User = info.User
		err = repo.CreatePaymentRequestHistoryPicture(pictureInfo)
		if err != nil {
			msg := "创建请款记录历史图片失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	type NewPaymentRequestCreated struct {
		PaymentRequestID int64 `json:"payment_request_id"`
	}
	var newEvent NewPaymentRequestCreated
	newEvent.PaymentRequestID = id
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newEvent)
	err = rabbit.Publish("NewPaymentRequestCreated", msg)
	if err != nil {
		msg := "create event NewPaymentRequestCreated error"
		return errors.New(msg)
	}
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
	if oldPaymentRequest.Status == 1 && oldPaymentRequest.AuditLevel != 1 {
		msg := "当前正在审核，请勿修改"
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
		err = repo.CreatePaymentRequestPicture(pictureInfo)
		if err != nil {
			return err
		}
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = id
	history.OrganizationID = oldPaymentRequest.OrganizationID
	history.User = info.User
	history.Action = "更新"
	history.Content = ""
	history.Remark = "请款已更新，当前状态为待审核"
	historyID, err := repo.CreatePaymentRequestHistory(history)
	if err != nil {
		return err
	}
	for _, picture := range info.Picture {
		var pictureInfo ReqPaymentRequestHistoryPictureNew
		pictureInfo.PaymentRequestHistoryID = historyID
		pictureInfo.Picture = picture
		pictureInfo.User = info.User
		err = repo.CreatePaymentRequestHistoryPicture(pictureInfo)
		if err != nil {
			msg := "创建请款记录历史图片失败"
			return errors.New(msg)
		}
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
		audits, err := query.GetPaymentRequestAuditList(budget.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[key].Audit = *audits
	}
	return count, list, nil
}

func (s *costControlService) GetPaymentRequestByID(id, organizationID int64) (*RespPaymentRequest, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	paymentRequest, err := query.GetPaymentRequestByID(id)
	if err != nil {
		return nil, err
	}
	if paymentRequest.OrganizationID != organizationID && organizationID != 0 {
		msg := "预算记录不存在或无权限"
		return nil, errors.New(msg)
	}
	pictures, err := query.GetPaymentRequestPictureList(id)
	if err != nil {
		return nil, err
	}
	audits, err := query.GetPaymentRequestAuditList(id)
	if err != nil {
		return nil, err
	}
	paymentRequest.Picture = *pictures
	paymentRequest.Audit = *audits
	return paymentRequest, nil
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
	if oldPaymentRequest.Status != 1 && oldPaymentRequest.Status != 2 && oldPaymentRequest.Status != 3 {
		msg := "请款记录无法删除，可能已付款"
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

	if oldPaymentRequest.BudgetID != 0 {
		oldBudget, err := repo.GetBudgetByID(oldPaymentRequest.BudgetID)
		if err != nil {
			msg := "获取预算失败"
			return errors.New(msg)
		}
		var budgetUpdate ReqBudgetPaid
		budgetUpdate.Used = oldBudget.Used - oldPaymentRequest.Total
		budgetUpdate.Balance = oldBudget.Balance + oldPaymentRequest.Total
		budgetUpdate.User = user
		err = repo.UpdateBudgitUsed(oldPaymentRequest.BudgetID, budgetUpdate)
		if err != nil {
			msg := "更新预算信息失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	return nil
}

func (s *costControlService) UpdatePaymentRequestType(info ReqPaymentRequestTypeUpdate) error {
	if info.OrganizationID == 0 {
		msg := "必须指定组织"
		return errors.New(msg)
	}
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	positionRepo := position.NewPositionRepository(tx)
	userRepo := auth.NewAuthRepository(tx)
	err = repo.DeletePaymentRequestTypeAudit(info.ReqPaymentRequestType, info.OrganizationID, info.User)
	if err != nil {
		msg := "更新审核设置失败"
		return errors.New(msg)
	}
	for _, audit := range info.AuditInfo {
		for _, auditTo := range audit.AuditTo {
			var auditInfo ReqPaymentRequestTypeAudit
			auditInfo.PaymentRequestType = info.ReqPaymentRequestType
			auditInfo.OrganizationID = info.OrganizationID
			auditInfo.AuditLevel = audit.AuditLevel
			auditInfo.AuditType = audit.AuditType
			auditInfo.AuditTo = auditTo
			if auditInfo.AuditType == 1 {
				_, err := positionRepo.GetPositionByID(auditTo, info.OrganizationID)
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
				if userInfo.OrganizationID != info.OrganizationID {
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
	type2.PaymentRequestTypeName = "工款类"
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

func (s *costControlService) AuditPaymentRequest(paymentRequestID int64, info ReqPaymentRequestAudit) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	paymentRequest, err := repo.GetPaymentRequestByID(paymentRequestID)
	if err != nil {
		msg := "获取请款记录失败"
		return errors.New(msg)
	}
	if paymentRequest.Status != 1 {
		msg := "此请款无法审核"
		return errors.New(msg)
	}
	assignExist, err := repo.CheckAudit(paymentRequestID, info.UserID, info.PositionID, paymentRequest.AuditLevel)
	fmt.Println(info.PositionID, paymentRequest.AuditLevel)
	if err != nil {
		msg := "检查审核设置失败"
		return errors.New(msg)
	}
	if assignExist == 0 {
		msg := "此请款审核未分配给你"
		return errors.New(msg)
	}
	nextLevel := 0
	result := "审核通过"
	remark := "审核已通过"
	status := 2
	if info.Result == 1 {
		nextLevel, err = repo.GetNextLevel(paymentRequestID, paymentRequest.AuditLevel)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				nextLevel = 0
				remark += "，没有下一层审核，当前状态为审核通过"
			} else {
				msg := "获取下一层审核失败"
				return errors.New(msg)
			}
		} else {
			remark += "，当前状态为待审核，下一层审核为第" + fmt.Sprintf("%d", nextLevel) + "层"
			status = 1
		}
	} else {
		result = "审核不通过"
		remark = "审核不通过，当前状态为审核驳回"
		status = 3
		nextLevel = 1
	}
	// remark += "\n " + info.Content
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = paymentRequestID
	history.OrganizationID = paymentRequest.OrganizationID
	history.User = info.User
	history.Action = result
	history.Content = info.Content
	history.Remark = remark
	historyID, err := repo.CreatePaymentRequestHistory(history)
	if err != nil {
		return err
	}
	err = repo.AuditPaymentRequest(paymentRequestID, nextLevel, status, info.User)
	if err != nil {
		msg := "更新请款状态失败"
		return errors.New(msg)
	}
	for _, link := range info.File {
		var paymentRequestHistoryPicture ReqPaymentRequestHistoryPictureNew
		paymentRequestHistoryPicture.PaymentRequestHistoryID = historyID
		paymentRequestHistoryPicture.Picture = link
		paymentRequestHistoryPicture.User = info.User
		err = repo.CreatePaymentRequestHistoryPicture(paymentRequestHistoryPicture)
		if err != nil {
			msg := "创建文件失败"
			return errors.New(msg)
		}
	}
	if nextLevel == 0 && paymentRequest.BudgetID != 0 {
		oldBudget, err := repo.GetBudgetByID(paymentRequest.BudgetID)
		if err != nil {
			msg := "获取预算失败"
			return errors.New(msg)
		}
		var budgetUpdate ReqBudgetPaid
		budgetUpdate.Used = oldBudget.Used + paymentRequest.Total
		budgetUpdate.Balance = oldBudget.Balance - paymentRequest.Total
		budgetUpdate.User = info.User
		err = repo.UpdateBudgitUsed(paymentRequest.BudgetID, budgetUpdate)
		if err != nil {
			msg := "更新预算信息失败"
			return errors.New(msg)
		}
	}

	tx.Commit()
	type NewPaymentRequestAudited struct {
		PaymentRequestID int64 `json:"paymentRequest_id"`
	}
	var newPaymentRequest NewPaymentRequestAudited
	newPaymentRequest.PaymentRequestID = paymentRequestID
	rabbit, _ := queue.GetConn()
	msg, _ := json.Marshal(newPaymentRequest)
	err = rabbit.Publish("NewPaymentRequestAudited", msg)
	if err != nil {
		msg := "create paymentRequest NewPaymentRequestAudited error"
		return errors.New(msg)
	}
	return nil
}

func (s *costControlService) GetPaymentRequestHistoryList(filter ReqPaymentRequestHistoryFilter) (*[]RespPaymentRequestHistory, error) {

	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	paymentRequest, err := query.GetPaymentRequestByID(filter.PaymentRequestID)
	if err != nil {
		msg := "获取请款记录失败"
		return nil, errors.New(msg)
	}
	if paymentRequest.OrganizationID != filter.OrganizationID && filter.OrganizationID != 0 {
		msg := "请款记录不存在"
		return nil, errors.New(msg)
	}
	list, err := query.GetPaymentRequestHistoryList(filter.PaymentRequestID)
	if err != nil {
		return nil, err
	}
	for key, history := range *list {
		pictures, err := query.GetPaymentRequestHistoryPictureList(history.ID)
		if err != nil {
			return nil, err
		}
		(*list)[key].Picture = *pictures
	}
	return list, nil
}

func (s *costControlService) UpdatePaymentRequestAudit(id int64, info ReqPaymentRequestAuditUpdate, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	positionRepo := position.NewPositionRepository(tx)
	userRepo := auth.NewAuthRepository(tx)
	err = repo.DeletePaymentRequestAudit(id, organizationID, info.User)
	if err != nil {
		msg := "更新请款审核失败"
		return errors.New(msg)
	}
	for _, audit := range info.AuditInfo {
		for _, auditTo := range audit.AuditTo {
			var auditInfo ReqPaymentRequestAuditNew
			auditInfo.PaymentRequestID = id
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
			err = repo.CreatePaymentRequestAudit(auditInfo)
			if err != nil {
				return err
			}
		}
	}
	err = repo.AuditPaymentRequest(id, 1, 1, info.User)
	if err != nil {
		msg := "更新请款审核失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil

}

func (s *costControlService) NewPayment(paymentRequestID int64, info ReqPaymentNew, organizationID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	paymentRequest, err := repo.GetPaymentRequestByID(paymentRequestID)
	if err != nil {
		msg := "获取请款记录失败"
		return errors.New(msg)
	}
	if paymentRequest.OrganizationID != organizationID {
		msg := "请款记录不存在"
		return errors.New(msg)
	}
	if paymentRequest.Status != 2 && paymentRequest.Status != 4 {
		msg := "请款记录状态不正确"
		return errors.New(msg)
	}
	if info.Amount > paymentRequest.Due {
		msg := "此次付款金额大于未付款金额"
		return errors.New(msg)
	}
	info.OrganizationID = paymentRequest.OrganizationID
	info.ProjectID = paymentRequest.ProjectID
	paymentID, err := repo.CreatePayment(paymentRequestID, info)
	if err != nil {
		msg := "生成付款记录失败"
		return errors.New(msg)
	}
	for _, picture := range info.Picture {
		var paymentPicture ReqPaymentPictureNew
		paymentPicture.PaymentID = paymentID
		paymentPicture.Picture = picture
		paymentPicture.User = info.User
		err = repo.CreatePaymentPicture(paymentPicture)
		if err != nil {
			msg := "创建付款记录文件失败"
			return errors.New(msg)
		}
	}
	var paymentRequestUpdate ReqPaymentRequestPaid
	paymentRequestUpdate.Paid = paymentRequest.Paid + info.Amount
	paymentRequestUpdate.Due = paymentRequest.Due - info.Amount
	if paymentRequestUpdate.Due == 0 {
		paymentRequestUpdate.Status = 5
	} else {
		paymentRequestUpdate.Status = 4
	}
	paymentRequestUpdate.User = info.User
	err = repo.UpdatePaymentRequestPaid(paymentRequestID, paymentRequestUpdate)
	if err != nil {
		msg := "更新请款信息失败"
		return errors.New(msg)
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = paymentRequestID
	history.OrganizationID = paymentRequest.OrganizationID
	history.User = info.User
	history.Action = "付款"
	history.Remark = "本次付款金额" + strconv.FormatFloat(info.Amount, 'f', 2, 64) + "元，"
	if paymentRequestUpdate.Due == 0 {
		history.Remark += "已完全付款，当前状态为已付款"
	} else {
		history.Remark += "未完全付款，当前状态为部分付款"
	}
	history.Content = info.Remark
	historyID, err := repo.CreatePaymentRequestHistory(history)
	if err != nil {
		msg := "生成付款记录失败"
		return errors.New(msg)
	}
	for _, link := range info.Picture {
		var paymentRequestHistoryPicture ReqPaymentRequestHistoryPictureNew
		paymentRequestHistoryPicture.PaymentRequestHistoryID = historyID
		paymentRequestHistoryPicture.Picture = link
		paymentRequestHistoryPicture.User = info.User
		err = repo.CreatePaymentRequestHistoryPicture(paymentRequestHistoryPicture)
		if err != nil {
			msg := "创建付款记录文件失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	return nil
}

func (s *costControlService) UpdatePayment(id int64, info ReqPaymentUpdate) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldPayment, err := repo.GetPaymentByID(id)
	fmt.Println(oldPayment.PaymentRequestID)
	if err != nil {
		msg := "获取付款信息失败"
		return errors.New(msg)
	}
	if oldPayment.UserID != info.UserID {
		msg := "只能更新自己的付款"
		return errors.New(msg)
	}
	if oldPayment.OrganizationID != info.OrganizationID {
		msg := "付款信息不存在"
		return errors.New(msg)
	}
	paymentRequest, err := repo.GetPaymentRequestByID(oldPayment.PaymentRequestID)
	if err != nil {
		msg := "获取请款记录失败"
		return errors.New(msg)
	}
	if paymentRequest.OrganizationID != info.OrganizationID {
		msg := "请款记录不存在"
		return errors.New(msg)
	}
	paymentRequest.Due = paymentRequest.Due + oldPayment.Amount
	paymentRequest.Paid = paymentRequest.Paid - oldPayment.Amount
	fmt.Println(paymentRequest.Due, paymentRequest.Paid)
	if info.Amount > paymentRequest.Due {
		msg := "此次付款金额大于未付款金额"
		return errors.New(msg)
	}
	err = repo.DeletePaymentPicture(id)
	if err != nil {
		msg := "删除付款图片失败"
		return errors.New(msg)
	}
	err = repo.UpdatePayment(id, info)
	if err != nil {
		msg := "更新付款记录失败"
		return errors.New(msg)
	}
	for _, picture := range info.Picture {
		var paymentPicture ReqPaymentPictureNew
		paymentPicture.PaymentID = id
		paymentPicture.Picture = picture
		paymentPicture.User = info.User
		err = repo.CreatePaymentPicture(paymentPicture)
		if err != nil {
			msg := "创建付款记录文件失败"
			return errors.New(msg)
		}
	}
	var paymentRequestUpdate ReqPaymentRequestPaid
	paymentRequestUpdate.Paid = paymentRequest.Paid + info.Amount
	paymentRequestUpdate.Due = paymentRequest.Due - info.Amount
	if paymentRequestUpdate.Due == 0 {
		paymentRequestUpdate.Status = 5
	} else {
		paymentRequestUpdate.Status = 4
	}
	paymentRequestUpdate.User = info.User
	err = repo.UpdatePaymentRequestPaid(oldPayment.PaymentRequestID, paymentRequestUpdate)
	if err != nil {
		msg := "更新请款信息失败"
		return errors.New(msg)
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = oldPayment.PaymentRequestID
	history.OrganizationID = paymentRequest.OrganizationID
	history.User = info.User
	history.Action = "更新付款"
	history.Remark = "本次付款金额由" + strconv.FormatFloat(oldPayment.Amount, 'f', 2, 64) + "元更新为" + strconv.FormatFloat(info.Amount, 'f', 2, 64) + "元，"
	if paymentRequestUpdate.Due == 0 {
		history.Remark += "已完全付款，当前状态为已付款"
	} else {
		history.Remark += "未完全付款，当前状态为部分付款"
	}
	history.Content = info.Remark
	historyID, err := repo.CreatePaymentRequestHistory(history)
	if err != nil {
		msg := "生成付款记录失败"
		return errors.New(msg)
	}
	for _, link := range info.Picture {
		var paymentRequestHistoryPicture ReqPaymentRequestHistoryPictureNew
		paymentRequestHistoryPicture.PaymentRequestHistoryID = historyID
		paymentRequestHistoryPicture.Picture = link
		paymentRequestHistoryPicture.User = info.User
		err = repo.CreatePaymentRequestHistoryPicture(paymentRequestHistoryPicture)
		if err != nil {
			msg := "创建付款记录文件失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	return nil
}

func (s *costControlService) GetPaymentList(filter ReqPaymentFilter) (int, *[]RespPayment, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	count, err := query.GetPaymentCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetPaymentList(filter)
	if err != nil {
		return 0, nil, err
	}
	for key, budget := range *list {
		pictures, err := query.GetPaymentPictureList(budget.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[key].Picture = *pictures
	}
	return count, list, nil
}

func (s *costControlService) GetPaymentByID(id, organizationID int64) (*RespPayment, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	payment, err := query.GetPaymentByID(id)
	if err != nil {
		return nil, err
	}
	if payment.OrganizationID != organizationID && organizationID != 0 {
		msg := "付款不存在或无权限"
		return nil, errors.New(msg)
	}
	pictures, err := query.GetPaymentPictureList(id)
	if err != nil {
		return nil, err
	}
	payment.Picture = *pictures
	return payment, nil
}

func (s *costControlService) DeletePayment(id, organizationID int64, user string, userID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldPayment, err := repo.GetPaymentByID(id)
	if err != nil {
		msg := "获取付款记录失败"
		return errors.New(msg)
	}
	if oldPayment.OrganizationID != organizationID && organizationID != 0 {
		msg := "付款记录不存在或无权限"
		return errors.New(msg)
	}
	if oldPayment.UserID != userID {
		msg := "只能删除自己创建的付款"
		return errors.New(msg)
	}
	err = repo.DeletePayment(id, user)
	if err != nil {
		msg := "删除付款失败"
		return errors.New(msg)
	}
	err = repo.DeletePaymentPicture(id)
	if err != nil {
		msg := "删除付款图片失败"
		return errors.New(msg)
	}
	paymentRequest, err := repo.GetPaymentRequestByID(oldPayment.PaymentRequestID)
	if err != nil {
		msg := "获取请款信息失败"
		return errors.New(msg)
	}
	var paymentRequestUpdate ReqPaymentRequestPaid
	paymentRequestUpdate.Paid = paymentRequest.Paid - oldPayment.Amount
	paymentRequestUpdate.Due = paymentRequest.Due + oldPayment.Amount
	if paymentRequestUpdate.Paid == 0 {
		paymentRequestUpdate.Status = 2
	} else {
		paymentRequestUpdate.Status = 4
	}
	paymentRequestUpdate.User = user
	err = repo.UpdatePaymentRequestPaid(paymentRequest.ID, paymentRequestUpdate)
	if err != nil {
		msg := "更新请款信息失败"
		return errors.New(msg)
	}
	var history ReqPaymentRequestHistoryNew
	history.PaymentRequestID = paymentRequest.ID
	history.OrganizationID = paymentRequest.OrganizationID
	history.User = user
	history.Action = "删除付款"
	history.Remark = "删除付款金额为" + strconv.FormatFloat(oldPayment.Amount, 'f', 2, 64) + "元，"
	if paymentRequestUpdate.Paid == 0 {
		history.Remark += "未付款，当前状态为审核通过"
	} else {
		history.Remark += "未完全付款，当前状态为部分付款"
	}
	history.Content = ""
	_, err = repo.CreatePaymentRequestHistory(history)
	if err != nil {
		msg := "生成付款记录失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}

func (s *costControlService) NewIncome(info ReqIncomeNew) error {
	if info.OrganizationID == 0 {
		msg := "组织ID错误"
		return errors.New(msg)
	}
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
		msg := "获取项目信息失败"
		return errors.New(msg)
	}
	incomeID, err := repo.CreateIncome(info)
	if err != nil {
		msg := "生成收入记录失败"
		return errors.New(msg)
	}
	for _, picture := range info.Picture {
		var paymentPicture ReqIncomePictureNew
		paymentPicture.IncomeID = incomeID
		paymentPicture.Picture = picture
		paymentPicture.User = info.User
		err = repo.CreateIncomePicture(paymentPicture)
		if err != nil {
			msg := "创建收入记录文件失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	return nil

}

func (s *costControlService) UpdateIncome(id int64, info ReqIncomeUpdate) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldIncome, err := repo.GetIncomeByID(id)
	if err != nil {
		msg := "获取收入信息失败"
		return errors.New(msg)
	}
	if oldIncome.UserID != info.UserID {
		msg := "只能更新自己的收入"
		return errors.New(msg)
	}
	if oldIncome.OrganizationID != info.OrganizationID {
		msg := "收入信息不存在"
		return errors.New(msg)
	}
	err = repo.DeleteIncomePicture(id)
	if err != nil {
		msg := "删除收入图片失败"
		return errors.New(msg)
	}
	err = repo.UpdateIncome(id, info)
	if err != nil {
		msg := "更新收入记录失败"
		return errors.New(msg)
	}
	for _, picture := range info.Picture {
		var paymentPicture ReqIncomePictureNew
		paymentPicture.IncomeID = id
		paymentPicture.Picture = picture
		paymentPicture.User = info.User
		err = repo.CreateIncomePicture(paymentPicture)
		if err != nil {
			msg := "创建收入记录文件失败"
			return errors.New(msg)
		}
	}
	tx.Commit()
	return nil
}

func (s *costControlService) GetIncomeList(filter ReqIncomeFilter) (int, *[]RespIncome, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	count, err := query.GetIncomeCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetIncomeList(filter)
	if err != nil {
		return 0, nil, err
	}
	for key, budget := range *list {
		pictures, err := query.GetIncomePictureList(budget.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[key].Picture = *pictures
	}
	return count, list, nil
}

func (s *costControlService) GetIncomeByID(id, organizationID int64) (*RespIncome, error) {
	db := database.InitMySQL()
	query := NewCostControlQuery(db)
	payment, err := query.GetIncomeByID(id)
	if err != nil {
		return nil, err
	}
	if payment.OrganizationID != organizationID && organizationID != 0 {
		msg := "收入不存在或无权限"
		return nil, errors.New(msg)
	}
	pictures, err := query.GetIncomePictureList(id)
	if err != nil {
		return nil, err
	}
	payment.Picture = *pictures
	return payment, nil
}

func (s *costControlService) DeleteIncome(id, organizationID int64, user string, userID int64) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewCostControlRepository(tx)
	oldIncome, err := repo.GetIncomeByID(id)
	if err != nil {
		msg := "获取收入记录失败"
		return errors.New(msg)
	}
	if oldIncome.OrganizationID != organizationID && organizationID != 0 {
		msg := "收入记录不存在或无权限"
		return errors.New(msg)
	}
	if oldIncome.UserID != userID {
		msg := "只能删除自己创建的收入"
		return errors.New(msg)
	}
	err = repo.DeleteIncome(id, user)
	if err != nil {
		msg := "删除收入失败"
		return errors.New(msg)
	}
	err = repo.DeleteIncomePicture(id)
	if err != nil {
		msg := "删除收入图片失败"
		return errors.New(msg)
	}
	tx.Commit()
	return nil
}
