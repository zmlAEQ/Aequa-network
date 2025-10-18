package payload

type StandardBlock struct{}
func (StandardBlock) Kind() string { return "standard_block" }
