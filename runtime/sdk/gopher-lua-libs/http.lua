---@class httpModule
local httpModule = {}

---@param opt table | nil
---@return httpClient
function httpModule.client(opt)
    return {}
end

---@param method string
---@param url string
---@vararg string
---@return httpRequest
function httpModule.request(method, url, ...)
    return {}
end

---@param opt string | table
---@return httpServer
---@return err
function httpModule.server(opt)
    return {}, nil
end

---@class httpRequest
local httpRequest = {}

---@vararg string
function httpRequest:set_basic_auth(...)

end

---@param key string
---@param value string
function httpRequest:header_set(key, value)

end

---@class httpResponse
---@field code number
---@field body string
local httpResponse = {}

---@class httpClient
local httpClient = {}

---@param req httpRequest
---@return httpResponse
---@return err
function httpClient:do_request(req)
    return {}, nil
end

---@class httpServer
local httpServer = {}

---@return httpServerRequest
---@return httpServerResponse
function httpServer:accept()
    return {}, {}
end

---@param str string
function httpServer:do_handle_string(str)

end

---@param handle fun(res :httpResponse,req :httpRequest):nil
function httpServer:do_handle_function(handle)

end

---@return string
function httpServer:addr()
    return ""
end

---@class httpServerRequest
---@field host string
---@field method string
---@field referer string
---@field proto string
---@field path string
---@field raw_path string
---@field raw_query string
---@field request_uri string
---@field request_addr string
---@field user_agent string
---@field headers table
---@field query table
local httpServerRequest = {}

---@return string
---@return err
function httpServerRequest.body()
    return "", nil
end

---@class httpServerResponse
local httpServerResponse = {}

---@param code number
function httpServerResponse:code(code)

end

---@param key string
---@param value string
function httpServerResponse:header(key, value)

end

---@param payload string
function httpServerResponse:write(payload)

end

function httpServerResponse:done()

end
