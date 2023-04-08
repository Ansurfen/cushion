---@class base64Module
---@field RawStdEncoding RawStdEncoding
---@field StdEncoding StdEncoding
---@field RawURLEncoding RawURLEncoding
---@field URLEncoding URLEncoding
local base64Module = {}

---@class RawStdEncoding
local RawStdEncoding = {}

---@class StdEncoding
local StdEncoding = {}

---@class RawURLEncoding
local RawURLEncoding = {}

---@class URLEncoding
local URLEncoding = {}

---@param str string
---@return string
function RawStdEncoding:encode_to_string(str)
    return ""
end

---@param str string
---@return string
function RawStdEncoding:decode_string(str)
    return ""
end

---@param str string
---@return string
function StdEncoding:encode_to_string(str)
    return ""
end

---@param str string
---@return string
function StdEncoding:decode_string(str)
    return ""
end

---@param str string
---@return string
function RawURLEncoding:encode_to_string(str)
    return ""
end

---@param str string
---@return string
function RawURLEncoding:decode_string(str)
    return ""
end

---@param str string
---@return string
function URLEncoding:encode_to_string(str)
    return ""
end

---@param str string
---@return string
function URLEncoding:decode_string(str)
    return ""
end
