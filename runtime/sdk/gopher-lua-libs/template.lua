---@class templateModule
local templateModule = {}

---@param target string
---@return templateEngine
---@return err
function templateModule.choose(target)
    return {}, nil
end

---@class templateEngine
local templateEngine = {}

---@param str string
---@param value table
---@return string
function templateEngine:render(str, value)
    return ""
end
