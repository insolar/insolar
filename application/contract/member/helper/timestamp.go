package helper

import (
	"fmt"
	"strconv"
	"time"
)

func ParseTimestamp(timeStr string) (time.Time, error) {

	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Unix(0, 0), fmt.Errorf("failed to parse time: %s", err.Error())
	}
	return time.Unix(i, 0), nil
}
