package utils

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func HexToBigInt(hex string) (*big.Int, error) {
	if hex == "0x0" {
		return big.NewInt(0), nil
	}
	bignum, ok := new(big.Int).SetString(strings.Replace(hex, "0x", "", -1), 16)
	if !ok {
		return nil, fmt.Errorf("failed to conver int")
	}
	return bignum, nil
}

func HexToInt(hex string) (int64, error) {
	i, err := strconv.ParseInt(strings.Replace(hex, "0x", "", -1), 16, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func IntToHex(i int) string {
	return fmt.Sprintf("0x%x", i)
}
