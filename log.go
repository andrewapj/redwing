package redwing

import "log"

func PrintLog(s string, r *Options) {
	if r.Logging {
		log.Printf("Redwing: " + s)
	}
}
