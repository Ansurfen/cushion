---@class CusionPath
local cusionPath

---@vararg string
---@return string
function cusionPath.Join(...)
    return ""
end

---@param path string
---@return boolean
function cusionPath.IsAbs(path)
    return false
end

---@param path string
---@return string
function cusionPath.Ext(path)
    return ""
end

---@param path string
---@return string
function cusionPath.Clean(path)
    return ""
end

---@param path string
---@return string
function cusionPath.Dir(path)
    return ""
end

---@param path string
---@return string, string
function cusionPath.Split(path)
    return "", ""
end

---@param pattern string
---@param name string
---@return boolean, error
function cusionPath.Match(pattern, name)
    return false, nil
end

---@param path string
---@return string
function cusionPath.Base(path)
    return ""
end

---@param fullpath string
---@return string
function cusionPath.Filename(fullpath)
    return ""
end