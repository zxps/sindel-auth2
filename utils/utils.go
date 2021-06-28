package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// IsFileExists - check if file exists
func IsFileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func StorageTimestamp() string {
	t := time.Now()

	return t.Format("2006-01-02 15:04:05")
}

func ConvertDatetimeToTimestamp(datetime string) uint32 {
	t, _ := time.Parse("2006-01-02 15:04:05", datetime)
	return uint32(t.Unix())
}

func ConvertTimestampToDatetime(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

func FormatBytesAsText(b uint64) string {
	value := float64(b) / 1024 / 1024
	result := fmt.Sprintf("%f", value)
	return result
}

func ConvertIpToString(ip uint64) string {
	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatUint((ip>>24)&0xff, 10)
	b1 := strconv.FormatUint((ip>>16)&0xff, 10)
	b2 := strconv.FormatUint((ip>>8)&0xff, 10)
	b3 := strconv.FormatUint(ip&0xff, 10)

	return b0 + "." + b1 + "." + b2 + "." + b3
}

func SimpleRandomToken(length int) string {
	rand.Seed(time.Now().Unix())
	var output strings.Builder
	charSet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteRune(randomChar)
	}

	return output.String()
}