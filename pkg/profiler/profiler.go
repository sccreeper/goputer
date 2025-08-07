package profiler

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"sccreeper/goputer/pkg/compiler"
	"sccreeper/goputer/pkg/constants"
	"sccreeper/goputer/pkg/util"
	"sccreeper/goputer/pkg/vm"
	"time"
)

const magicString string = "GPPR"

type ProfileEntry struct {
	TotalCycleTime     uint64
	TotalTimesExecuted uint64

	Address     uint32
	Instruction [5]byte

	StandardDeviation float64
	sumOfSquares      float64
}

type Profiler struct {
	vm          *vm.VM
	TotalCycles uint64
	Data        map[uint64]*ProfileEntry

	cycleAddress uint32
	cycleStart   time.Time
	cycleLength  uint64
}

func NewProfiler(machine *vm.VM) (*Profiler, error) {

	if machine == nil {
		return nil, errors.New("vm cannot be nil")
	}

	return &Profiler{
		vm:   machine,
		Data: make(map[uint64]*ProfileEntry),
	}, nil

}

func (p *Profiler) Finish() {

	for _, v := range p.Data {

		if v.TotalTimesExecuted == 1 {
			continue
		}

		variance := (v.sumOfSquares / float64(v.TotalTimesExecuted)) - math.Pow(float64(v.TotalCycleTime)/float64(v.TotalTimesExecuted), 2)

		if math.Signbit(variance) {
			variance = 0.0
		}

		v.StandardDeviation = math.Sqrt(variance)

	}

}

func (p *Profiler) Dump(w io.WriteSeeker) (int, error) {

	var numBytes int

	_, err := w.Seek(0, io.SeekEnd)
	if err != nil {
		return numBytes, err
	}

	// Header data

	var headerBytes []byte = make([]byte, 0)

	headerBytes = append(headerBytes, []byte(magicString)...)

	numEntries := len(p.Data)
	headerBytes = binary.LittleEndian.AppendUint64(headerBytes, uint64(numEntries))

	headerBytes = binary.LittleEndian.AppendUint64(headerBytes, uint64(p.TotalCycles))

	n, err := w.Write(headerBytes)
	if err != nil {
		return numBytes, err
	}
	numBytes += n

	// Entries

	for _, v := range p.Data {

		dataToWrite := make([]byte, 5)

		copy(dataToWrite[:5], v.Instruction[:])

		dataToWrite = binary.LittleEndian.AppendUint32(dataToWrite, v.Address)
		dataToWrite = binary.LittleEndian.AppendUint64(dataToWrite, v.TotalCycleTime)
		dataToWrite = binary.LittleEndian.AppendUint64(dataToWrite, v.TotalTimesExecuted)
		dataToWrite = binary.LittleEndian.AppendUint64(dataToWrite, math.Float64bits(v.StandardDeviation))

		_, err = w.Seek(0, io.SeekEnd)
		if err != nil {
			return numBytes, err
		}

		n, err = w.Write(dataToWrite)
		if err != nil {
			return numBytes, err
		}

		numBytes += n

	}

	return numBytes, nil

}

func (p *Profiler) Load(r io.ReadSeeker) (int, error) {

	p.Data = make(map[uint64]*ProfileEntry)

	var totalBytesRead int

	_, err := r.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	magicStringBytes := make([]byte, 4)

	n, err := r.Read(magicStringBytes)
	if err != nil {
		return totalBytesRead, err
	}
	totalBytesRead += n

	if string(magicStringBytes) != magicString {
		return totalBytesRead, errors.New("invalid header")
	}

	headerBytes := make([]byte, 16)

	n, err = r.Read(headerBytes)
	if err != nil {
		return totalBytesRead, err
	}
	totalBytesRead += n

	numEntries := binary.LittleEndian.Uint64(headerBytes[:8])
	p.TotalCycles = binary.LittleEndian.Uint64(headerBytes[8:])

	for i := 0; i < int(numEntries); i++ {

		var entryBytes [33]byte = [33]byte{}

		n, err = r.Read(entryBytes[:])
		if err == io.EOF {
			fmt.Println("EOF!")
			break
		} else if err != nil {
			return totalBytesRead, err
		}
		totalBytesRead += n

		instructionBytes := entryBytes[:5]
		addr := binary.LittleEndian.Uint32(entryBytes[5:9])

		p.Data[genKey(instructionBytes, addr)] = &ProfileEntry{
			Address:     addr,
			Instruction: [5]byte(instructionBytes),

			TotalCycleTime:     binary.LittleEndian.Uint64(entryBytes[9:17]),
			TotalTimesExecuted: binary.LittleEndian.Uint64(entryBytes[17:25]),
			StandardDeviation:  math.Float64frombits(binary.LittleEndian.Uint64(entryBytes[25:33])),
		}

	}

	return totalBytesRead, nil

}

func (p *Profiler) StartCycle() {

	p.cycleStart = time.Now()
	p.cycleAddress = p.vm.Registers[constants.RProgramCounter]

}

func (p *Profiler) Cycle() {

	p.StartCycle()
	p.vm.Cycle()
	p.EndCycle()

}

func (p *Profiler) EndCycle() {

	p.cycleLength = uint64(time.Since(p.cycleStart).Nanoseconds())

	p.TotalCycles++

	key := genKey(p.vm.CurrentInstruction, p.cycleAddress)

	if _, exists := p.Data[key]; exists {

		p.Data[key].TotalCycleTime += p.cycleLength
		p.Data[key].TotalTimesExecuted++
		p.Data[key].sumOfSquares += math.Pow(float64(p.cycleLength), 2)

	} else {

		p.Data[key] = &ProfileEntry{
			TotalCycleTime:     p.cycleLength,
			TotalTimesExecuted: 1,

			Address:     p.cycleAddress,
			Instruction: [5]byte(p.vm.CurrentInstruction),
		}

	}

}

func genKey(instruction []byte, programCounter uint32) (key uint64) {

	if len(instruction) != int(compiler.InstructionLength) {
		panic("invalid instruction length")
	}

	key |= uint64(instruction[0]) << 56
	key |= uint64(instruction[1]) << 48
	key |= uint64(instruction[2]) << 40
	key |= uint64(instruction[3]) << 32
	key |= uint64(instruction[4]) << 24

	key |= uint64(util.Clamp(programCounter, 0, 0xFFFFFF))

	return

}
