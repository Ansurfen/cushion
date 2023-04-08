---@class stringsModule
local stringsModule = {}

---@param str string
---@param prefix string
---@return table
function stringsModule.split(str, prefix)
    return {}
end

---@param str string
---@param prefix string
---@return boolean
function stringsModule.has_prefix(str, prefix)
    return false
end

---@param str string
---@param suffix string
---@return boolean
function stringsModule.has_suffix(str, suffix)
    return false
end

---@param str string
---@param cutset string
---@return string
function stringsModule.trim(str, cutset)
    return ""
end

---@param str string
---@param substr string
---@return boolean
function stringsModule.contains(str, substr)
    return false
end

---@return stringReader
function stringsModule.new_reader()
    return {}
end

---@return stringBuilder
function stringsModule.new_builder()
    return {}
end

---@class stringReader
local stringReader = {}

---@param str string
function stringReader:read(str)

end

---@class stringBuilder
local stringBuilder = {}

---@vararg any
function stringBuilder:write(...)

end
