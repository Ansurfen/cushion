---@class argparseModule
---@vararg string
---@return argparser
local argparseModule = function(...)
    return {}
end

---@class argparser
local argparser = {}

---@vararg string
function argparser:argument(...)
end

---@vararg string
function argparser:option(...)
end

function argparser:parse()
end
