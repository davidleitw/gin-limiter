package limiter

const Script = `
	local globalKey = KEYS[1]
	local singleKey = KEYS[2]
	local times = tonumber(ARGV[1])
	local globalLimit = tonumber(ARGV[2])
	local singleLimit = tonumber(ARGV[3])
	
	local IpGlobalInfo = redis.call('HGETALL', globalKey)
	local IpSingleInfo = redis.call('HGETALL', singleKey)
`
