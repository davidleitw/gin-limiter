package limiter

const ResetScript = `	
	local routeKey = KEYS[1]
	local staticKey = KEYS[2]
	local routeDeadline = ARGV[1]

	redis.call('HSET', staticKey, "Count", 1)
	redis.call('HSET', routeKey, "Count", 1, "Deadline", routeDeadline)
`

const Script = `
	local result = {}
	local routeKey = KEYS[1]
	local staticKey = KEYS[2]

	local routeLimit = ARGV[1]
	local staticLimit = ARGV[2]
	local routeDeadline = ARGV[3]

	local routeInfo = redis.call('HGETALL', routeKey)
	local staticCount = redis.call('HGET', staticKey)

	-- First time visit
	if (not staticCount) then 
		redis.call('HSET', staticKey, "Count", 1)
		redis.call('HSET', routeKey, "Count", 1, "Deadline", routeDeadline)
		result[1] = staticLimit - 1
		result[2] = routeLimit - 1
		return result
	end 

	if #routeInfo == 0 then 
		if tonumber(staticCount) < staticLimit then
			result[1] = staticLimit - redis.call('HINCRBY', staticKey, "Count", 1)
		else 
			result[1] = -1
		end
		redis.call('HSET', routeKey, "Count", 1, "Deadline", routeDeadline)
		result[2] = routeLimit - 1
		return result
	end

	local rCount = tonumber(routeInfo[2])
	local sCount = tonumber(staticCount)

	if sCount < staticLimit then 
		result[1] = staticLimit - redis.call('HINCRBY', staticKey, "Count", 1)
	else 
		result[1] = -1
	end

	if rCount < routeLimit then 
		result[2] = redis.call('HINCRBY', routeKey, "Count", 1)
	else 
		result[2] = -1
	end

	return result
`

const TestScript = `
	local result = {}
	local test = redis.call('HGETALL', "test")
	local t = redis.call('HGET', "test1", "val1")
	result[1] = test
	result[2] = type(t)

	return result
`
