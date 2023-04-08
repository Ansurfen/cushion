---@class CusionStrings
local cusionStrings

---@param s string
---@param sep string
---@return table
function cusionStrings.Split(s, sep)
    return {}
end

---@param s string
---@param sep string
---@return string, string, boolean
function cusionStrings.Cut(s, sep)
    return "", "", false
end

---@param s string
---@param substr string
---@return boolean
function cusionStrings.Contains(s, substr)
    return false
end

---@param s string
---@param prefix string
---@return boolean
function cusionStrings.HasPrefix(s, prefix)
    return false
end

---@param s string
---@param suffix string
---@return boolean
function cusionStrings.HasSuffix(s, suffix)
    return false
end
