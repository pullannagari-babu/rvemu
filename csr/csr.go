//-----------------------------------------------------------------------------
/*

RISC-V Control and Status Register Definitions

*/
//-----------------------------------------------------------------------------

package csr

import (
	"fmt"
	"strings"

	cli "github.com/deadsy/go-cli"
)

//-----------------------------------------------------------------------------

// Error is a CSR access error.
type Error struct {
	reg uint // csr register number
	n   uint // bitmap of csr access errors
}

// CSR error bits.
const (
	ErrTodo      = 1 << iota // csr not implemented
	ErrPrivilege             // insufficient privilege
	ErrReadOnly              // trying to write a read-only register
	ErrNoRead                // no read function (todo)
	ErrNoWrite               // no write function (todo)
)

func (e *Error) Error() string {
	s := []string{}
	if e.n&ErrTodo != 0 {
		s = append(s, "not implemented")
	}
	if e.n&ErrPrivilege != 0 {
		s = append(s, "insufficient privilege")
	}
	if e.n&ErrReadOnly != 0 {
		s = append(s, "read only")
	}
	if e.n&ErrNoRead != 0 {
		s = append(s, "no read function")
	}
	if e.n&ErrNoWrite != 0 {
		s = append(s, "no write function")
	}
	return fmt.Sprintf("%s: %s", Name(e.reg), strings.Join(s, ","))
}

//-----------------------------------------------------------------------------

// Privilege Levels.
const (
	PrivU = 0 // user
	PrivS = 1 // supervisor
	PrivM = 3 // machine
)

//-----------------------------------------------------------------------------

// Register numbers for specific CSRs.
const (
	FFLAGS  = 0x001
	FRM     = 0x002
	FCSR    = 0x003
	SSTATUS = 0x100
	MSTATUS = 0x300
)

//-----------------------------------------------------------------------------

// wrIgnore is a no-op write function.
func wrIgnore(s *State, val uint) {
}

// rdZero always reads the CSR as zero.
func rdZero(s *State) uint {
	return 0
}

//-----------------------------------------------------------------------------
// mcause register

// mcause interrupts
const (
	IntUserSoftware       = 0  // User software interrupt
	IntSupervisorSoftware = 1  // Supervisor software interrupt
	IntMachineSoftware    = 3  // Machine software interrupt
	IntUserTimer          = 4  // User timer interrupt
	IntSupervisorTimer    = 5  // Supervisor timer interrupt
	IntMachineTimer       = 7  // Machine timer interrupt
	IntUserExternal       = 8  // User external interrupt
	IntSupervisorExternal = 9  // Supervisor external interrupt
	IntMachineExternal    = 11 // Machine external interrupt
)

// SetInterrupt sets an interrupt code in mcause.
func (s *State) SetInterrupt(x uint) {
	s.mcause = (1 << (s.mxlen - 1)) | x
}

// mcause exceptions
const (
	ExInsAddrMisaligned         = 0  // Instruction address misaligned
	ExInsAccessFault            = 1  // Instruction access fault
	ExInsIllegal                = 2  // Illegal instruction
	ExBreakpoint                = 3  // Breakpoint
	ExLoadAddrMisaligned        = 4  // Load address misaligned
	ExLoadAccessFault           = 5  // Load access fault
	ExStoreAddrMisaligned       = 6  // Store/AMO address misaligned
	ExStoreAccessFault          = 7  // Store/AMO access fault
	ExEnvCallFromUserMode       = 8  // Environment call from U-mode
	ExEnvCallFromSupervisorMode = 9  // Environment call from S-mode
	ExEnvCallFromMachineMode    = 11 // Environment call from M-mode
	ExInsPageFault              = 12 // Instruction page fault
	ExLoadPageFault             = 13 // Load page fault
	ExStorePageFault            = 15 // Store/AMO page fault
)

// SetException sets an exception code in mcause.
func (s *State) SetException(x uint) {
	s.mcause = (0 << (s.mxlen - 1)) | x
}

func rdMCAUSE(s *State) uint {
	return s.mcause
}

//-----------------------------------------------------------------------------
// machine isa register

func mxl(xlen uint) uint {
	return map[uint]uint{32: 1, 64: 2, 128: 3}[xlen]
}

func initMISA(s *State) {
	s.misa = mxl(s.xlen) << (s.mxlen - 2)
}

func rdMISA(s *State) uint {
	return s.misa
}

//-----------------------------------------------------------------------------
// machine exception program counter

func wrMEPC(s *State, val uint) {
	s.mepc = val & ^uint(1)
}

func rdMEPC(s *State) uint {
	if s.ialign == 32 {
		return s.mepc & ^uint(3)
	}
	return s.mepc
}

//-----------------------------------------------------------------------------
// machine trap value register

func wrMTVAL(s *State, val uint) {
	s.mtval = val
}

func rdMTVAL(s *State) uint {
	return s.mtval
}

//-----------------------------------------------------------------------------
// machine trap vector

func wrMTVEC(s *State, val uint) {
	s.mtvec = val
}

func rdMTVEC(s *State) uint {
	return s.mtvec
}

//-----------------------------------------------------------------------------
// mscratch

func wrMSCRATCH(s *State, val uint) {
	s.mscratch = val
}

func rdMSCRATCH(s *State) uint {
	return s.mscratch
}

//-----------------------------------------------------------------------------
// sscratch

func wrSSCRATCH(s *State, val uint) {
	s.sscratch = val
}

func rdSSCRATCH(s *State) uint {
	return s.sscratch
}

//-----------------------------------------------------------------------------
// supervisor exception program counter

func wrSEPC(s *State, val uint) {
	s.sepc = val & ^uint(1)
}

func rdSEPC(s *State) uint {
	if s.ialign == 32 {
		return s.sepc & ^uint(3)
	}
	return s.sepc
}

//-----------------------------------------------------------------------------
// fcsr

const frmMask = uint(7 << 5)
const fflagsMask = uint(31 << 0)
const fcsrMask = frmMask | fflagsMask

func wrFCSR(s *State, val uint) {
	s.fcsr = val & fcsrMask
}

func rdFCSR(s *State) uint {
	return s.fcsr
}

func wrFRM(s *State, val uint) {
	s.fcsr &= ^frmMask
	s.fcsr |= (val & 7) << 5
}

func rdFRM(s *State) uint {
	return (s.fcsr & frmMask) >> 5
}

func wrFFLAGS(s *State, val uint) {
	s.fcsr &= ^fflagsMask
	s.fcsr |= val & fflagsMask
}

func rdFFLAGS(s *State) uint {
	return s.fcsr & fflagsMask
}

//-----------------------------------------------------------------------------
// mstatus

func wrMSTATUS(s *State, x uint) {
	s.mstatus = x
}

func rdMSTATUS(s *State) uint {
	return s.mstatus
}

func (s *State) mstatusRdMPP() uint {
	return (s.mstatus >> 11) & 3
}

func (s *State) mstatusWrMPP(x uint) {
	s.mstatus &= ^uint(3 << 11)
	s.mstatus |= (x & 3) << 11
}

func (s *State) mstatusWrMIE(x uint) {
	s.mstatus &= ^uint(1 << 3)
	s.mstatus |= (x & 1) << 3
}

func (s *State) mstatusRdMPIE() uint {
	return (s.mstatus >> 7) & 1
}

func (s *State) mstatusWrMPIE(x uint) {
	s.mstatus &= ^uint(1 << 7)
	s.mstatus |= (x & 1) << 7
}

//-----------------------------------------------------------------------------
// machine exception/interrupt delegation registers

func wrMEDELEG(s *State, x uint) {
	s.medeleg = x
}

func rdMEDELEG(s *State) uint {
	return s.medeleg
}

func wrMIDELEG(s *State, x uint) {
	s.mideleg = x
}

func rdMIDELEG(s *State) uint {
	return s.mideleg
}

//-----------------------------------------------------------------------------
// sstatus

func wrSSTATUS(s *State, x uint) {
	s.sstatus = x
}

func rdSSTATUS(s *State) uint {
	return s.sstatus
}

func (s *State) sstatusRdSPP() uint {
	return (s.sstatus >> 8) & 1
}

func (s *State) sstatusWrSPP(x uint) {
	s.sstatus &= ^uint(1 << 8)
	s.sstatus |= (x & 1) << 8
}

func (s *State) sstatusWrSIE(x uint) {
	s.sstatus &= ^uint(1 << 1)
	s.sstatus |= (x & 1) << 1
}

func (s *State) sstatusRdSPIE() uint {
	return (s.sstatus >> 5) & 1
}

func (s *State) sstatusWrSPIE(x uint) {
	s.sstatus &= ^uint(1 << 5)
	s.sstatus |= (x & 1) << 5
}

//-----------------------------------------------------------------------------

type wrFunc func(s *State, val uint)
type rdFunc func(s *State) uint

type csrDefn struct {
	name string // name of CSR
	wr   wrFunc // write function for CSR
	rd   rdFunc // read function for CSR
}

var lookup = map[uint]csrDefn{
	// User CSRs 0x000 - 0x0ff (read/write)
	0x000: {"ustatus", nil, nil},
	0x001: {"fflags", wrFFLAGS, rdFFLAGS},
	0x002: {"frm", wrFRM, rdFRM},
	0x003: {"fcsr", wrFCSR, rdFCSR},
	0x004: {"uie", nil, nil},
	0x005: {"utvec", nil, nil},
	0x040: {"uscratch", nil, nil},
	0x041: {"uepc", nil, nil},
	0x042: {"ucause", nil, nil},
	0x043: {"utval", nil, nil},
	0x044: {"uip", nil, nil},
	// User CSRs 0xc00 - 0xc7f (read only)
	0xc00: {"cycle", nil, nil},
	0xc01: {"time", nil, nil},
	0xc02: {"instret", nil, nil},
	0xc03: {"hpmcounter3", nil, nil},
	0xc04: {"hpmcounter4", nil, nil},
	0xc05: {"hpmcounter5", nil, nil},
	0xc06: {"hpmcounter6", nil, nil},
	0xc07: {"hpmcounter7", nil, nil},
	0xc08: {"hpmcounter8", nil, nil},
	0xc09: {"hpmcounter9", nil, nil},
	0xc0a: {"hpmcounter10", nil, nil},
	0xc0b: {"hpmcounter11", nil, nil},
	0xc0c: {"hpmcounter12", nil, nil},
	0xc0d: {"hpmcounter13", nil, nil},
	0xc0e: {"hpmcounter14", nil, nil},
	0xc0f: {"hpmcounter15", nil, nil},
	0xc10: {"hpmcounter16", nil, nil},
	0xc11: {"hpmcounter17", nil, nil},
	0xc12: {"hpmcounter18", nil, nil},
	0xc13: {"hpmcounter19", nil, nil},
	0xc14: {"hpmcounter20", nil, nil},
	0xc15: {"hpmcounter21", nil, nil},
	0xc16: {"hpmcounter22", nil, nil},
	0xc17: {"hpmcounter23", nil, nil},
	0xc18: {"hpmcounter24", nil, nil},
	0xc19: {"hpmcounter25", nil, nil},
	0xc1a: {"hpmcounter26", nil, nil},
	0xc1b: {"hpmcounter27", nil, nil},
	0xc1c: {"hpmcounter28", nil, nil},
	0xc1d: {"hpmcounter29", nil, nil},
	0xc1e: {"hpmcounter30", nil, nil},
	0xc1f: {"hpmcounter31", nil, nil},
	// User CSRs 0xc80 - 0xcbf (read only)
	0xc80: {"cycleh", nil, nil},
	0xc81: {"timeh", nil, nil},
	0xc82: {"instreth", nil, nil},
	0xc83: {"hpmcounter3h", nil, nil},
	0xc84: {"hpmcounter4h", nil, nil},
	0xc85: {"hpmcounter5h", nil, nil},
	0xc86: {"hpmcounter6h", nil, nil},
	0xc87: {"hpmcounter7h", nil, nil},
	0xc88: {"hpmcounter8h", nil, nil},
	0xc89: {"hpmcounter9h", nil, nil},
	0xc8a: {"hpmcounter10h", nil, nil},
	0xc8b: {"hpmcounter11h", nil, nil},
	0xc8c: {"hpmcounter12h", nil, nil},
	0xc8d: {"hpmcounter13h", nil, nil},
	0xc8e: {"hpmcounter14h", nil, nil},
	0xc8f: {"hpmcounter15h", nil, nil},
	0xc90: {"hpmcounter16h", nil, nil},
	0xc91: {"hpmcounter17h", nil, nil},
	0xc92: {"hpmcounter18h", nil, nil},
	0xc93: {"hpmcounter19h", nil, nil},
	0xc94: {"hpmcounter20h", nil, nil},
	0xc95: {"hpmcounter21h", nil, nil},
	0xc96: {"hpmcounter22h", nil, nil},
	0xc97: {"hpmcounter23h", nil, nil},
	0xc98: {"hpmcounter24h", nil, nil},
	0xc99: {"hpmcounter25h", nil, nil},
	0xc9a: {"hpmcounter26h", nil, nil},
	0xc9b: {"hpmcounter27h", nil, nil},
	0xc9c: {"hpmcounter28h", nil, nil},
	0xc9d: {"hpmcounter29h", nil, nil},
	0xc9e: {"hpmcounter30h", nil, nil},
	0xc9f: {"hpmcounter31h", nil, nil},
	// Supervisor CSRs 0x100 - 0x1ff (read/write)
	0x100: {"sstatus", wrSSTATUS, rdSSTATUS},
	0x102: {"sedeleg", wrIgnore, rdZero},
	0x103: {"sideleg", wrIgnore, rdZero},
	0x104: {"sie", wrIgnore, rdZero},
	0x105: {"stvec", wrIgnore, rdZero},
	0x106: {"scounteren", nil, nil},
	0x140: {"sscratch", wrSSCRATCH, rdSSCRATCH},
	0x141: {"sepc", wrSEPC, rdSEPC},
	0x142: {"scause", nil, nil},
	0x143: {"stval", nil, nil},
	0x144: {"sip", wrIgnore, rdZero},
	0x180: {"satp", wrIgnore, nil},
	// Machine CSRs 0xf00 - 0xf7f (read only)
	0xf11: {"mvendorid", nil, rdZero},
	0xf12: {"marchid", nil, rdZero},
	0xf13: {"mimpid", nil, rdZero},
	0xf14: {"mhartid", nil, rdZero},
	// Machine CSRs 0x300 - 0x3ff (read/write)
	0x300: {"mstatus", wrMSTATUS, rdMSTATUS},
	0x301: {"misa", wrIgnore, rdMISA},
	0x302: {"medeleg", wrMEDELEG, rdMEDELEG},
	0x303: {"mideleg", wrMIDELEG, rdMIDELEG},
	0x304: {"mie", wrIgnore, nil},
	0x305: {"mtvec", wrMTVEC, rdMTVEC},
	0x306: {"mcounteren", nil, nil},
	0x320: {"mucounteren", nil, nil},
	0x321: {"mscounteren", nil, nil},
	0x322: {"mhcounteren", nil, nil},
	0x323: {"mhpmevent3", nil, nil},
	0x324: {"mhpmevent4", nil, nil},
	0x325: {"mhpmevent5", nil, nil},
	0x326: {"mhpmevent6", nil, nil},
	0x327: {"mhpmevent7", nil, nil},
	0x328: {"mhpmevent8", nil, nil},
	0x329: {"mhpmevent9", nil, nil},
	0x32a: {"mhpmevent10", nil, nil},
	0x32b: {"mhpmevent11", nil, nil},
	0x32c: {"mhpmevent12", nil, nil},
	0x32d: {"mhpmevent13", nil, nil},
	0x32e: {"mhpmevent14", nil, nil},
	0x32f: {"mhpmevent15", nil, nil},
	0x330: {"mhpmevent16", nil, nil},
	0x331: {"mhpmevent17", nil, nil},
	0x332: {"mhpmevent18", nil, nil},
	0x333: {"mhpmevent19", nil, nil},
	0x334: {"mhpmevent20", nil, nil},
	0x335: {"mhpmevent21", nil, nil},
	0x336: {"mhpmevent22", nil, nil},
	0x337: {"mhpmevent23", nil, nil},
	0x338: {"mhpmevent24", nil, nil},
	0x339: {"mhpmevent25", nil, nil},
	0x33a: {"mhpmevent26", nil, nil},
	0x33b: {"mhpmevent27", nil, nil},
	0x33c: {"mhpmevent28", nil, nil},
	0x33d: {"mhpmevent29", nil, nil},
	0x33e: {"mhpmevent30", nil, nil},
	0x33f: {"mhpmevent31", nil, nil},
	0x340: {"mscratch", wrMSCRATCH, rdMSCRATCH},
	0x341: {"mepc", wrMEPC, rdMEPC},
	0x342: {"mcause", nil, rdMCAUSE},
	0x343: {"mtval", wrMTVAL, rdMTVAL},
	0x344: {"mip", nil, nil},
	0x380: {"mbase", nil, nil},
	0x381: {"mbound", nil, nil},
	0x382: {"mibase", nil, nil},
	0x383: {"mibound", nil, nil},
	0x384: {"mdbase", nil, nil},
	0x385: {"mdbound", nil, nil},
	0x3a0: {"pmpcfg0", wrIgnore, nil},
	0x3a1: {"pmpcfg1", wrIgnore, nil},
	0x3a2: {"pmpcfg2", wrIgnore, nil},
	0x3a3: {"pmpcfg3", wrIgnore, nil},
	0x3b0: {"pmpaddr0", wrIgnore, nil},
	0x3b1: {"pmpaddr1", wrIgnore, nil},
	0x3b2: {"pmpaddr2", wrIgnore, nil},
	0x3b3: {"pmpaddr3", wrIgnore, nil},
	0x3b4: {"pmpaddr4", wrIgnore, nil},
	0x3b5: {"pmpaddr5", wrIgnore, nil},
	0x3b6: {"pmpaddr6", wrIgnore, nil},
	0x3b7: {"pmpaddr7", wrIgnore, nil},
	0x3b8: {"pmpaddr8", wrIgnore, nil},
	0x3b9: {"pmpaddr9", wrIgnore, nil},
	0x3ba: {"pmpaddr10", wrIgnore, nil},
	0x3bb: {"pmpaddr11", wrIgnore, nil},
	0x3bc: {"pmpaddr12", wrIgnore, nil},
	0x3bd: {"pmpaddr13", wrIgnore, nil},
	0x3be: {"pmpaddr14", wrIgnore, nil},
	0x3bf: {"pmpaddr15", wrIgnore, nil},
	// Machine CSRs 0xb00 - 0xb7f (read/write)
	0xb00: {"mcycle", nil, nil},
	0xb02: {"minstret", nil, nil},
	0xb03: {"mhpmcounter3", nil, nil},
	0xb04: {"mhpmcounter4", nil, nil},
	0xb05: {"mhpmcounter5", nil, nil},
	0xb06: {"mhpmcounter6", nil, nil},
	0xb07: {"mhpmcounter7", nil, nil},
	0xb08: {"mhpmcounter8", nil, nil},
	0xb09: {"mhpmcounter9", nil, nil},
	0xb0a: {"mhpmcounter10", nil, nil},
	0xb0b: {"mhpmcounter11", nil, nil},
	0xb0c: {"mhpmcounter12", nil, nil},
	0xb0d: {"mhpmcounter13", nil, nil},
	0xb0e: {"mhpmcounter14", nil, nil},
	0xb0f: {"mhpmcounter15", nil, nil},
	0xb10: {"mhpmcounter16", nil, nil},
	0xb11: {"mhpmcounter17", nil, nil},
	0xb12: {"mhpmcounter18", nil, nil},
	0xb13: {"mhpmcounter19", nil, nil},
	0xb14: {"mhpmcounter20", nil, nil},
	0xb15: {"mhpmcounter21", nil, nil},
	0xb16: {"mhpmcounter22", nil, nil},
	0xb17: {"mhpmcounter23", nil, nil},
	0xb18: {"mhpmcounter24", nil, nil},
	0xb19: {"mhpmcounter25", nil, nil},
	0xb1a: {"mhpmcounter26", nil, nil},
	0xb1b: {"mhpmcounter27", nil, nil},
	0xb1c: {"mhpmcounter28", nil, nil},
	0xb1d: {"mhpmcounter29", nil, nil},
	0xb1e: {"mhpmcounter30", nil, nil},
	0xb1f: {"mhpmcounter31", nil, nil},
	// Machine CSRs 0xb80 - 0xbbf (read/write)
	0xb80: {"mcycleh", nil, nil},
	0xb82: {"minstreth", nil, nil},
	0xb83: {"mhpmcounter3h", nil, nil},
	0xb84: {"mhpmcounter4h", nil, nil},
	0xb85: {"mhpmcounter5h", nil, nil},
	0xb86: {"mhpmcounter6h", nil, nil},
	0xb87: {"mhpmcounter7h", nil, nil},
	0xb88: {"mhpmcounter8h", nil, nil},
	0xb89: {"mhpmcounter9h", nil, nil},
	0xb8a: {"mhpmcounter10h", nil, nil},
	0xb8b: {"mhpmcounter11h", nil, nil},
	0xb8c: {"mhpmcounter12h", nil, nil},
	0xb8d: {"mhpmcounter13h", nil, nil},
	0xb8e: {"mhpmcounter14h", nil, nil},
	0xb8f: {"mhpmcounter15h", nil, nil},
	0xb90: {"mhpmcounter16h", nil, nil},
	0xb91: {"mhpmcounter17h", nil, nil},
	0xb92: {"mhpmcounter18h", nil, nil},
	0xb93: {"mhpmcounter19h", nil, nil},
	0xb94: {"mhpmcounter20h", nil, nil},
	0xb95: {"mhpmcounter21h", nil, nil},
	0xb96: {"mhpmcounter22h", nil, nil},
	0xb97: {"mhpmcounter23h", nil, nil},
	0xb98: {"mhpmcounter24h", nil, nil},
	0xb99: {"mhpmcounter25h", nil, nil},
	0xb9a: {"mhpmcounter26h", nil, nil},
	0xb9b: {"mhpmcounter27h", nil, nil},
	0xb9c: {"mhpmcounter28h", nil, nil},
	0xb9d: {"mhpmcounter29h", nil, nil},
	0xb9e: {"mhpmcounter30h", nil, nil},
	0xb9f: {"mhpmcounter31h", nil, nil},
	// Machine Debug CSRs 0x7a0 - 0x7af (read/write)
	0x7a0: {"tselect", nil, nil},
	0x7a1: {"tdata1", nil, nil},
	0x7a2: {"tdata2", nil, nil},
	0x7a3: {"tdata3", nil, nil},
	// Machine Debug Mode Only CSRs 0x7b0 - 0x7bf (read/write)
	0x7b0: {"dcsr", nil, nil},
	0x7b1: {"dpc", nil, nil},
	0x7b2: {"dscratch", nil, nil},
	// Hypervisor CSRs 0x200 - 0x2ff (read/write)
	0x200: {"hstatus", nil, nil},
	0x202: {"hedeleg", nil, nil},
	0x203: {"hideleg", nil, nil},
	0x204: {"hie", nil, nil},
	0x205: {"htvec", nil, nil},
	0x240: {"hscratch", nil, nil},
	0x241: {"hepc", nil, nil},
	0x242: {"hcause", nil, nil},
	0x243: {"hbadaddr", nil, nil},
	0x244: {"hip", nil, nil},
}

// Name returns the name of a given CSR.
func Name(reg uint) string {
	if x, ok := lookup[reg]; ok {
		return x.name
	}
	return fmt.Sprintf("0x%03x", reg)
}

// access returns the access string of a given CSR.
func access(reg uint) string {
	mode := [4]string{"u", "s", "h", "m"}[(reg>>8)&3]
	rw := [4]string{"rw", "rw", "rw", "r_"}[(reg>>10)&3]
	return mode + rw
}

//-----------------------------------------------------------------------------

// canAccess returns true if the register can be accessed at the current privilege level.
func (s *State) canAccess(reg uint) bool {
	priv := (reg >> 8) & 3
	return s.Priv >= priv
}

// canWr returns true if the register can be written.
func canWr(reg uint) bool {
	rw := (reg >> 10) & 3
	return rw != 3
}

//-----------------------------------------------------------------------------

// State stores the CSR state for the CPU.
type State struct {
	Priv     uint // current privilege level
	xlen     uint // cpu register length 32/64/128
	mxlen    uint // machine register length
	uxlen    uint // user register length
	sxlen    uint // supervisor register length
	ialign   uint // instruction alignment 16/32
	fcsr     uint // floating point control and status register
	mepc     uint // machine exception program counter
	mtvec    uint // machine trap vector base address register
	mtval    uint // machine trap value register
	mstatus  uint // machine status
	mscratch uint // machine scratch
	misa     uint // machine isa register
	mcause   uint // machine cause register
	medeleg  uint // machine exception delegation register
	mideleg  uint // machine interrupt delegation register
	sscratch uint // supervisor scratch
	sepc     uint // supervisor exception program counter
	sstatus  uint // supervisor status
}

// NewState returns a CSR state object.
func NewState(xlen uint) *State {
	s := &State{
		Priv:   PrivM, // start at machine level
		xlen:   xlen,
		mxlen:  xlen,
		uxlen:  xlen,
		sxlen:  xlen,
		ialign: 16, // TODO
	}
	initMISA(s)
	return s
}

// Rd reads from a CSR.
func (s *State) Rd(reg uint) (uint, error) {
	if !s.canAccess(reg) {
		return 0, &Error{reg, ErrPrivilege}
	}
	if x, ok := lookup[reg]; ok {
		if x.rd == nil {
			return 0, &Error{reg, ErrNoRead}
		}
		return x.rd(s), nil
	}
	return 0, &Error{reg, ErrTodo}
}

// Wr writes to a CSR.
func (s *State) Wr(reg, val uint) error {
	if !canWr(reg) {
		return &Error{reg, ErrReadOnly}
	}
	if !s.canAccess(reg) {
		return &Error{reg, ErrPrivilege}
	}
	if x, ok := lookup[reg]; ok {
		if x.wr == nil {
			return &Error{reg, ErrNoWrite}
		}
		x.wr(s, val)
		return nil
	}
	return &Error{reg, ErrTodo}
}

// Set sets bits in a CSR.
func (s *State) Set(reg, bits uint) error {
	val, err := s.Rd(reg)
	if err != nil {
		return err
	}
	return s.Wr(reg, val|bits)
}

// Clr clears bits in a CSR.
func (s *State) Clr(reg, bits uint) error {
	val, err := s.Rd(reg)
	if err != nil {
		return err
	}
	return s.Wr(reg, val & ^bits)
}

// Display displays the CSR state.
func (s *State) Display() string {
	// allow all reads
	savedPriv := s.Priv
	s.Priv = PrivM
	// read all registers
	x := [][]string{}
	for reg := uint(0); reg < 4096; reg++ {
		val, err := s.Rd(reg)
		if err != nil {
			e := err.(*Error)
			if e.n == ErrTodo || e.n == ErrNoRead {
				continue
			}
		}
		regStr := fmt.Sprintf("%03x %s %s", reg, access(reg), Name(reg))
		valStr := "0"
		if val != 0 {
			valStr = fmt.Sprintf("%08x", val)
		}
		x = append(x, []string{regStr, valStr})
	}
	// restore privilege
	s.Priv = savedPriv
	// return the table string
	return cli.TableString(x, []int{0, 0}, 1)
}

//-----------------------------------------------------------------------------

// MRET performs an MRET operation.
func (s *State) MRET() (uint, error) {
	if !s.canAccess(MSTATUS) {
		return 0, &Error{MSTATUS, ErrPrivilege}
	}
	s.Priv = s.mstatusRdMPP()
	s.mstatusWrMIE(s.mstatusRdMPIE())
	s.mstatusWrMPIE(1)
	s.mstatusWrMPP(PrivU)
	return rdMEPC(s), nil
}

// SRET performs an SRET operation.
func (s *State) SRET() (uint, error) {
	if !s.canAccess(SSTATUS) {
		return 0, &Error{SSTATUS, ErrPrivilege}
	}
	s.Priv = s.sstatusRdSPP()
	s.sstatusWrSIE(s.sstatusRdSPIE())
	s.sstatusWrSPIE(1)
	s.sstatusWrSPP(0)
	return rdSEPC(s), nil
}

//-----------------------------------------------------------------------------
