package organization

import (
	"bpm/core/config"
	"bpm/core/database"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

type organizationService struct {
}

func NewOrganizationService() *organizationService {
	return &organizationService{}
}

func (s *organizationService) GetOrganizationByID(id int64) (*OrganizationResponse, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	organization, err := query.GetOrganizationByID(id)
	if err != nil {
		return nil, err
	}
	qrcodes, err := query.GetOrganizationQrcode(id)
	if err != nil {
		return nil, err
	}
	organization.Qrcode = *qrcodes
	return organization, nil
}

func (s *organizationService) NewOrganization(info OrganizationNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	organizationID, err := repo.CreateOrganization(info)
	if err != nil {
		return err
	}
	if len(info.Qrcode) > 0 {
		for _, qrcode := range info.Qrcode {
			err = repo.CreateOrganizationQrcode(organizationID, qrcode, info.User)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *organizationService) GetOrganizationList(filter OrganizationFilter) (int, *[]OrganizationResponse, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	count, err := query.GetOrganizationCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetOrganizationList(filter)
	if err != nil {
		return 0, nil, err
	}
	for k, v := range *list {
		qrcodes, err := query.GetOrganizationQrcode(v.ID)
		if err != nil {
			return 0, nil, err
		}
		(*list)[k].Qrcode = *qrcodes
	}
	return count, list, err
}

func (s *organizationService) UpdateOrganization(organizationID int64, info OrganizationNew) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	err = repo.DeleteOrganizationQrcode(organizationID, info.User)
	if err != nil {
		return err
	}
	_, err = repo.UpdateOrganization(organizationID, info)
	if err != nil {
		return err
	}
	if len(info.Qrcode) > 0 {
		for _, qrcode := range info.Qrcode {
			err = repo.CreateOrganizationQrcode(organizationID, qrcode, info.User)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	return nil
}

func (s *organizationService) GetQrCodeByPath(path, source string) (string, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	res, err := query.GetQrCodeByPath(path, source)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return "", err
		} else {
			accessToken, err := query.GetAccessToken(source)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return "", err
				} else {
					var tokenRes WechatToken
					httpClient := &http.Client{}
					var appID, appSecret string
					token_uri := config.ReadConfig("Wechat.token_uri")
					if source == "bpm" {
						appID = config.ReadConfig("Wechat.app_id")
						appSecret = config.ReadConfig("Wechat.app_secret")
					} else if source == "portal" {
						appID = config.ReadConfig("PortalWechat.app_id")
						appSecret = config.ReadConfig("PortalWechat.app_secret")
					}
					uri := token_uri + "?appid=" + appID + "&secret=" + appSecret + "&grant_type=client_credential"
					req, err := http.NewRequest("GET", uri, nil)
					if err != nil {
						return "", err
					}
					res, err := httpClient.Do(req)
					if err != nil {
						return "", err
					}
					defer res.Body.Close()
					body, err := ioutil.ReadAll(res.Body)
					if err != nil {
						return "", err
					}
					err = json.Unmarshal(body, &tokenRes)
					if err != nil {
						return "", err
					}
					tx, err := db.Begin()
					if err != nil {
						return "", err
					}
					defer tx.Rollback()
					repo := NewOrganizationRepository(tx)
					err = repo.NewAccessToken(source, tokenRes.AccessToken)
					if err != nil {
						return "", err
					}
					tx.Commit()
					accessToken = tokenRes.AccessToken
				}
			}
			jsonReq := []byte(`{ "path" : "` + path + `", "width" : 430 }`)
			qrcode_uri := config.ReadConfig("Wechat.qrcode_uri") + "?access_token=" + accessToken
			req, err := http.NewRequest("POST", qrcode_uri, bytes.NewBuffer(jsonReq))
			if err != nil {
				return "", err
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}
			dest := config.ReadConfig("file.path")
			extension := ".png"
			newName := uuid.NewString() + extension
			imgPath := dest + newName
			err = ioutil.WriteFile(imgPath, body, 0666)
			if err != nil {
				return "", err
			}
			tx, err := db.Begin()
			if err != nil {
				return "", err
			}
			defer tx.Rollback()
			repo := NewOrganizationRepository(tx)
			err = repo.NewQrcode(path, source, newName)
			if err != nil {
				return "", err
			}
			tx.Commit()
			return newName, nil
		}
	}
	return res, nil
}

func (s *organizationService) GetPortalOrganizationList(filter OrganizationFilter) (int, *[]OrganizationExampleResponse, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	count, err := query.GetOrganizationCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetOrganizationList(filter)
	if err != nil {
		return 0, nil, err
	}
	var res []OrganizationExampleResponse
	for _, organization := range *list {
		var resRow OrganizationExampleResponse
		resRow.ID = organization.ID
		resRow.Name = organization.Name
		resRow.Logo = organization.Logo
		resRow.Description = organization.Description
		resRow.Phone = organization.Phone
		resRow.Contact = organization.Contact
		resRow.Address = organization.Address
		resRow.City = organization.City
		resRow.Type = organization.Type
		resRow.Status = organization.Status
		examples, err := query.GetOrganizationTopExamples(organization.ID)
		if err != nil {
			return 0, nil, err
		}
		qrcodes, err := query.GetOrganizationQrcode(organization.ID)
		if err != nil {
			return 0, nil, err
		}
		resRow.Examples = *examples
		resRow.Qrcode = *qrcodes
		res = append(res, resRow)
	}
	return count, &res, err
}
