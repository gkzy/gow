package oss

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gkzy/gow/lib/util"
	"io"
	"net/http"
	"time"
)

//AliClient AliClient
type AliClient struct {
	AccessKeyId string
	Secret      string
	EndPoint    string
	BucketName  string
	ServerUrl   string
}

//NewAliClient NewAliClient
func NewAliClient(accessKeyId, secret, endPoint, bucketName, serverUrl string) *AliClient {
	//末尾添加/
	if len(serverUrl) > 0 && serverUrl[len(serverUrl)-1:] != "/" {
		serverUrl = serverUrl + "/"
	}
	return &AliClient{
		AccessKeyId: accessKeyId,
		Secret:      secret,
		EndPoint:    endPoint,
		BucketName:  bucketName,
		ServerUrl:   serverUrl,
	}
}

// UploadFile 上传一个文件，返回远程地址及错误
//	不强制重命名文件名
//	like:   /dir/20210309/filename
func (c *AliClient) UploadFile(reader io.Reader, dir, filename string) (url string, err error) {
	if filename == "" {
		err = errors.New("请传入上传后的文件名")
		return
	}
	if dir == "" {
		err = errors.New("请传入上传后的目录")
	}
	client, err := oss.New(c.EndPoint, c.AccessKeyId, c.Secret)
	if err != nil {
		err = fmt.Errorf("[client]init失败:%v", err)
		return
	}
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return
	}

	// like /topic/20210309/filename
	filePath := fmt.Sprintf("%s/%s/%s", dir, time.Now().Format("20060102"), filename)
	err = bucket.PutObject(filePath, reader)
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s%s", c.ServerUrl, filePath)
	return
}

// Upload 上传文件	返回远程地址及错误
//		url,err:=client.Upload(reader,dir,ext)
//		会强制重命名文件名
func (c *AliClient) Upload(reader io.Reader, dir, ext string) (url string, err error) {
	if ext == "" {
		ext = ".jpg"
	}
	client, err := oss.New(c.EndPoint, c.AccessKeyId, c.Secret)
	if err != nil {
		err = fmt.Errorf("[client]init失败:%v", err)
		return
	}
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return
	}
	uuid, _ := util.GetUUID()
	filePath := fmt.Sprintf("%s/%s/%s", dir, time.Now().Format("20060102"), uuid+ext)
	err = bucket.PutObject(filePath, reader)
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s%s", c.ServerUrl, filePath)
	return
}

// UploadRemoteFile 上传网络图片到oss
func (c *AliClient) UploadRemoteFile(httpUrl, dir string) (url string, err error) {
	resp, err := http.Get(httpUrl)
	if err != nil {
		err = fmt.Errorf("远程图片获取失败:%v", httpUrl)
		return
	}
	defer resp.Body.Close()
	return c.Upload(resp.Body, dir, "")
}
