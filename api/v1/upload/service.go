package upload

import (
	"bpm/core/config"
	"bpm/core/database"
	"fmt"
	"time"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type uploadService struct {
}

func NewUploadService() uploadService {
	return uploadService{}
}

func (s *uploadService) NewUpload(path string, organizationID int64, userName string) error {
	db := database.InitMySQL()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	repo := NewUploadRepository(tx)
	err = repo.CreateUpload(path, organizationID, userName)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (s *uploadService) GetUploadList(filter UploadFilter, organizationID int64) (int, *[]Upload, error) {
	if organizationID != 0 {
		filter.OrganizationID = organizationID
	}
	db := database.InitMySQL()
	query := NewUploadQuery(db)
	count, err := query.GetUploadCount(filter)
	if err != nil {
		return 0, nil, err
	}
	list, err := query.GetUploadList(filter)
	if err != nil {
		return 0, nil, err
	}
	return count, list, nil
}

func (s *uploadService) GetUploadKey(filter KeyFilter) (*KeyRes, error) {
	appid := filter.APPID
	bucket := filter.Bucket
	secretID := config.ReadConfig("Upload.secret_id")
	secretKey := config.ReadConfig("Upload.secret_key")
	c := sts.NewClient(
		// 通过环境变量获取密钥, os.Getenv 方法表示获取环境变量
		secretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		secretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		nil,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          "ap-guangzhou",
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					// 密钥的权限列表。简单上传和分片需要以下的权限，其他权限列表请看 https://cloud.tencent.com/document/product/436/31923
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						// 这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						// 存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						"qcs::cos:ap-guangzhou:uid/" + appid + ":" + bucket + "/exampleobject",
					},
					// 开始构建生效条件 condition
					// 关于 condition 的详细设置规则和COS支持的condition类型可以参考https://cloud.tencent.com/document/product/436/71306
					Condition: map[string]map[string]interface{}{
						"ip_equal": map[string]interface{}{
							"qcs:ip": []string{
								"10.217.182.3/24",
								"111.21.33.72/24",
							},
						},
					},
				},
			},
		},
	}

	// case 1 请求临时密钥
	res, err := c.GetCredential(opt)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.Credentials)
	var key KeyRes
	key.TmpSecretId = res.Credentials.TmpSecretID
	key.TmpSecretKey = res.Credentials.TmpSecretKey
	key.StartTime = res.StartTime
	key.Expiration = res.Expiration
	key.ExpiredTime = res.ExpiredTime
	key.Token = res.Credentials.SessionToken
	return &key, nil
}
