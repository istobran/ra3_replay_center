package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	AdapterQCloud = "qcloud"
	QAPIDomain    = "monitor.api.qcloud.com"
	QAPIPath      = "/v2/index.php"
	QAPIAction    = "SendCustomAlarmMsg"
)

// 基于腾讯云自定义消息策略的日志上报机制
// https://monitor.api.qcloud.com/v2/index.php?
// &<公共请求参数>
// &policyId=cm-ts6c8ad7
// &msg=hello custom msg
type QCloudWriter struct {
	PolicyId  string `json:"policyId"`  // 策略id
	Region    string `json:"region"`    // 区域参数
	SecretId  string `json:"secretId"`  // 云 API 密钥上的 SecretId
	SecretKey string `json:"secretKey"` // 云 API 密钥上的 SecretKey
	Level     int    `json:"level"`     // 上报等级
}

func NewQCloud() logs.Logger {
	return &QCloudWriter{Level: logs.LevelTrace}
}

func (c *QCloudWriter) Init(jsonConfig string) error {
	return json.Unmarshal([]byte(jsonConfig), c)
}

func (c *QCloudWriter) WriteMsg(when time.Time, msg string, level int) error {
	if level > c.Level {
		return nil
	}
	params := url.Values{
		"Action":    {QAPIAction},
		"Region":    {c.Region},
		"Timestamp": {strconv.FormatInt(when.Unix(), 10)},
		"Nonce":     {strconv.FormatInt(int64(rand.Intn(int(^uint(0)>>1))), 10)},
		"SecretId":  {c.SecretId},
		"policyId":  {c.PolicyId},
		"msg":       {msg},
	}
	// 生成签名
	err := c.generateSignature(&params)
	if err != nil {
		return err
	}
	resp, err := http.PostForm("https://"+QAPIDomain+QAPIPath, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *QCloudWriter) Destroy() {
	return
}

func (c *QCloudWriter) Flush() {
	return
}

func (c *QCloudWriter) generateSignature(params *url.Values) error {
	sortedParams, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return err
	}
	signSource := strings.Join([]string{"POST", QAPIDomain, QAPIPath, "?", sortedParams}, "")
	fmt.Println(signSource)
	secretKey := c.SecretKey
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signSource))
	sha1Sum := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(sha1Sum)
	params.Set("Signature", signature)
	return nil
}

func init() {
	logs.Register(AdapterQCloud, NewQCloud)
}
