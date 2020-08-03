package limiter

const Script = `
	local result = {}
		
	local now = tonumber(ARGV[1]) -- 現在時間
	
	local globalKey = KEYS[1] 
	local singleKey = KEYS[2]
	
	local globalLimit = tonumber(ARGV[2])
	local singleLimit = tonumber(ARGV[3])
	local globalExpired = ARGV[4]
	local singleExpired = ARGV[5]

	local IpGlobalInfo = redis.call('HGETALL', globalKey)
	local IpSingleInfo = redis.call('HGETALL', singleKey)
	
	if globalExpired == "true" then 
		redis.call('HMSET', globalKey, "Count", 1)
		result[1] = globalLimit - 1
	end

	if singleExpired == "true" then
		redis.call('HMSET', singleKey, "Count", 1)
		result[2] = singleLimit - 1
	end

	-- 該Ip第一次訪問
	if #IpGlobalInfo == 0 then
		redis.call('HMSET', globalKey, "Count", 1)
		redis.call('HMSET', singleKey, "Count", 1)
		result[1] = globalLimit - 1 
		result[2] = singleLimit - 1 
		return result
	end

	if #IpSingleInfo == 0 then 
		redis.call('HMSET', singleKey, "Count", 1)
		result[2] = singleLimit - 1
	end

	local gc = tonumber(IpGlobalInfo[2])
	local sc = tonumber(IpSingleInfo[2])


`

const TestScript = `
	local result = {}

	-- 測試註解
	result[5] = 28
	result[1] = 10
	result[2] = 20 + 2
	result[3] = KEYS[1]
	result[4] = ARGV[1]
	
	return result
`
