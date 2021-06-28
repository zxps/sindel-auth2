package main

import (
	"fmt"
	"strconv"
)

func main() {
	var i uint64 = 1 << 63
	fmt.Println(strconv.FormatUint(i, 2))
	fmt.Println(strconv.FormatUint(i, 10))

	var i2 int64 = 1 << 62
	fmt.Println(strconv.FormatInt(i2, 2))
	fmt.Println(strconv.FormatInt(i2, 10))
}

func pow(n uint64, p int) uint64 {
	var result uint64 = n
	for i := 0; i < p; i++ {
		result = result * (uint64)(n)
		fmt.Printf("result: %d\n", result)
	}

	return result
}
