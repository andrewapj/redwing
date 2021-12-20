package redwing

import "log"

func PrintLog(s string, r *Options) {
	if r.logging {
		log.Printf("Redwing: " + s)
	}
}
