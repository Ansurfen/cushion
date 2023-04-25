---@param tb table
---@return table | nil
function table.deepCopy(tb)
    if tb == nil then
        return nil
    end
    local copy = {}
    for k, v in pairs(tb) do
        if type(v) == 'table' then
            copy[k] = table.deepCopy(v)
        else
            copy[k] = v
        end
    end
    setmetatable(copy, table.deepCopy(getmetatable(tb)))
    return copy
end

---@param tbl table
---@param level number | nil
---@param filteDefault boolean | nil
function table.dump(tbl, level, filteDefault)
    filteDefault = filteDefault or true --default filter keywords（DeleteMe, _class_type）
    level = level or 1
    local indent_str = ""
    for i = 1, level do
        indent_str = indent_str .. "  "
    end

    print(indent_str .. "{")
    for k, v in pairs(tbl) do
        if filteDefault then
            if k ~= "_class_type" and k ~= "DeleteMe" then
                local item_str = string.format("%s%s = %s", indent_str .. " ", tostring(k), tostring(v))
                print(item_str)
                if type(v) == "table" then
                    table.dump(v, level + 1)
                end
            end
        else
            local item_str = string.format("%s%s = %s", indent_str .. " ", tostring(k), tostring(v))
            print(item_str)
            if type(v) == "table" then
                table.dump(v, level + 1)
            end
        end
    end
    print(indent_str .. "}")
end

---
--- filter keys of tbl according to `conditions`
---
---@param tbl table
---@param conditions table
---@return table
function table.filter(tbl, conditions)
    local matchCond = function(v, cond)
        for k, value in pairs(cond) do
            if v[k] ~= value then
                return false
            end
        end
        return true
    end
    local ret = {}
    for _, v in pairs(tbl) do
        if type(v) == "table" then
            if matchCond(v, conditions) == false then
                table.insert(ret, v)
            else
                table.filter(v, conditions)
            end
        end
    end
    return ret
end

---@param tb table
---@return table
function table.keys(tb)
    local keys = {}
    for key, _ in pairs(tb) do
        table.insert(keys, key)
    end
    return keys
end

---@param tbls table
---@return table
function table.Unpack(tbls)
    local dump = {}
    for _, tbl in ipairs(tbls) do
        for key, value in pairs(tbl) do
            dump[key] = value
        end
    end
    return dump
end

---@param tbA table
---@param tbB table
function table.merge(tbA, tbB)
    for _, value in pairs(tbB) do
        table.insert(tbA, value)
    end
end
