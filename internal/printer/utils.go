package printer

import (
	"runtime"
	"strconv"
	"strings"
)

func getGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	s := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.ParseUint(s, 10, 64)
	return id
}
