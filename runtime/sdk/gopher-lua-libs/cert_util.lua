---@class cert_utilModule
local cert_utilModule = {}

---@param hostname string
---@param ip string
---@return number
---@return err
function cert_utilModule:not_after(hostname, ip)
    return 0, nil
end
