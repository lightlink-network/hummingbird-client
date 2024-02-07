package cannon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log/slog"
	"os"
	"strconv"

	uc "github.com/CryptoKass/unicorn/bindings/go/unicorn"
	"github.com/pellartech/mipsevm"
)

type MipsEvmOpts struct {
	Logger      *slog.Logger
	BaseDir     string // Directory to read inputs, write outputs, and cache preimage oracle data. Defaults '/tmp/cannon'
	ProgramPath string // Path to binary file containing the program to run. Default '/assets/minigeth.bin'
}

type MipsEvm struct {
	Opts    *MipsEvmOpts
	regfalt int
}

func NewMipsEvm(opts *MipsEvmOpts) *MipsEvm {

	regfault := -1
	regfault_str, regfault_valid := os.LookupEnv("REGFAULT")
	if regfault_valid {
		regfault, _ = strconv.Atoi(regfault_str)
	}

	return &MipsEvm{Opts: opts, regfalt: regfault}
}

type MipsProgram struct {
	root     string
	regfault int
	target   int
	ram      map[uint32](uint32)
	unicorn  uc.Unicorn
	lastStep int
}

func (m *MipsEvm) prepare(blockNum int64, target int) (*MipsProgram, error) {
	var root string
	var regfault = m.regfalt

	if blockNum >= 0 {
		root = fmt.Sprintf("%s/%d_%d", m.Opts.BaseDir, 0, blockNum)
	} else {
		root = m.Opts.BaseDir
	}

	m.Opts.Logger.Info("Preparing program", "blockNum", blockNum, "target", target, "root", root)
	ram := make(map[uint32](uint32))
	lastStep := 1

	// GetHookedUnicorn returns a unicorn instance with hooks for regfault and target
	mu := mipsevm.GetHookedUnicorn(root, ram, func(step int, mu uc.Unicorn, ram map[uint32](uint32)) {
		if step == regfault {
			m.Opts.Logger.Info("regfault at step", "step", step)
			mu.RegWrite(uc.MIPS_REG_V0, 0xbabababa)
		}
		if step == target {
			mipsevm.SyncRegs(mu, ram)
			fn := fmt.Sprintf("%s/checkpoint_%d.json", root, step)
			mipsevm.WriteCheckpoint(ram, fn, step)
			if step == target {
				// done
				mu.RegWrite(uc.MIPS_REG_PC, 0x5ead0004)
			}
		}
		lastStep = step + 1
	})

	// reset the RAM registers
	mipsevm.ZeroRegisters(ram)

	// Load the program
	mipsevm.LoadMappedFileUnicorn(mu, m.Opts.ProgramPath, ram, 0)

	return &MipsProgram{
		root:     root,
		regfault: regfault,
		target:   target,
		ram:      ram,
		unicorn:  mu,
		lastStep: lastStep,
	}, nil
}

// WriteGolden writes the golden snapshot of the program to the base directory.
// Golden snapshots refers to a snapshot of the program state before any execution
// has taken place.
func (m *MipsEvm) WriteGolden(blockNum int64, target int) error {
	program, err := m.prepare(blockNum, target)
	if err != nil {
		return err
	}

	m.Opts.Logger.Info("Writing golden snapshot", "blockNum", blockNum, "target", target, "root", program.root)
	mipsevm.WriteCheckpoint(program.ram, fmt.Sprintf("%s/golden.json", program.root), -1)

	m.Opts.Logger.Info("Golden snapshot written", "file", fmt.Sprintf("%s/golden.json", program.root))
	return nil
}

func (m *MipsEvm) WriteCheckpoints(blockNum int64, target int) error {
	program, err := m.prepare(blockNum, target)
	if err != nil {
		return err
	}

	m.Opts.Logger.Info("Running program", "blockNum", blockNum, "target", target, "root", program.root)
	mipsevm.LoadMappedFileUnicorn(program.unicorn, fmt.Sprintf("%s/input", program.root), program.ram, 0x30000000)
	err = program.unicorn.Start(0, 0x5ead0004)
	if err != nil {
		return err
	}

	mipsevm.SyncRegs(program.unicorn, program.ram)
	if program.target == -1 {
		if program.ram[0x30000800] != 0x1337f00d {
			return fmt.Errorf("Failed to output stateroot")
		}

		output_filename := fmt.Sprintf("%s/output", program.root)
		outputs, err := ioutil.ReadFile(output_filename)
		if err != nil {
			return fmt.Errorf("Failed to read output file %w", err)
		}
		real := append([]byte{0x13, 0x37, 0xf0, 0x0d}, outputs...)

		output := []byte{}
		for i := 0; i < 0x44; i += 4 {
			t := make([]byte, 4)
			binary.BigEndian.PutUint32(t, program.ram[uint32(0x30000800+i)])
			output = append(output, t...)
		}

		if bytes.Compare(real, output) != 0 {
			return fmt.Errorf("Output mismatch, overwriting")
		}
		m.Opts.Logger.Debug("Output matches")

		mipsevm.WriteCheckpoint(program.ram, fmt.Sprintf("%s/checkpoint_final.json", program.root), program.lastStep)
		m.Opts.Logger.Info("Final checkpoint written", "file", fmt.Sprintf("%s/checkpoint_final.json", program.root), "steps", program.lastStep)
	}

	return nil
}
