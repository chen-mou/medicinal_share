package resource

import "os"

var Machine string

var Mode int

const UriPre = "http://localhost:15888"

func init() {
	Machine, _ = os.Hostname()
}
