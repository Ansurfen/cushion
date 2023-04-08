---@class regexpModule
local regexpModule = {}

---@param regexp string
---@param data string
---@return boolean
---@return err
function regexpModule.match(regexp, data)
    return false, nil
end

---@param regexp string
---@param data string
---@return table
---@return err
function regexpModule.find_all_string_submatch(regexp, data)
    return {}, nil
end

---@param regexp string
---@return regexp
---@return err
function regexpModule.compile(regexp)
    return {}, nil
end

---@class regexp
local regexp = {}

---@param data string
---@return boolean
function regexp:match(data)
    return false
end

---@param data string
---@return table
function regexp:find_all_string_submatch(data)
    return {}
end
