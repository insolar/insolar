package instracer

import (
	"os"
)

func hostname() (h string) {
	h, _ = os.Hostname()
	return
}
