package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"sort"
)

func bootstrap() {
	fmt.Printf("App starting. Strict mode: %t, EV brk: %s\n", flagStrict(), flagExpBreak())
	glob := flagGlob()
	fmt.Printf("Analysing glob: %s\n", glob)

	// Load all files matching specified glob
	filePaths, err := getFileNames(glob)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("File list loaded, %d files found. Validating ...\n", len(filePaths))

	// Validate all loaded files
	var validFilePaths sort.StringSlice
	for _, b := range filePaths {
		if isValidFile(b) {
			validFilePaths = append(validFilePaths, b)
		}
	}

	// Sort files
	validFilePaths.Sort()

	fmt.Printf("Found %d valid files. Loading files into sequences ...\n", len(validFilePaths))

	// Load all EXIF data from files and put them into sequences
	sequences := loadSequences(validFilePaths, flagSequences())

	fmt.Printf("\nFound %d possible seqences, analyzing them now ...\n", len(sequences))

	// Run HDR validation for all sequences
	validSequences := validateSequences(sequences, flagStrict(), flagExpBreak())

	fmt.Printf("\nAnalysis complete. Found %d HDR sequences ...\n", len(validSequences))

	// We get results backwards, as procedure runs from end, reverse entire slice
	for i, j := 0, len(validSequences)-1; i < j; i, j = i+1, j-1 {
		validSequences[i], validSequences[j] = validSequences[j], validSequences[i]
	}

	// Print all found sequences
	for i, seq := range validSequences {
		fmt.Printf("%d. %s (%d)\n", i+1, seq[0].file.Name(), len(seq))
	}

	fmt.Println("Complete")
}

func loadSequences(filePaths []string, seqLengths []int) [][]*exifFile {
	var exifFiles []*exifFile
	var sequences [][]*exifFile
	for idx, filePath := range filePaths {
		xf, err := newExifFile(filePath)
		if err != nil {
			log.Println(err)
			continue
		}

		exifFiles = append(exifFiles, xf)

		for _, c := range seqLengths {
			if len(exifFiles) > c {
				seq := exifFiles[len(exifFiles)-c:]
				sequences = append(sequences, seq)
			}
		}

		fmt.Printf("\r%d/%d ...", idx+1, len(filePaths))
	}

	return sequences
}

func validateSequences(sequences [][]*exifFile, strict bool, brkDiff *big.Rat) [][]*exifFile {
	var validSequences [][]*exifFile
	for i := len(sequences) - 1; i >= 0; i-- {
		seq := sequences[i]
		isHdr, err := isHDRSequence(seq, strict, brkDiff)
		if err != nil {
			log.Println(err)
			continue
		}

		if isHdr {
			// Ok, its HDR, not put it into slice, or replace old shorter sequence
			if len(validSequences) > 0 {
				if sequenceContainsAnother(validSequences[len(validSequences)-1], seq) {
					continue
				}
			}

			validSequences = append(validSequences, seq)
		}

		fmt.Printf("\r%d/%d ...", len(sequences)-i, len(sequences))
	}

	return validSequences
}

func sequenceContainsAnother(haystack []*exifFile, needle []*exifFile) bool {
	for _, nxf := range needle {
		found := false
		for _, hxf := range haystack {
			if hxf == nxf {
				found = true
				continue
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func getFileNames(glob string) ([]string, error) {
	return filepath.Glob(glob)
}

func isValidFile(filePath string) bool {
	file, err := os.Stat(filePath)
	if err != nil || file.IsDir() {
		return false
	}

	return true
}
