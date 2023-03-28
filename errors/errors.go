package errors

type DappFileError struct {
}

func (err *DappFileError) Error() string {
	return "failed to process dapp_file"
}

type WhoIsLookupError struct {
}

func (err *WhoIsLookupError) Error() string {
	return "failed to get domain information"
}

type WhoIsParsingError struct {
}

func (err *WhoIsParsingError) Error() string {
	return "failed to parse domain information"
}

type BlockchainError struct {
}

func (err *BlockchainError) Error() string {
	return "failed to get blockchain information"
}
