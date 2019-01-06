package memtable

import (
	"github.com/fabulousduck/smol/errors"
)

/*
MemTable is a simple collection of memory regions in use
*/
type MemTable map[string]*MemRegion

/*
MemRegion represents a region of memory in the IR
*/
type MemRegion struct {
	Addr, Size, Value int
}

/*
Put places a variable on a memoryTable

notes

chip-8's blocks are 8 bit, so 1 byte.
with a total of 4096 bytes
*/
func (table *MemTable) Put(name string, value int) *MemRegion {
	region := new(MemRegion)
	//check if there is any memory left for our variable
	currentMemSize := table.getSize()
	if currentMemSize >= 95 {
		errors.OutOfMemoryError()
	}

	region.Addr = table.findNextEmptyAddr()
	region.Size = 2
	region.Value = value

	(*table)[name] = region
	return region
}

func (table *MemTable) lookupVariable(name string) *MemRegion {
	if val, ok := (*table)[name]; ok {
		return val
	}
	errors.UndefinedVariableError(name)
	return nil
}

func (table MemTable) findNextEmptyAddr() int {
	//Note: all ints are padded to by 2 bytes long

	varAddrSpaceStart := 0xEA0
	varAddrSpaceEnd := 0xEFF

	currentSpaceUsed := 0

	for i := 0; i < len(table); i++ {
		currentSpaceUsed += 0x2
	}

	if varAddrSpaceStart+currentSpaceUsed+0x2 > varAddrSpaceEnd {
		errors.OutOfMemoryError()
	}

	return varAddrSpaceStart + currentSpaceUsed
}

func (table MemTable) getSize() int {
	return len(table)
}