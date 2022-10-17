package message

import (
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
type messageToSend struct {
	OpenID string `json:"open_id"`
	Thing2 string `json:"thing2"`
	Thing5 string `json:"thing5"`
	Name7  string `json:"name7"`
	Date3  string `json:"date3"`
	Thing8 string `json:"thing8"`
}

type messageRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func Subscribe(conn *queue.Conn) {
	conn.StartConsumer("NewTodo", "NewProjectCreated", NewTodo)
	conn.StartConsumer("NewEventTodo", "NewEventUpdated", NewEventTodo)
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
	var toSends []messageToSend
	db := database.InitMySQL()
	query := NewMessageQuery(db)
	eventQuery := event.NewEventQuery(db)
	projectQuery := project.NewProjectQuery(db)
	project, err := projectQuery.GetProjectByID(NewProjectCreated.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return false
	}
	var filter event.MyEventFilter
	filter.ProjectID = NewProjectCreated.ProjectID
	filter.Status = "active"
	events, err := eventQuery.GetProjectEvent(filter)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return false
	}
	for _, event := range *events {
		// fmt.Println(event.ProjectName)
		// return false
		active, err := eventQuery.CheckActive(event.ID)
		if err != nil {
			fmt.Println(err.Error() + "3")
			return false
		}
		if !active {
			continue
		}
		assigned, err := eventQuery.GetAssignsByEventID(event.ID)
		if err != nil {
			fmt.Println(err.Error() + "2")
			return false
		}
		for _, assignTo := range *assigned {
			if assignTo.AssignType == 1 { //user
				users, err := query.GetUserByPosition(assignTo.AssignTo)
				if err != nil {
					fmt.Println(err.Error() + "1")
					return false
				}
				for _, user := range *users {
					if !checkExist(toSends, user) {
						var msg messageToSend
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
				openID, err := query.GetUserByID(assignTo.AssignTo)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "6")
						return false
					}
				}
				repeat := checkExist(toSends, openID)
				if !repeat {
					var msg messageToSend
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
					return false
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
						return false
					}
				}
				res, err := httpClient.Do(req)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "9")
						return false
					}
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "10")
						return false
					}
				}
				err = json.Unmarshal(body, &tokenRes)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "11")
						return false
					}
				}
				tx, err := db.Begin()
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "12")
						return false
					}
				}
				defer tx.Rollback()
				repo := organization.NewOrganizationRepository(tx)
				err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "13")
						return false
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
				return false
			}
		}
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "15")
				return false
			}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "16")
				return false
			}
		}
		var res messageRes
		err = json.Unmarshal(body, &res)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "17")
				return false
			}
		}
		if res.Errcode != 0 {
			fmt.Println(res.Errmsg)
			return true
		}
	}

	return true
}

func checkExist(slice []messageToSend, find string) bool {
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
	var toSends []messageToSend
	db := database.InitMySQL()
	query := NewMessageQuery(db)
	eventQuery := event.NewEventQuery(db)
	projectQuery := project.NewProjectQuery(db)
	event, err := eventQuery.GetEventByID(NewEventUpdated.EventID, 0)
	if err != nil {
		fmt.Println(err.Error() + "18")
		return false
	}
	project, err := projectQuery.GetProjectByID(event.ProjectID, 0)
	if err != nil {
		fmt.Println(err.Error() + "4")
		return false
	}
	active, err := eventQuery.CheckActive(event.ID)
	if err != nil {
		fmt.Println(err.Error() + "3")
		return false
	}
	if !active {
		return true
	}
	assigned, err := eventQuery.GetAssignsByEventID(event.ID)
	if err != nil {
		fmt.Println(err.Error() + "2")
		return false
	}
	for _, assignTo := range *assigned {
		if assignTo.AssignType == 1 { //user
			users, err := query.GetUserByPosition(assignTo.AssignTo)
			if err != nil {
				fmt.Println(err.Error() + "1")
				return false
			}
			for _, user := range *users {
				if !checkExist(toSends, user) {
					var msg messageToSend
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
			openID, err := query.GetUserByID(assignTo.AssignTo)
			if err != nil {
				if err != nil {
					fmt.Println(err.Error() + "6")
					return false
				}
			}
			repeat := checkExist(toSends, openID)
			if !repeat {
				var msg messageToSend
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
					return false
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
						return false
					}
				}
				res, err := httpClient.Do(req)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "9")
						return false
					}
				}
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "10")
						return false
					}
				}
				err = json.Unmarshal(body, &tokenRes)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "11")
						return false
					}
				}
				tx, err := db.Begin()
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "12")
						return false
					}
				}
				defer tx.Rollback()
				repo := organization.NewOrganizationRepository(tx)
				err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
				if err != nil {
					if err != nil {
						fmt.Println(err.Error() + "13")
						return false
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
				return false
			}
		}
		q := req.URL.Query()
		q.Add("access_token", accessToken)
		req.URL.RawQuery = q.Encode()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "15")
				return false
			}
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "16")
				return false
			}
		}
		var res messageRes
		err = json.Unmarshal(body, &res)
		if err != nil {
			if err != nil {
				fmt.Println(err.Error() + "17")
				return false
			}
		}
		if res.Errcode != 0 {
			fmt.Println(res.Errmsg)
			return true
		}
	}

	return true
}
