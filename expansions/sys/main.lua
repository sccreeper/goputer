local attributeNames = {
    "name",
    "display_width",
    "display_height",
    "memory",
    "expansions"
}

local attributeIds = {
    name = 0,
    displayWidth = 1,
    displayHeight = 2,
    memory = 3,
}

--- @type table<table<integer>>
local attributeValues = {
    ["name"] = {string.byte("goputer.sys", 1, -1)},
    ["display_width"] = Gp.toLittleEndian(320),
    ["display_height"] = Gp.toLittleEndian(240),
    ["memory"] = Gp.toLittleEndian((2 ^ 16) + (320 * 240 * 3)),
    ["expansions"] = {string.byte("goputer.sys\x00\x00\x00", 1, -1)}
}

Gp.hooks.addHook("vm_finish", "finish", function ()
    Gp.log("VM Finished")
end)

---Is v in a
---@param a table
---@param v any
---@return boolean
local function arrContainsValue(a, v)

    for _, value in pairs(a) do
        if value == v then
            return true
        end
    end

    return false

end

---Main handler for data
---@param data table<integer>
---@return table<integer>
function Handler(data)

    if data[2] == attributeIds.name then
        return attributeValues["name"]
    elseif data[2] == attributeIds.displayWidth then
        return attributeValues["display_width"]
    elseif data[2] == attributeIds.displayHeight then
        return attributeValues["display_height"]
    elseif data[2] == attributeIds.memory then
        return attributeValues["memory"]
    end

    return {0}

end

---Sets an attribute
---@param name string
---@param data table<integer>
function SetAttribute(name, data)

    assert(type(name), "string")

    if not arrContainsValue(attributeNames, name) then
        print(string.format("no attribute called %s", name))
        return
    end

    attributeValues[name] = data

end

---Gets an attribute
---@param name string
---@return any
function GetAttribute(name)

    assert(type(name), "string")

    if not arrContainsValue(attributeNames, name) then
        return {}
    end

    return attributeValues[name]
end
