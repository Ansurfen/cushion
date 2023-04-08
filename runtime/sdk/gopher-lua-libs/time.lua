---@class timeModule
local timeModule = {}

---@return number
function timeModule.unix()
    return 0
end

---@param value string
---@param layout string
---@return number
---@return err
function timeModule.parse(value, layout)
    return 0, nil
end

---@param value number
---@param layout string
---@param location string
---@return string
---@return err
function timeModule.format(value, layout, location)
    return "", nil
end
