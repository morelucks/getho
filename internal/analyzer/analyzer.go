package analyzer

// Analyzer provides gas and fee analysis capabilities
type Analyzer interface {
	// AnalyzeGas analyzes gas usage and fees
	AnalyzeGas(txHash string) (*GasAnalysis, error)
}

// GasAnalysis represents gas and fee breakdown
type GasAnalysis struct {
}
