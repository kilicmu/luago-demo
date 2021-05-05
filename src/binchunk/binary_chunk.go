package binchunk

type binaryChunk struct {
	header, // 头部
	sizeUpvalues byte // 主函数 binaryChunk 数量
	mainFunc *Prototype // 主函数原型
}

type header struct { // about 30 bytes
	signature       [4]byte // 魔数 chuck签名
	version         byte    // 版本号 大版本 * 16 + 小版本 => 5.4.3 => 5 * 16 + 4 => 0x54
	format          byte    // 格式号：检验虚拟机匹配
	luaData         [6]byte // LuaC_DATA: 进一步校验 0x 19 93 0D 0A 1993 + 回车（0x0D） + 换行（0x0A）
	cintSize        byte    // cint 长度
	csizetSize      byte    // size_t 长度
	instructionSize byte    //
	luaIntegerSize  byte    // lua int 长度
	luaNumberSize   byte    // lua number长度
	luacInt         int64   // 存储 lua 整数值 0x5678 用于定位主机内存的大小端 0x78 56 => 小端 0x56 78 => 大端
	luacNum         float64 // 存储 lua 的浮点数 370.5 检测二进制chunk使用的浮点格式
}

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSZIET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type Upvalue struct {
	Instack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

type Prototype struct {
	Source          string // 源文件名 debug  这个字段非主函数无值
	LineDefined     uint32 //
	LastLineDefined uint32
	NumberParams    byte     // 固定参数个数
	IsVararg        byte     // 是否有变长参数
	MaxStackSize    byte     // 记录执行期间用到的寄存器数量
	Code            []uint32 // 指令表 （每条指令4 byte
	// 常量表(nil boolean int float string) tag(1byte)
	/*
		0x00 => nil 不存储
		0x01 => boolean byte
		0x03 => number  LUA_NUMBER_SIZE
		0x13 => integer LUA_INTEGER_SIZE
		0x04 => string 短字符串
		0x14 => string 长字符串
	*/
	Constants    []interface{}
	Upvalues     []Upvalue // upvalue表
	Protos       []*Prototype
	LineInfo     []uint32 // 行号表
	LocVars      []LocVar // 局部变量表
	UpvalueNames []string // Upvalue名表 与 Upvalue表中元素一一对应记录每个Upvalue在源码的名字
}

/*
resolve binary lua chunk
*/
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte()
	return reader.readProto("")
}
