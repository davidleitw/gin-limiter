package limiter

const Script = `
	local result = {}
		
	local globalKey = KEYS[1] 
	local singleKey = KEYS[2]
	
	local globalLimit = tonumber(ARGV[1])
	local singleLimit = tonumber(ARGV[2])

	local IpGlobalInfo = redis.call('HGETALL', globalKey) 
	local IpSingleInfo = redis.call('HGETALL', singleKey)

	local globalExpired = ARGV[3] -- if true, IpGlobalInfo = "1", else = "0"
	local singleExpired = ARGV[4]

	if #IpGlobalInfo == 0 or globalExpired == "1" then 
		redis.call('HMSET', globalKey, "Count", 1)
		redis.call('HMSET', singleKey, "Count", 1)
		result[1] = globalLimit - 1
		result[2] = singleLimit - 1
		return result
	end

	if #IpSingleInfo == 0 or singleExpired == "1" then 
		if tonumber(IpGlobalInfo[2]) < globalLimit then 
			result[1] = globalLimit - redis.call('HINCRBY', globalKey, "Count", 1)
		end
		redis.call('HMSET', singleKey, "Count", 1)
		result[2] = singleLimit - 1
		return result
	end

	local gc = tonumber(IpGlobalInfo[2]) -- global count 
	local sc = tonumber(IpSingleInfo[2]) -- single count 

	if gc < globalLimit then 
		result[1] = globalLimit - redis.call('HINCRBY', globalKey, "Count", 1)
	else 
		result[1] = -1
	end

	if sc < singleLimit then 
		result[2] = singleLimit - redis.call('HINCRBY', singleKey, "Count", 1)
	else
		result[2] = -1
	end

	return result
`

const TestScript = `
	local result = {}
	local test = ARGV[1]
	local test2 = "TestHash"

	-- local c = redis.call('HMSET', test2, "Count", 1)
	local t = redis.call('HINCRBY', test2, "Count", 1) 

	result[1] = redis.call('HINCRBY', test2, "Count", 1)
	result[2] = 4 - redis.call('HINCRBY', test2, "Count", 1)
	return result
`
