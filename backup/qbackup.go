package backup

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"log"

	"qiniupkg.com/api.v7/conf"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
)

type PutRet struct {
	Hash     string `json:"hash"`
	Key      string `json:"key"`
	Filesize int    `json:"filesize"`
}

type QBackup struct {
	client *kodo.Client
	bucket string
	// domain for download
	domain string
}

func NewQBackup(key string, secret string, domain, bucket string) BackupProvider {
	conf.ACCESS_KEY = key
	conf.SECRET_KEY = secret

	c := kodo.New(0, nil)
	return &QBackup{
		client: c,
		domain: domain,
		bucket: bucket,
	}
}

func (q *QBackup) Info(key string) (entry kodo.Entry, err error) {
	p := q.client.Bucket(q.bucket)

	entry, err = p.Stat(nil, key)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			err = ErrNotFound
		}
		return
	}

	return
}

func (q *QBackup) Get(key string, dirpath string) error {
	baseURL := kodo.MakeBaseUrl(q.domain, key)
	downloadURL := q.client.MakePrivateUrl(baseURL, &kodo.GetPolicy{})

	filepath := fmt.Sprintf("%s/%s", dirpath, key)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("download file failed, status_code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}

func (q *QBackup) Put(filepath string, key string) error {
	if !FileExists(filepath) {
		return ErrFileNotExists
	}

	if key == "" {
		var err error
		key, err = FileKey(filepath)
		if err != nil {
			return err
		}
	}

	policy := kodo.PutPolicy{
		Scope: q.bucket + ":" + key,
	}

	token, err := q.client.MakeUptokenWithSafe(&policy)
	if err != nil {
		return err
	}

	uploader := kodocli.NewUploaderWithoutZone(nil)

	var res PutRet
	err = uploader.PutFile(context.TODO(), &res, token, key, filepath, nil)
	if err != nil {
		return err
	}

	return nil
}

func (q *QBackup) Remove(key string) error {
	p := q.client.Bucket(q.bucket)

	return p.Delete(nil, key)
}

func (q *QBackup) List(prefix string, marker string, limit int) error {
	p := q.client.Bucket(q.bucket)

	entries, cps, marker, err := p.List(context.TODO(), prefix, "", "", limit)
	if err != nil && err != io.EOF {
		return err
	}

	log.Printf("entries: %+v, commonPrefixes: %+v, marker: %s", entries, cps, marker)

	return nil
}
