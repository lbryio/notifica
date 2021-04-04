package store

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/lbryio/notifica/app/env"
	"github.com/lbryio/notifica/app/types"
)

var storagePath string
var bucketPrefixLength int = 1

type Store interface {
	Add(*types.Notification) error
	Next() (*types.Notification, error)
}

func Init(conf *env.Config) {
	storagePath, _ = os.Getwd()
	storagePath = storagePath + "/data"
	if conf.BucketPath != "" {
		storagePath = conf.BucketPath
	}
	if conf.BucketPrefixLength != 0 {
		bucketPrefixLength = conf.BucketPrefixLength
	}
	logrus.Debug("storage path: ", storagePath)
	logrus.Debug("bucket prefix length: ", bucketPrefixLength)
}

func Retrieve(bucket string) Store {
	return newFileStore(bucket)
}
