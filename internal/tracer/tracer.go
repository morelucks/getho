package tracer

// Tracer provides execution tracing capabilities
type Tracer interface {
	// Trace generates an opcode-level execution trace
	Trace(txHash string) (*Trace, error)
}

// Trace represents an execution trace
type Trace struct {
}
