
-- http://lua-users.org/wiki/AutomagicTables
-- https://www.lua.org/pil/13.4.1.html
-- https://www.lua.org/pil/contents.html

Numeric = {}      -- The numeric namespace
Numeric.prototype = { _parent = nil, _tag = nil, _suffix = '<suffix>', _value = nil }
Numeric.mt = {}   -- metatable for numeric

function Numeric.variable(tag) -- create a new numeric variable
    local n = { _tag = nil, _suffix = tag }
    setmetatable(n, Numeric.mt)
    return n
end

function Numeric.mt.__index(n, suffix)
    local atsuffix = '@'..suffix
    local suffix_var = rawget(n, atsuffix)
    if suffix_var then
        return suffix_var
    end
    return Numeric.prototype[suffix]
end

function Numeric.mt.__tostring(n)
    return "<numeric "..n:fullname().."="..(n._value or "<unknown>")..">"
end

function Numeric.mt.__newindex(n, suffix, v)
    suffix = suffix or 'value'
    local atsuffix = '@' .. suffix
    local var = Numeric.variable(suffix)
    rawset(var, "_tag", n._tag or n)
    rawset(var, "_parent", n)
    rawset(var, '_value', v)
    rawset(n, atsuffix, var)
    return n
end

Numeric.prototype.tag = function(n)
    if not n._tag then
        return n
    end
    return n._tag
end

Numeric.prototype.isknown = function(n)
    if n._value == nil then
        return false
    end
    return true
end

Numeric.prototype.value = function(n, x)
    if x then
        n._value = x
    else
        return n._value
    end
end

Numeric.prototype.fullname = function(n)
    s = ""
    repeat
        if type(n._suffix) == "number" then
            s = '['..n._suffix..']' .. s
        else
            s = n._suffix .. s
            if not (n._parent == nil) then
                s = "."..s
            end
        end
        n = n._parent
    until not n
    return s
end


return Numeric
