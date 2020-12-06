package main

import (
	"flag"
	"log"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

var glob = flag.String("glob", "", "File mask used to search files (required)")
var strictMode = flag.Bool("strict", true, "Type of EV measuring")
var expBreak = flag.String("expbrk", "1/3", "Expected EV diff")
var sequences = flag.String("seq", "3,5,9", "Tested length of sequences")

func parseFlags() error {
	flag.Parse()
	testFlags()

	return nil
}

func testFlags() {
	if flagGlob() == "" {
		log.Fatalln("Parameter `glob` is required")
	}

	_, err := flagExpBrkOptional()
	if err != nil {
		log.Fatalln("Parameter `expbrk` is invalid. Use format `1/3`")
	}

	_, err = flagSequencesOptional()
	if err != nil {
		log.Fatalln("Parameter `seq` is invalid. Use format `3,5,9`")
	}
}

func flagGlob() string {
	return *glob
}

func flagStrict() bool {
	return *strictMode
}

func flagExpBreak() *big.Rat {
	eb, err := flagExpBrkOptional()
	if err != nil {
		panic(err)
	}

	return eb
}

func flagSequences() []int {
	seqs, err := flagSequencesOptional()
	if err != nil {
		panic(err)
	}

	return seqs
}

func flagExpBrkOptional() (r *big.Rat, err error) {
	parts := strings.SplitN(*expBreak, "/", 2)
	nom, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}

	denom, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	r = big.NewRat(int64(nom), int64(denom))
	return
}

func flagSequencesOptional() ([]int, error) {
	parts := strings.Split(*sequences, ",")
	var seqs []int
	for _, part := range parts {
		sn, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}

		seqs = append(seqs, sn)
	}

	sort.Ints(seqs)

	return seqs, nil
}
