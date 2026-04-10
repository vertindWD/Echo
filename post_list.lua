-- ==========================================
-- File: post_list.lua
-- 针对 Echo 社区帖子列表接口的压测脚本
-- ==========================================

-- 1. 配置请求方法
wrk.method = "GET"

-- 2. 配置请求头 (重点：填入你的真实 Token)
local token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozNDc0MDM5MTM4MzE0MjQwMCwidXNlcm5hbWUiOiJ6aGFvIiwiaXNzIjoiZWNobyIsImV4cCI6MTc3NTc0MTc5Nn0.hFuKuy0xt1aiXShulRe7d6ovUsrZccxa-dgo-CuJZek"
wrk.headers["Authorization"] = "Bearer " .. token
wrk.headers["Accept"] = "application/json"

-- 3. 配置请求路径和 Query 参数
-- 这里测试按热度(score)排序的第一页
wrk.path = "/api/v1/post/list?page=1&size=10&order=score"

-- (可选) 每次请求发送前的钩子函数
-- request 函数会返回最终的 HTTP 请求字符串
request = function()
    return wrk.format(wrk.method, wrk.path, wrk.headers, nil)
end