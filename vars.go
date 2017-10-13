package aoj

import (
	"log"
	"os"
)

var (
	// Endpoint of testcase
	BASE_ENDPOINT = "http://analytic.u-aizu.ac.jp:8080/aoj"

	// Cache directory
	CACHE_DIR = ".cache/aoj"

	// DEBUG Flag
	DEBUG = false

	// Acceptable error range of float
	Epsilon = 0.0001
)

func Debug(toggle bool) {
	DEBUG = toggle
}

func logdbg(fmt string, args ...interface{}) {
	if DEBUG {
		log.Printf("D "+fmt, args...)
	}
}

func logerr(fmt string, args ...interface{}) {
	log.Printf("E "+fmt, args...)
}

func logfatal(fmt string, args ...interface{}) {
	log.Printf("F "+fmt, args...)
	os.Exit(1)
}
