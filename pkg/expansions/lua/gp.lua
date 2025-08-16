-- Expansion stubs and annotations

--- @meta

Gp = {}

-- Utils

---Converts a number to little endian 
---@param num integer
---@return table<integer>
Gp.toLittleEndian = function (num)
end

---Returns a number from little endian, must be 4 bytes long
---@param bytes table<integer>
---@return integer
Gp.fromLittleEndian = function (bytes)
end

-- Memory

---Set a memory address value
---@param addr integer
---@param val integer
---@return boolean success
Gp.setMemoryAddress = function (addr, val)
end

---Get a value from a memory address
---@param addr integer
---@return integer|nil
Gp.getMemoryAddress = function (addr)
end

---Set a range of memory
---@param addr integer
---@param values table<integer>
---@return boolean success
Gp.setMemoryRange = function (addr, values)
end

---Get a range of memory
---@param addr integer
---@param size integer
---@return table<integer>|nil
Gp.getMemoryRange = function (addr, size)
end

---Sets a range of memory to `value`
---@param addr integer
---@param size integer
---@param value integer
---@return boolean success
Gp.clearMemoryRange = function (addr, size, value)
end

-- Registers

---@param reg string
---@return integer|nil
Gp.getRegister = function (reg)
end

---@param reg string
---@param val integer
---@return boolean success
Gp.setRegister = function (reg, val)
end

---@param buf string
---@param val table<integer>
---@param offset integer
---@return boolean success
Gp.setBuffer = function (buf, val, offset)
end

---@param buf string
---@param offset integer
---@param length integer
---@return table<integer>|nil
Gp.getBuffer = function (buf, val, offset, length)
end

---Stop the VM. Internally sets `vm.Finished` to true.
Gp.stop = function ()
end

---@param addr integer
---@return string|nil
Gp.decodeInstructionToString = function (addr)
end