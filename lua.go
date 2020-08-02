package limiter

const Script = `
	local result = {}
		
	local now = tonumber(ARGV[1])
	
	local globalKey = KEYS[1]
	local singleKey = KEYS[2]
	
	local globalLimit = tonumber(ARGV[2])
	local singleLimit = tonumber(ARGV[3])

	local IpGlobalInfo = redis.call('HGETALL', globalKey)
	local IpSingleInfo = redis.call('HGETALL', singleKey)
	
	-- 該Ip第一次訪問
	if #IpGlobalInfo == 0 then
		return 0
	end
`

const TestScript = `
	local result = {}

	-- 測試註解
	result[1] = 10
	result[2] = 20 + 2
	result[3] = KEYS[1]
	result[4] = ARGV[1]
	
	return result
`
