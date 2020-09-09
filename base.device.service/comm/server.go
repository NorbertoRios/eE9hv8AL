package comm

import "expvar"

var (
	//ServerCounters composes service metrics
	ServerCounters = expvar.NewMap("Server")
)
