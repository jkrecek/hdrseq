package main

import (
	"math"
	"math/big"
	"time"
)

func isHDRSequence(xfs []*exifFile, strict bool, brkDiff *big.Rat) (bool, error) {
	if len(xfs) < 1 {
		panic("invalid exif count")
	}

	expTimeCheck, err := hdrCheckExpTime(xfs)
	if err != nil || !expTimeCheck {
		return expTimeCheck, err
	}

	expBiasCheck, err := hdrExpBiasCheck(xfs, strict, brkDiff)
	return expBiasCheck, err
}

func hdrCheckExpTime(xfs []*exifFile) (bool, error) {
	totalDuration := &big.Rat{}
	for _, xf := range xfs {
		totalDuration.Add(totalDuration, xf.ExposureTime)
	}

	totalExposure, _ := totalDuration.Float64()
	tolerance := math.Ceil(totalExposure + (float64(len(xfs)) * 0.2))

	sd := xfs[0].DateTime
	ed := xfs[len(xfs)-1].DateTime

	diff := ed.Sub(sd)

	inTolerance := diff <= (time.Duration(tolerance) * time.Second)

	return inTolerance, nil
}

func hdrExpBiasCheck(xfs []*exifFile, strict bool, brkDiff *big.Rat) (bool, error) {
	br := &big.Rat{}
	initExp, _ := xfs[0].ExposureBiasValue.Float64()
	invExp, _ := br.Inv(brkDiff).Float64()
	baseExp := br.Mul(big.NewRat(int64(math.Round(initExp*invExp)), int64(1)), brkDiff)

	if strict {
		for i, xf := range xfs {
			r := &big.Rat{}
			expectedEB := r.Add(baseExp, big.NewRat(int64(math.Pow(-1, float64(i))*float64(brkDiff.Num().Int64())*(math.Ceil(0.5*float64(i)))), brkDiff.Denom().Int64()))

			ebvf, _ := xf.ExposureBiasValue.Float64()
			eebf, _ := expectedEB.Float64()
			if math.Round(ebvf*float64(10)) != math.Round(eebf*float64(10)) {
				return false, nil
			}

		}

		return true, nil
	} else {
		calcBias := &big.Rat{}

		for i, xf := range xfs {
			for y, xfa := range xfs {
				// Check that sequence does not contain same exposure
				if i != y && xfa.ExposureBiasValue.Cmp(xf.ExposureBiasValue) == 0 {
					return false, nil
				}
			}

			calcBias = calcBias.Add(calcBias, xf.ExposureBiasValue)
		}

		calcBias = calcBias.Quo(calcBias, big.NewRat(int64(len(xfs)), 1))
		return calcBias.Cmp(baseExp) == 0, nil

	}
}
