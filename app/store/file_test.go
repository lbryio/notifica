package store

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"

	"github.com/pierrec/lz4"
)

const testString = "Blah blah blah blah blah"

func TestLZ4(t *testing.T) {
	//Write it
	test := []byte(testString)
	r := bytes.NewReader(test)
	var zout bytes.Buffer
	zw := lz4.NewWriter(&zout)
	//Compress it
	_, err := io.Copy(zw, r)
	if err != nil {
		t.Error(err)
	}
	err = zw.Close()
	if err != nil {
		t.Fatal(err)
	}
	var fileBuf bytes.Buffer
	b := make([]byte, 2)
	bytesToWrite := uint16(zout.Len())
	binary.BigEndian.PutUint16(b, bytesToWrite)
	fileBuf.Write(b)
	fileBuf.Write(zout.Bytes())

	fr := bytes.NewBuffer(fileBuf.Bytes())
	b2readBuf := make([]byte, 2)
	_, err = fr.Read(b2readBuf)
	if err != nil {
		t.Error(err)
	}
	bytesToRead := binary.BigEndian.Uint16(b2readBuf)
	if bytesToRead != bytesToWrite {
		t.Errorf("expected %d but found %d", bytesToWrite, bytesToRead)
	}
	compNotifBuf := make([]byte, bytesToRead)
	_, err = fr.Read(compNotifBuf)
	if err != nil {
		t.Error(err)
	}
	var uncompNotifBuf bytes.Buffer
	compressedBuffer := bytes.NewBuffer(compNotifBuf)
	zr := lz4.NewReader(compressedBuffer)
	read, err := io.Copy(&uncompNotifBuf, zr)
	if err != nil {
		t.Error(err, "read ", read, "bytes")
	}
	if testString != string(uncompNotifBuf.Bytes()) {
		t.Errorf("expected '%s' got '%s'", testString, string(uncompNotifBuf.Bytes()))
	}
	println(string(uncompNotifBuf.Bytes()))
}

func TestLZ4Basic(t *testing.T) {
	raw := []byte(testString)
	r := bytes.NewReader(raw)

	// Compress.
	var zout bytes.Buffer
	zw := lz4.NewWriter(&zout)
	_, err := io.Copy(zw, r)
	if err != nil {
		t.Fatal(err)
	}
	err = zw.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Uncompress.
	var out bytes.Buffer
	zr := lz4.NewReader(&zout)
	n, err := io.Copy(&out, zr)
	if err != nil {
		t.Fatal(err)
	}

	// The uncompressed data must be the same as the initial input.
	if got, want := int(n), len(raw); got != want {
		t.Errorf("invalid sizes: got %d; want %d", got, want)
	}

	if got, want := out.Bytes(), raw; !reflect.DeepEqual(got, want) {
		t.Fatal("uncompressed data does not match original")
	}
}
