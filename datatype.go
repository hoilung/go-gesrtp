package gogesrtp

const (
	R  = 0x08 // Register (Word)
	AI = 0x0a // Analog Input (Word)
	AQ = 0x0c // Analog Output (Word)
	/* Byte selector - address */
	I_BYTE  = 0x10 // Descrete Inputs (Byte)
	Q_BYTE  = 0x12 // Descrete Outputs (Byte)
	T_BYTE  = 0x14 // Descrete Temporary Bits (Byte)
	M_BYTE  = 0x16 // Descrete Markers (Byte)
	SA_BYTE = 0x18 // System Bits A-part (Byte)
	SB_BYTE = 0x1a // 0x20,   // System Bits B-part (Byte)
	SC_BYTE = 0x1c // 0x22,   // System Bits C-part (Byte)
	G_BYTE  = 0x38 // Genius Global (Byte)
	/* Binary selector - address */
	I_BIT  = 0x46 // Descrete Input (Bit)
	Q_BIT  = 0x48 // Descrete Output (Bit)
	T_BIT  = 0x4a // Descrete Temporary (Bit)
	M_BIT  = 0x4c // Descrete Marker (Bit)
	SA_BIT = 0x4e // System Bit A-part (Bit)
	SB_BIT = 0x50 // System Bit B-part (Bit)
	SC_BIT = 0x52 // System Bit C-part (Bit)
	G_BIT  = 0x56 // Genius Global (Bit)
)
