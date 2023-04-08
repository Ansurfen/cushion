---@class pluginModule
local pluginModule = {}

---@param code string
function pluginModule.do_string(code)

end

---@param filename string
function pluginModule.do_file(filename)

end

---@param code string
---@param payload string
function pluginModule.do_string_with_payload(code, payload)

end

---@param filename string
---@param payload string
function pluginModule.do_file_with_payload(filename, payload)

end

---@class plugin
local plugin = {}

function plugin:run()

end

function plugin:stop()

end

---@return err
function plugin:is_running()
    return nil
end
