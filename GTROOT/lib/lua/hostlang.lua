
-- http://lua-users.org/wiki/AutomagicTables
-- https://www.lua.org/pil/13.4.1.html
-- https://www.lua.org/pil/contents.html

HostLang = {}      -- The numeric namespace

HostVarRef = {}
HostVarRef.prototype = {
    _parent = nil,
    _tag = nil,
    _suffix = '<suffix>',
    _value = nil,
    type = "numeric"
}

HostVarRef.mt = {}   -- metatable for numeric

function HostLang._variable(tag) -- create a new variable
    local v = { _tag = nil, _suffix = tag }
    v.type = "numeric"
    setmetatable(v, HostVarRef.mt)
    return v
end

function HostLang.numeric(tag) -- create a new numeric variable
    local n = HostLang._variable(tag)
    n.type = "numeric"
    return n
end

function HostLang.pair(tag) -- create a new pair variable
    local n = HostLang._variable(tag)
    pr.type = "pair"
    return n
end

function HostVarRef.mt.__index(n, suffix)
    local atsuffix = '@'..suffix
    local suffix_var = rawget(n, atsuffix)
    if suffix_var then
        return suffix_var
    end
    return HostVarRef.prototype[suffix]
end

function HostVarRef.mt.__tostring(n)
    return "<numeric "..n:fullname().."="..(n._value or "<unknown>")..">"
end

function HostVarRef.mt.__newindex(n, suffix, v)
    suffix = suffix or 'value'
    local atsuffix = '@' .. suffix
    local var = HostLang._variable(suffix)
    rawset(var, "_tag", n._tag or n)
    rawset(var, "_parent", n)
    rawset(var, '_value', v)
    rawset(n, atsuffix, var)
    return n
end

HostVarRef.prototype.tag = function(n)
    if not n._tag then
        return n
    end
    return n._tag
end

HostVarRef.prototype.isknown = function(n)
    if n._value == nil then
        return false
    end
    return true
end

HostVarRef.prototype.value = function(n, x)
    if x then
        n._value = x
    else
        return n._value
    end
end

HostVarRef.prototype.fullname = function(n)
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

-- --- predefined functions ---------------------------------------------

vardef = {}

function vardef.z(suffixes)
    trace(0, "z() with suffixes "..(suffixes or ""))
    x = varref.refer_to("x"..(suffixes or ""))
    y = varref.refer_to("y"..(suffixes or ""))
    p = pair.new{x, y}
    return p
end

return HostLang
