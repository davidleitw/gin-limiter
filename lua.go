package limiter

const ResetScript = `	
	local routeKey = KEYS[1]
	local staticKey = KEYS[2]
	local routeDeadline = ARGV[1]

	redis.call('HMSET', staticKey, "Count", 1)
	redis.call('HMSET', routeKey, "Count", 1, "Deadline", routeDeadline)
`

const Script = `

`
