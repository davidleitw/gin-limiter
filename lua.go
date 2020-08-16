package limiter

const ResetScript = `	
	local routeKey = KEYS[1]
	local staticKey = KEYS[2]

	redis.call('HMSET', staticKey, "Count", 1)
	redis.call('HMSET', routeKey, "Count", 1)

`

const Script = `

`
