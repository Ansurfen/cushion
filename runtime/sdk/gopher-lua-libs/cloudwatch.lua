---@class cloudwatchModule
local cloudwatchModule = {}

---@return cloudwatchClient
---@return err
function cloudwatchModule.new()
    return {}, nil
end

---@class cloudwatchClient
local cloudwatchClient = {}

---@param query table
---@return table
---@return err
function cloudwatchClient:get_metric_data(query)
    return {}, nil
end

---@vararg any
---@return err
function cloudwatchClient:download(...)
    return nil
end
