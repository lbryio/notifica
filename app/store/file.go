package store

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/notifica/app/types"
	"github.com/pierrec/lz4"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func newFileStore(bucket string) *file {
	return &file{
		bucket:       bucket,
		bucketPrefix: bucket[:bucketPrefixLength],
		bucketSuffix: bucket[bucketPrefixLength:],
	}
}

type file struct {
	r            *bytes.Reader
	bucket       string
	bucketPrefix string
	bucketSuffix string
}

func (f *file) Add(notification *types.Notification) error {

	n, err := proto.Marshal(notification)
	if err != nil {
		return errors.Err(err)
	}
	// Compress.
	r := bytes.NewReader(n)
	logrus.Debugf("notification size %d bytes", len(n))
	var zout bytes.Buffer
	zw := lz4.NewWriter(&zout)
	_, err = io.Copy(zw, r)
	if err != nil {
		return errors.Err(err)
	}
	err = zw.Close()
	if err != nil {
		return errors.Err(err)
	}
	var fileBuf bytes.Buffer
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(zout.Len()))
	fileBuf.Write(b)
	fileBuf.Write(zout.Bytes())
	file, err := os.OpenFile(f.filePath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Err(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	written, err := w.Write(fileBuf.Bytes())
	if err != nil {
		return errors.Err(err)
	}
	err = w.Flush()
	if err != nil {
		return errors.Err(err)
	}
	logrus.Debugf("appended %d bytes to %s", written, f.filePath())
	return nil
}

func (f *file) filePath() string {
	filePathDir := storagePath + "/" + f.bucketPrefix
	filePath := filePathDir + "/" + f.bucket + ".beams"
	err := os.MkdirAll(filePathDir, 0755)
	if err != nil {
		logrus.Error(errors.FullTrace(err))
	}
	return filePath
}

func (f *file) Next() (*types.Notification, error) {
	if f.r == nil {
		buf, err := ioutil.ReadFile(f.filePath())
		if err != nil {
			buf = make([]byte, 0)
		}
		f.r = bytes.NewReader(buf)
	}
	if f.r.Len() == 0 {
		return nil, nil
	}
	b2readBuf := make([]byte, 2)
	_, err := f.r.Read(b2readBuf)
	if err != nil {
		return nil, errors.Err(err)
	}
	bytesToRead := int16(binary.LittleEndian.Uint16(b2readBuf))
	compNotifBuf := make([]byte, bytesToRead)
	_, err = f.r.Read(compNotifBuf)
	if err != nil {
		return nil, errors.Err(err)
	}
	var uncompNotifBuf bytes.Buffer
	zr := lz4.NewReader(bytes.NewBuffer(compNotifBuf))
	_, err = io.Copy(&uncompNotifBuf, zr)
	if err != nil {
		return nil, errors.Err(err)
	}
	notification := types.Notification{}
	err = proto.Unmarshal(uncompNotifBuf.Bytes(), &notification)
	if err != nil {
		return nil, errors.Err(err)
	}

	return &notification, nil
}
