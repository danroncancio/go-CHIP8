package chip8

type Display struct {
	WindowWidth int
	WindowHeight int
	ResolutionWidth int
	ResolutionHeight int

	RGBABuffer   [4 * 64 * 32]byte
	BinaryBuffer [64 * 32]byte
}
