package message

import (
	"bpm/api/v1/assignment"
	"bpm/api/v1/auth"
	"bpm/api/v1/event"
	"bpm/api/v1/organization"
	"bpm/api/v1/project"
	"bpm/core/config"
	"bpm/core/database"
	"bpm/core/queue"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/streadway/amqp"
)

type NewProjectCreated struct {
	ProjectID int64 `json:"project_id"`
}

type NewEventUpdated struct {
	EventID int64 `json:"event_id"`
}
type NewEventCompleted struct {
	EventID int64 `json:"event_id"`
}
type NewEventAudited struct {
	EventID int64 `json:"event_id"`
}
type NewProjectReportCreated struct {
	ProjectReportID int64 `json:"project_report_id"`
}
type NewAssignmentCreated struct {
	AssignmentID int64 `json:"assignment_id"`
}
type NewAssignmentCompleted struct {
	AssignmentID int64 `json:"assignment_id"`
}
type todoToSend struct {
	OpenID string `json:"open_id"`
	Thing2 string `json:"thing2"`
	Thing5 string `json:"thing5"`
	Name7  string `json:"name7"`
	Date3  string `json:"date3"`
	Thing8 string `json:"thing8"`
}
type auditToSend struct {
	OpenID  string `json:"open_id"`
	Thing1  string `json:"thing1"`
	Thing2  string `json:"thing2"`
	Thing11 string `json:"thing11"`
	Thing6  string `json:"thing6"`
	Time12  string `json:"time12"`
}
type reportToSend struct {
	OpenID string `json:"open_id"`
	Thing1 string `json:"thing1"`
	Thing3 string `json:"thing3"`
	Thing4 string `json:"thing4"`
	Thing5 string `json:"thing5"`
	Time2  string `json:"time2"`
}
type messageRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type assignmentToSend struct {
	OpenID string `json:"open_id"`
	Thing4 string `json:"thing4"`
	Date3  string `json:"date3"`
	Thing6 string `json:"thing6"`
}
type assignmentAuditToSend struct {
	OpenID string `json:"open_id"`
	Thing1 string `json:"thing1"`
	Thing2 string `json:"thing2"`
	Name4  string `json:"name4"`
	Time5  string `json:"time5"`
}

func Subscribe(conn *queue.Conn) {
	// conn.StartConsumer("NewTodo", "NewProjectCreated", NewTodo)
	conn.StartConsumer("NewTodo", "NewProjectMember", NewTodo)
	conn.StartConsumer("NewEventTodo", "NewEventUpdated", NewEventTodo)
	conn.StartConsumer("NewEventAudit", "NewEventCompleted", NewEventAudit)
	conn.StartConsumer("NewEventAudited", "NewEventAudited", NextEventTodo)
	conn.StartConsumer("NewProjectReportCreated", "NewProjectReportCreated", NewReportTodo)
	conn.StartConsumer("NewAssignmentCreated", "NewAssignmentCreated", NewAssignmentTodo)
	conn.StartConsumer("NewAssignmentCompleted", "NewAssignmentCompleted", NewAssignmentAuditTodo)
}

func NewTodo(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewProjectCreated NewProjectCreated
	err := json.Unmarshal(d.Body, &NewProjectCreated)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	err = sendMessageToActive(NewProjectCreated.ProjectID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func checkExist(slice []todoToSend, find string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i].OpenID == find {
			return true
		}
	}
	return false
}

func checkExist2(slice []auditToSend, find string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i].OpenID == find {
			return true
		}
	}
	return false
}

func checkExist3(slice []reportToSend, find string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i].OpenID == find {
			return true
		}
	}
	return false
}
func NewEventTodo(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewEventUpdated NewEventUpdated
	err := json.Unmarshal(d.Body, &NewEventUpdated)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	err = sendMessageToEvent(NewEventUpdated.EventID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func NewEventAudit(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewEventCompleted NewEventCompleted
	err := json.Unmarshal(d.Body, &NewEventCompleted)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	db := database.InitMySQL()
	eventQuery := event.NewEventQuery(db)
	event, err := eventQuery.GetEventByID(NewEventCompleted.EventID, 0)
	if err != nil {
		fmt.Println(err.Error() + "18")
		return false
	}
	if event.Status != 2 {
		err = sendMessageToActive(event.ProjectID)
		if err != nil {
			fmt.Println(err.Error())
			return false
		} else {
			return true
		}
	} else {
		err = sendMessageToAudit(event.ID)
		if err != nil {
			fmt.Println(err.Error())
			return false
		} else {
			return true
		}
	}
}

func NextEventTodo(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewEventAudited NewEventAudited
	err := json.Unmarshal(d.Body, &NewEventAudited)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	db := database.InitMySQL()
	eventQuery := event.NewEventQuery(db)
	event, err := eventQuery.GetEventByID(NewEventAudited.EventID, 0)
	if err != nil {
		fmt.Println(err.Error() + "18")
		return false
	}
	if event.Status == 3 {
		err = sendMessageToEvent(event.ID)
		if err != nil {
			fmt.Println(err.Error())
			return false
		} else {
			return true
		}
	} else if event.Status == 9 {
		err = sendMessageToActive(event.ProjectID)
		if err != nil {
			fmt.Println(err.Error())
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

func sendMessageToActive(projectID int64) error {
	var toSends []todoToSend
	db := database.InitMySQL()
	query := NewMessageQuery(db)
	eventQuery := event.NewEventQuery(db)
	projectQuery := project.NewProjectQuery(db)
	project, err := projectQuery.GetProjectByID(projectID, 0)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			fmt.Println("项目不存在")
			return nil
		}
		fmt.Println(err.Error() + "4")
		return err
	}
	var filter event.MyEventFilter
	filter.ProjectID = projectID
	filter.Status = "active"
	events, err := eventQuery.GetProjectEvent(filter)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return err
	}
	for _, event := range *events {
		// fmt.Println(event.ProjectName)
		// return false
		active, err := eventQuery.CheckActive(event.ID)
		if err != nil {
			fmt.Println(err.Error() + "3")
			return err
		}
		if !active {
			continue
		}
		assigned, err := eventQuery.GetAssignsByEventID(event.ID)
		if err != nil {
			fmt.Println(err.Error() + "2")
			return err
		}
		for _, assignTo := range *assigned {
			if assignTo.AssignType == 1 { //user
				users, err := query.GetUserByPositionAndProject(assignTo.AssignTo, projectID)
				if err != nil {
					fmt.Println(err.Error() + "1")
					return err
				}
				for _, user := range *users {
					if !checkExist(toSends, user) {
						var msg todoToSend
						msg.OpenID = user
						msg.Thing2 = event.ProjectName
						msg.Thing5 = event.Name
						msg.Name7 = project.CreatedBy
						msg.Date3 = project.Created.Format("2006-01-02 15:04:05")
						if event.Deadline == "" {
							msg.Thing8 = "无备注"
						} else {
							msg.Thing8 = "请在" + event.Deadline + "之前完成"
						}
						toSends = append(toSends, msg)
					}
				}
			} else {
				openID, err := query.GetUserByIDAndProject(assignTo.AssignTo, projectID)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "6")
						return err
					}
				}
				repeat := checkExist(toSends, openID)
				if !repeat {
					var msg todoToSend
					msg.OpenID = openID
					msg.Thing2 = event.ProjectName
					msg.Thing5 = event.Name
					msg.Name7 = project.CreatedBy
					msg.Date3 = project.Created.Format("2006-01-02 15:04:05")
					if event.Deadline == "" {
						msg.Thing8 = "无备注"
					} else {
						msg.Thing8 = "请在" + event.Deadline + "之前完成"
					}
					toSends = append(toSends, msg)
				}
			}
		}
	}
	organizationQuery := organization.NewOrganizationQuery(db)
	for _, toSend := range toSends {
		accessToken, err := organizationQuery.GetAccessToken("bpm")
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				if err != nil {
					fmt.Println(err.Error() + "7")
					return err
				}
			} else {
				var tokenRes organization.WechatToken
				httpClient := &http.Client{}
				token_uri := config.ReadConfig("Wechat.token_uri")
				appID := config.ReadConfig("Wechat.app_id")
				appSecret := config.ReadConfig("Wechat.app_secret")
				uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
				req, err := http.NewRequest("GET", uri, nil)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "8")
						return err
					}
				}
				res, err := httpClient.Do(req)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "9")
						return err
					}
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "10")
						return err
					}
				}
				err = json.Unmarshal(body, &tokenRes)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "11")
						return err
					}
				}
				tx, err := db.Begin()
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "12")
						return err
					}
				}
				defer tx.Rollback()
				repo := organization.NewOrganizationRepository(tx)
				err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "13")
						return err
					}
				}
				tx.Commit()
				accessToken = tokenRes.AccessToken
			}
		}
		url := config.ReadConfig("Wechat.message_uri")
		templateID := config.ReadConfig("Wechat.daiban_template_id")
		state := config.ReadConfig("Wechat.state")
		jsonReq := []byte(`{ "touser" : "` + toSend.OpenID + `", "template_id" : "` + templateID + `", "page" : "pages/index/index","miniprogram_state" : "` + state + `","lang" : "zh_CN","data" : {  "thing2" : { "value": "` + toSend.Thing2 + `"}, "thing5": { "value": "` + toSend.Thing5 + `"}, "name7": { "value": "` + toSend.Name7 + `"}, "date3": { "value": "` + toSend.Date3 + `"}, "thing8": { "value": "` + toSend.Thing8 + `" } } }`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "14")
				return err
			}
		}
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "15")
				return err
			}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "16")
				return err
			}
		}
		var res messageRes
		err = json.Unmarshal(body, &res)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "17")
				return err
			}
		}
		if res.Errcode != 0 {
			fmt.Println(res.Errmsg)
		}
	}
	return nil
}

func sendMessageToAudit(eventID int64) error {
	var toSends []auditToSend
	db := database.InitMySQL()
	query := NewMessageQuery(db)
	eventQuery := event.NewEventQuery(db)
	projectQuery := project.NewProjectQuery(db)
	event, err := eventQuery.GetEventByID(eventID, 0)
	if err != nil {
		fmt.Println(err.Error() + "1")
		return err
	}
	project, err := projectQuery.GetProjectByID(event.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return err
	}
	assigned, err := eventQuery.GetAuditsByEventID(event.ID)
	if err != nil {
		fmt.Println(err.Error() + "2")
		return err
	}
	for _, assignTo := range *assigned {
		if assignTo.AuditType == 1 { //position
			users, err := query.GetUserByPositionAndProject(assignTo.AuditTo, event.ProjectID)
			if err != nil {
				fmt.Println(err.Error() + "1")
				return err
			}
			for _, user := range *users {
				if !checkExist2(toSends, user) {
					var msg auditToSend
					msg.OpenID = user
					msg.Thing1 = project.Name
					msg.Thing2 = event.UpdatedBy
					msg.Thing11 = event.Name
					msg.Thing6 = "有需要你审批的节点"
					msg.Time12 = event.Updated.Format("2006-01-02 15:03:04")
					toSends = append(toSends, msg)
				}
			}
		} else {
			openID, err := query.GetUserByIDAndProject(assignTo.AuditTo, event.ProjectID)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "6")
					return err
				}
			}
			repeat := checkExist2(toSends, openID)
			if !repeat {
				var msg auditToSend
				msg.OpenID = openID
				msg.Thing1 = project.Name
				msg.Thing2 = event.UpdatedBy
				msg.Thing11 = event.Name
				msg.Thing6 = "有需要你审批的节点"
				msg.Time12 = event.Updated.Format("2006-01-02 15:03:04")
				toSends = append(toSends, msg)
			}
		}
	}
	organizationQuery := organization.NewOrganizationQuery(db)
	for _, toSend := range toSends {
		accessToken, err := organizationQuery.GetAccessToken("bpm")
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				if err != nil {
					fmt.Println(err.Error() + "7")
					return err
				}
			} else {
				var tokenRes organization.WechatToken
				httpClient := &http.Client{}
				token_uri := config.ReadConfig("Wechat.token_uri")
				appID := config.ReadConfig("Wechat.app_id")
				appSecret := config.ReadConfig("Wechat.app_secret")
				uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
				req, err := http.NewRequest("GET", uri, nil)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "8")
						return err
					}
				}
				res, err := httpClient.Do(req)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "9")
						return err
					}
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "10")
						return err
					}
				}
				err = json.Unmarshal(body, &tokenRes)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "11")
						return err
					}
				}
				tx, err := db.Begin()
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "12")
						return err
					}
				}
				defer tx.Rollback()
				repo := organization.NewOrganizationRepository(tx)
				err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "13")
						return err
					}
				}
				tx.Commit()
				accessToken = tokenRes.AccessToken
			}
		}
		url := config.ReadConfig("Wechat.message_uri")
		templateID := config.ReadConfig("Wechat.shenpi_template_id")
		state := config.ReadConfig("Wechat.state")
		jsonReq := []byte(`{ "touser" : "` + toSend.OpenID + `", "template_id" : "` + templateID + `", "page" : "pages/index/index","miniprogram_state" : "` + state + `","lang" : "zh_CN","data" : {  "thing1" : { "value": "` + toSend.Thing1 + `"}, "thing2": { "value": "` + toSend.Thing2 + `"}, "thing11": { "value": "` + toSend.Thing11 + `"}, "thing6": { "value": "` + toSend.Thing6 + `"}, "time12": { "value": "` + toSend.Time12 + `" } } }`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "14")
				return err
			}
		}
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "15")
				return err
			}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "16")
				return err
			}
		}
		var res messageRes
		err = json.Unmarshal(body, &res)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "17")
				return err
			}
		}
		if res.Errcode != 0 {
			fmt.Println(res.Errmsg)
		}
	}
	return nil
}

func sendMessageToEvent(eventID int64) error {
	var toSends []todoToSend
	db := database.InitMySQL()
	query := NewMessageQuery(db)
	eventQuery := event.NewEventQuery(db)
	projectQuery := project.NewProjectQuery(db)
	event, err := eventQuery.GetEventByID(eventID, 0)
	if err != nil {
		fmt.Println(err.Error() + "18")
		return err
	}
	project, err := projectQuery.GetProjectByID(event.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return err
	}
	active, err := eventQuery.CheckActive(event.ID)
	if err != nil {
		fmt.Println(err.Error() + "3")
		return err
	}
	if !active {
		return nil
	}
	assigned, err := eventQuery.GetAssignsByEventID(event.ID)
	if err != nil {
		fmt.Println(err.Error() + "2")
		return err
	}
	for _, assignTo := range *assigned {
		if assignTo.AssignType == 1 { //user
			users, err := query.GetUserByPositionAndProject(assignTo.AssignTo, event.ProjectID)
			if err != nil {
				fmt.Println(err.Error() + "1")
				return err
			}
			for _, user := range *users {
				if !checkExist(toSends, user) {
					var msg todoToSend
					msg.OpenID = user
					msg.Thing2 = project.Name
					msg.Thing5 = event.Name
					msg.Name7 = event.UpdatedBy
					msg.Date3 = event.Updated.Format("2006-01-02 15:04:05")
					if event.Deadline == "" {
						msg.Thing8 = "无备注"
					} else {
						msg.Thing8 = "请在" + event.Deadline + "之前完成"
					}
					toSends = append(toSends, msg)
				}
			}
		} else {
			openID, err := query.GetUserByIDAndProject(assignTo.AssignTo, event.ProjectID)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "6")
					return err
				}
			}
			repeat := checkExist(toSends, openID)
			if !repeat {
				var msg todoToSend
				msg.OpenID = openID
				msg.Thing2 = project.Name
				msg.Thing5 = event.Name
				msg.Name7 = event.UpdatedBy
				msg.Date3 = event.Updated.Format("2006-01-02 15:04:05")
				if event.Deadline == "" {
					msg.Thing8 = "无备注"
				} else {
					msg.Thing8 = "请在" + event.Deadline + "之前完成"
				}
				toSends = append(toSends, msg)
			}
		}
	}
	organizationQuery := organization.NewOrganizationQuery(db)
	for _, toSend := range toSends {

		accessToken, err := organizationQuery.GetAccessToken("bpm")
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				if err != nil {
					fmt.Println(err.Error() + "7")
					return err
				}
			} else {
				var tokenRes organization.WechatToken
				httpClient := &http.Client{}
				token_uri := config.ReadConfig("Wechat.token_uri")
				appID := config.ReadConfig("Wechat.app_id")
				appSecret := config.ReadConfig("Wechat.app_secret")
				uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
				req, err := http.NewRequest("GET", uri, nil)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "8")
						return err
					}
				}
				res, err := httpClient.Do(req)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "9")
						return err
					}
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "10")
						return err
					}
				}
				err = json.Unmarshal(body, &tokenRes)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "11")
						return err
					}
				}
				tx, err := db.Begin()
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "12")
						return err
					}
				}
				defer tx.Rollback()
				repo := organization.NewOrganizationRepository(tx)
				err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "13")
						return err
					}
				}
				tx.Commit()
				accessToken = tokenRes.AccessToken
			}
		}
		url := config.ReadConfig("Wechat.message_uri")
		templateID := config.ReadConfig("Wechat.daiban_template_id")
		state := config.ReadConfig("Wechat.state")
		jsonReq := []byte(`{ "touser" : "` + toSend.OpenID + `", "template_id" : "` + templateID + `", "page" : "pages/index/index","miniprogram_state" : "` + state + `","lang" : "zh_CN","data" : {  "thing2" : { "value": "` + toSend.Thing2 + `"}, "thing5": { "value": "` + toSend.Thing5 + `"}, "name7": { "value": "` + toSend.Name7 + `"}, "date3": { "value": "` + toSend.Date3 + `"}, "thing8": { "value": "` + toSend.Thing8 + `" } } }`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "14")
				return err
			}
		}
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "15")
				return err
			}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "16")
				return err
			}
		}
		var res messageRes
		err = json.Unmarshal(body, &res)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "17")
				return err
			}
		}
		if res.Errcode != 0 {
			fmt.Println(res.Errmsg)
		}
	}
	return nil
}

func NewReportTodo(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewProjectReportCreated NewProjectReportCreated
	err := json.Unmarshal(d.Body, &NewProjectReportCreated)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	db := database.InitMySQL()
	projectQuery := project.NewProjectQuery(db)
	fmt.Println(NewProjectReportCreated.ProjectReportID)
	report, err := projectQuery.GetProjectReportByID(NewProjectReportCreated.ProjectReportID, 0)
	if err != nil {
		fmt.Println(err.Error() + "18")
		return false
	}
	err = sendMessageToReport(report.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}
}

func sendMessageToReport(reportID int64) error {
	var toSends []reportToSend
	db := database.InitMySQL()
	query := NewMessageQuery(db)
	projectQuery := project.NewProjectQuery(db)
	report, err := projectQuery.GetProjectReportByID(reportID, 0)
	if err != nil {
		fmt.Println(err.Error(), "get report err")
		return err

	}
	project, err := projectQuery.GetProjectByID(report.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "get project error")
		return err
	}
	users, err := query.GetOtherMemberByProject(project.ID, report.UserID)
	if err != nil {
		fmt.Println(err.Error() + "get user error")
		return err
	}
	for _, user := range users {
		if !checkExist3(toSends, user) {
			var msg reportToSend
			msg.OpenID = user
			msg.Thing1 = "内部报告"
			msg.Thing3 = "有新报告，请查看"
			msg.Thing4 = report.Username
			msg.Thing5 = report.Name
			msg.Time2 = report.Updated.Format("2006-01-02 15:04:05")
			toSends = append(toSends, msg)
		}
	}
	organizationQuery := organization.NewOrganizationQuery(db)
	for _, toSend := range toSends {
		accessToken, err := organizationQuery.GetAccessToken("bpm")
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				if err != nil {
					fmt.Println(err.Error() + "7")
					return err
				}
			} else {
				var tokenRes organization.WechatToken
				httpClient := &http.Client{}
				token_uri := config.ReadConfig("Wechat.token_uri")
				appID := config.ReadConfig("Wechat.app_id")
				appSecret := config.ReadConfig("Wechat.app_secret")
				uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
				req, err := http.NewRequest("GET", uri, nil)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "8")
						return err
					}
				}
				res, err := httpClient.Do(req)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "9")
						return err
					}
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "10")
						return err
					}
				}
				err = json.Unmarshal(body, &tokenRes)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "11")
						return err
					}
				}
				tx, err := db.Begin()
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "12")
						return err
					}
				}
				defer tx.Rollback()
				repo := organization.NewOrganizationRepository(tx)
				err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "13")
						return err
					}
				}
				tx.Commit()
				accessToken = tokenRes.AccessToken
			}
		}
		url := config.ReadConfig("Wechat.message_uri")
		templateID := config.ReadConfig("Wechat.report_template_id")
		state := config.ReadConfig("Wechat.state")
		jsonReq := []byte(`{ "touser" : "` + toSend.OpenID + `", "template_id" : "` + templateID + `", "page" : "pages/index/index","miniprogram_state" : "` + state + `","lang" : "zh_CN","data" : {  "thing1" : { "value": "` + toSend.Thing1 + `"}, "thing3": { "value": "` + toSend.Thing3 + `"}, "thing4": { "value": "` + toSend.Thing4 + `"}, "time2": { "value": "` + toSend.Time2 + `"}, "thing5": { "value": "` + toSend.Thing5 + `" } } }`)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "14")
				return err
			}
		}
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "15")
				return err
			}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "16")
				return err
			}
		}
		var res messageRes
		err = json.Unmarshal(body, &res)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "17")
				return err
			}
		}
		if res.Errcode != 0 {
			fmt.Println(res.Errmsg)
		}
	}
	return nil
}

func NewAssignmentTodo(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewAssignmentCreated NewAssignmentCreated
	err := json.Unmarshal(d.Body, &NewAssignmentCreated)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	db := database.InitMySQL()
	assignmentQuery := assignment.NewAssignmentQuery(db)
	fmt.Println(NewAssignmentCreated.AssignmentID)
	assignment, err := assignmentQuery.GetAssignmentByID(NewAssignmentCreated.AssignmentID, 0)
	if err != nil {
		fmt.Println(err.Error() + "28")
		return true
	}
	err = sendMessageToAssignment(assignment.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}
}

func sendMessageToAssignment(assignmentID int64) error {
	db := database.InitMySQL()
	assignmentQuery := assignment.NewAssignmentQuery(db)
	authQuery := auth.NewAuthQuery(db)
	projectQuery := project.NewProjectQuery(db)
	organizationQuery := organization.NewOrganizationQuery(db)
	assignment, err := assignmentQuery.GetAssignmentByID(assignmentID, 0)
	if err != nil {
		fmt.Println(err.Error(), "get assignment err")
		return err

	}
	user, err := authQuery.GetUserByID(assignment.AssignTo, 0)
	if err != nil {
		fmt.Println(err.Error() + "get user error")
		return err
	}
	project, err := projectQuery.GetProjectByID(assignment.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "get project error")
		return err
	}
	var msgToSend assignmentToSend
	msgToSend.OpenID = user.Identifier
	msgToSend.Thing4 = assignment.Name
	msgToSend.Thing6 = project.Name
	msgToSend.Date3 = assignment.Created.Format("2006-01-02 15:04:05")
	accessToken, err := organizationQuery.GetAccessToken("bpm")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			if err != nil {
				fmt.Println(err.Error() + "7")
				return err
			}
		} else {
			var tokenRes organization.WechatToken
			httpClient := &http.Client{}
			token_uri := config.ReadConfig("Wechat.token_uri")
			appID := config.ReadConfig("Wechat.app_id")
			appSecret := config.ReadConfig("Wechat.app_secret")
			uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "8")
					return err
				}
			}
			res, err := httpClient.Do(req)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "9")
					return err
				}
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "10")
					return err
				}
			}
			err = json.Unmarshal(body, &tokenRes)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "11")
					return err
				}
			}
			tx, err := db.Begin()
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "12")
					return err
				}
			}
			defer tx.Rollback()
			repo := organization.NewOrganizationRepository(tx)
			err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "13")
					return err
				}
			}
			tx.Commit()
			accessToken = tokenRes.AccessToken
		}
	}
	url := config.ReadConfig("Wechat.message_uri")
	templateID := config.ReadConfig("Wechat.assignment_template_id")
	state := config.ReadConfig("Wechat.state")
	jsonReq := []byte(`{ "touser" : "` + msgToSend.OpenID + `", "template_id" : "` + templateID + `", "page" : "pages/index/index","miniprogram_state" : "` + state + `","lang" : "zh_CN","data" : {  "thing4" : { "value": "` + msgToSend.Thing4 + `"}, "thing6": { "value": "` + msgToSend.Thing6 + `"}, "date3": { "value": "` + msgToSend.Date3 + `"} } }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "14")
			return err
		}
	}
	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "15")
			return err
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "16")
			return err
		}
	}
	var res messageRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "17")
			return err
		}
	}
	if res.Errcode != 0 {
		fmt.Println(res.Errmsg)
	}
	return nil
}

func NewAssignmentAuditTodo(d amqp.Delivery) bool {
	if d.Body == nil {
		return false
	}
	var NewAssignmentCompleted NewAssignmentCompleted
	err := json.Unmarshal(d.Body, &NewAssignmentCompleted)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "5")
			return false
		}
	}
	db := database.InitMySQL()
	assignmentQuery := assignment.NewAssignmentQuery(db)
	fmt.Println(NewAssignmentCompleted.AssignmentID)
	assignment, err := assignmentQuery.GetAssignmentByID(NewAssignmentCompleted.AssignmentID, 0)
	if err != nil {
		fmt.Println(err.Error() + "28")
		return true
	}
	err = sendMessageToAssignmentAudit(assignment.ID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		return true
	}
}

func sendMessageToAssignmentAudit(assignmentID int64) error {
	db := database.InitMySQL()
	assignmentQuery := assignment.NewAssignmentQuery(db)
	authQuery := auth.NewAuthQuery(db)
	projectQuery := project.NewProjectQuery(db)
	organizationQuery := organization.NewOrganizationQuery(db)
	assignment, err := assignmentQuery.GetAssignmentByID(assignmentID, 0)
	if err != nil {
		fmt.Println(err.Error(), "get assignment err")
		return err

	}
	user, err := authQuery.GetUserByID(assignment.AuditTo, 0)
	if err != nil {
		fmt.Println(err.Error() + "get user error")
		return err
	}
	project, err := projectQuery.GetProjectByID(assignment.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "get project error")
		return err
	}
	var msgToSend assignmentAuditToSend
	msgToSend.OpenID = user.Identifier
	msgToSend.Thing1 = assignment.Name
	msgToSend.Thing2 = project.Name
	msgToSend.Time5 = assignment.CompleteTime[0:19]
	accessToken, err := organizationQuery.GetAccessToken("bpm")
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			if err != nil {
				fmt.Println(err.Error() + "7")
				return err
			}
		} else {
			var tokenRes organization.WechatToken
			httpClient := &http.Client{}
			token_uri := config.ReadConfig("Wechat.token_uri")
			appID := config.ReadConfig("Wechat.app_id")
			appSecret := config.ReadConfig("Wechat.app_secret")
			uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
			req, err := http.NewRequest("GET", uri, nil)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "8")
					return err
				}
			}
			res, err := httpClient.Do(req)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "9")
					return err
				}
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "10")
					return err
				}
			}
			err = json.Unmarshal(body, &tokenRes)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "11")
					return err
				}
			}
			tx, err := db.Begin()
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "12")
					return err
				}
			}
			defer tx.Rollback()
			repo := organization.NewOrganizationRepository(tx)
			err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "13")
					return err
				}
			}
			tx.Commit()
			accessToken = tokenRes.AccessToken
		}
	}
	url := config.ReadConfig("Wechat.message_uri")
	templateID := config.ReadConfig("Wechat.assignment_audit_template_id")
	state := config.ReadConfig("Wechat.state")
	jsonReq := []byte(`{ "touser" : "` + msgToSend.OpenID + `", "template_id" : "` + templateID + `", "page" : "pages/index/index","miniprogram_state" : "` + state + `","lang" : "zh_CN","data" : {  "thing1" : { "value": "` + msgToSend.Thing1 + `"}, "thing2": { "value": "` + msgToSend.Thing2 + `"}, "time5": { "value": "` + msgToSend.Time5 + `"}, "name4": { "value": "` + msgToSend.Name4 + `"} } }`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "14")
			return err
		}
	}
	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "15")
			return err
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "16")
			return err
		}
	}
	var res messageRes
	err = json.Unmarshal(body, &res)
	if err != nil {
		if err != nil {
			fmt.Println(err.Error() + "17")
			return err
		}
	}
	if res.Errcode != 0 {
		fmt.Println(res.Errmsg)
	}
	return nil
}
