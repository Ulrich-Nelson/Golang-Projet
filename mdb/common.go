package mdb

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ",")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func GenerateID() string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func GenerateTimestamping() time.Time {
	return time.Now().UTC()
}

func FormatTruncDate(t time.Time) string {
	return t.Format("02-Jan-2006")
}

func ParseTruncDate(s string) time.Time {
	t, err := time.Parse("02-Jan-2006", s)
	if err != nil {
		panic(err)
	}
	return t
}

func ParseUnixDate(s string) time.Time {
	// Convert string as Unix timestamp to time.Time
	date64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(date64, 0)
}
