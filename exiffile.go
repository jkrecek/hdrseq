package main

import (
	"github.com/rwcarlsen/goexif/exif"
	"math/big"
	"os"
	"time"
)

type exifFile struct {
	file *os.File

	ExposureTime      *big.Rat
	ExposureBiasValue *big.Rat
	DateTime          time.Time
}

func newExifFile(filePath string) (*exifFile, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	et, err := getRat(x, exif.ExposureTime)
	if err != nil {
		return nil, err
	}

	ebv, err := getRat(x, exif.ExposureBiasValue)
	if err != nil {
		return nil, err
	}

	dt, err := x.DateTime()
	if err != nil {
		return nil, err
	}

	ef := &exifFile{
		file:              f,
		ExposureTime:      et,
		ExposureBiasValue: ebv,
		DateTime:          dt,
	}

	return ef, nil
}

func getRat(x *exif.Exif, name exif.FieldName) (*big.Rat, error) {
	ebv, err := x.Get(name)
	if err != nil {
		return nil, err
	}

	return ebv.Rat(0)
}
