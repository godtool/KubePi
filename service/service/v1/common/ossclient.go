package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSClient struct {
	client *oss.Client

	bucket string
}

func InitOssClient(endpoint string) (*OSSClient, error) {
	const accessIDKey = "accessID"
	const accessIDSecret = "accessSecret"

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse oss url(%s) failed. %v", endpoint, err)
	}

	host := u.Host
	accessId := u.Query().Get(accessIDKey)
	accessKey := u.Query().Get(accessIDSecret)
	bucket := u.Query().Get("bucket")

	client, err := oss.New(host, accessId, accessKey)
	if err != nil {
		str := fmt.Sprintf("init oss client fail.endpoint:%v accessId:%v", host, accessId)
		return nil, fmt.Errorf(str)
	}
	cli := &OSSClient{client: client, bucket: bucket}

	return cli, nil
}

func (o *OSSClient) PutObject(url string, data []byte, retryCount int, contentType string, ctx context.Context) error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = o.putObject(url, data, contentType)
		if err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 200)
	}

	return err

}

func (o *OSSClient) putObject(url string, data []byte, contentType string) error {
	err, bucketName, object := o.parseUrl(url)
	if err != nil {
		return err
	}

	if o.bucket != "" {
		bucketName = o.bucket
	}

	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Case 1: Upload an object from a string

	var options []oss.Option

	if contentType != "" {
		options = append(options, oss.ContentType(contentType))
	}

	err = bucket.PutObject(object, bytes.NewReader(data), options...)
	if err != nil {
		return err
	}

	return nil
}

func (o *OSSClient) AuthURL(url string) (string, error) {
	err, bucketName, object := o.parseUrl(url)
	if err != nil {
		return "", err
	}

	if o.bucket != "" {
		bucketName = o.bucket
	}

	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	return bucket.SignURL(object, oss.HTTPGet, int64(7*24*time.Hour))
}

func (o *OSSClient) parseUrl(url string) (error, string, string) {
	prefix := "oss://"
	if !strings.HasPrefix(url, prefix) {
		return fmt.Errorf("invalid oss Url, not start with oss://. Url:%s", url), "", ""
	}

	s := len(prefix)
	e := strings.Index(url[s:], "/")
	if e == -1 {
		return fmt.Errorf("invalid oss Url. get bucket fail. Url:%s", url), "", ""
	}

	bucketName := url[s : s+e]

	object := url[s+e+1:]

	return nil, bucketName, object
}

func (o *OSSClient) GetObject(url string, ctx context.Context) (string, error) {
	// 1. get data from oss bucket
	data, err := o.getObject(url)

	// 2.1 if got err: log err & return err
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (o *OSSClient) getObject(url string) ([]byte, error) {
	err, bucketName, object := o.parseUrl(url)
	if err != nil {
		return nil, err
	}

	if o.bucket != "" {
		bucketName = o.bucket
	}

	bucket, err := o.client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	body, err := bucket.GetObject(object)
	if err != nil {
		return nil, err
	}

	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
