package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "the number to be parsed needs to be given as the first argument")
	}
	floatBits := 53
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "no amount of IEEE754 fraction bits specified, assuming 53 (double)")
	} else {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid argument", os.Args[2])
		}
		floatBits = i
	}
	calcRepresentationError(os.Args[1], floatBits)
}

func calcRepresentationError(decimal string, maxFractionalBits int) (string, error) {
	pieces := strings.SplitN(decimal, ".", 2)
	if len(pieces) == 1 {
		pieces = append(pieces, "0")
	}
	integralStr := pieces[0]
	integralPart, err := strconv.Atoi(integralStr)
	if err != nil {
		return "", err
	}
	fractionalStr := pieces[1]
	fmt.Println(strconv.Itoa(integralPart) + "." + fractionalStr + " to float")
	integralBits := fmt.Sprintf("%b", integralPart)
	integralBitLength := len(integralBits)
	fmt.Println("integral part takes " + strconv.Itoa(integralBitLength) + " bits")

	errorMagnitude := 0
	rest := ""
	fractionalBits := ""
	firstBitFound := integralPart > 0
	if firstBitFound {
		maxFractionalBits -= integralBitLength - 1
	}
	for rest = fractionalStr; errorMagnitude < maxFractionalBits; errorMagnitude++ {
		restLength := len(rest)
		restInt, err := strconv.Atoi(rest)
		if err != nil {
			return "", err
		}
		if restInt == 0 {
			break
		}
		double := strconv.Itoa(restInt * 2)
		desiredLength := len(rest) + 1
		if restLength > len(rest) {
			desiredLength = restLength
		}
		if desiredLength-len(double) > 0 {
			double = strings.Repeat("0", desiredLength-len(double)) + double
		}
		bit := string(double[0])
		rest = string(double[1:])
		fractionalBits += bit
	}
	fmt.Println(integralBits, ".", fractionalBits)
	fmt.Println("error 0.", rest, "*2^-", errorMagnitude)
	return "0." + rest + "*2^-" + strconv.Itoa(errorMagnitude), nil
}
