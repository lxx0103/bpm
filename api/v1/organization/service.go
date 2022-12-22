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

func (s *organizationService) GetOrganizationByID(id int64) (*Organization, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	organization, err := query.GetOrganizationByID(id)
	return organization, err
}

func (s *organizationService) NewOrganization(info OrganizationNew) (*Organization, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	organizationID, err := repo.CreateOrganization(info)
	if err != nil {
		return nil, err
	}
	organization, err := repo.GetOrganizationByID(organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return organization, err
}

func (s *organizationService) GetOrganizationList(filter OrganizationFilter) (int, *[]Organization, error) {
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
	return count, list, err
}

func (s *organizationService) UpdateOrganization(organizationID int64, info OrganizationNew) (*Organization, error) {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	repo := NewOrganizationRepository(tx)
	_, err = repo.UpdateOrganization(organizationID, info)
	if err != nil {
		return nil, err
	}
	organization, err := repo.GetOrganizationByID(organizationID)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return organization, err
}

func (s *organizationService) GetQrCodeByPath(path string) (string, error) {
	db := database.InitMySQL()
	query := NewOrganizationQuery(db)
	res, err := query.GetQrCodeByPath(path)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return "", err
		} else {
			accessToken, err := query.GetAccessToken("bpm")
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return "", err
				} else {
					var tokenRes WechatToken
					httpClient := &http.Client{}
					token_uri := config.ReadConfig("Wechat.token_uri")
					appID := config.ReadConfig("Wechat.app_id")
					appSecret := config.ReadConfig("Wechat.app_secret")
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
					err = repo.NewAccessToken("bpm", tokenRes.AccessToken)
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
			err = repo.NewQrcode(path, newName)
			if err != nil {
				return "", err
			}
			tx.Commit()
			return newName, nil
		}
	}
	return res, nil
}
