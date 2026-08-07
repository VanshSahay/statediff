package main

type Transaction struct {
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Nonce            string `json:"nonce"`
	GasPrice         string `json:"gasPrice"`
	Input            string `json:"input"`
	TransactionIndex string `json:"transactionIndex"`
	Type             string `json:"type"`
}

type Block struct {
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Transactions     []Transaction `json:"transactions"`
	Withdrawals      []Withdrawal  `json:"withdrawals"`
}

type Withdrawal struct {
	Address        string `json:"address"`
	Amount         string `json:"amount"`
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
}

type ProcessedBlock struct {
	Block
	GasDelta     int64
	TimeDelta    int64
	TxAdded      []Transaction
	TxRemoved    []Transaction
	NewContracts []Transaction
}

type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type RPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  Block     `json:"result"`
	Error   *RPCError `json:"error"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BALResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  []AccountAccess `json:"result"`
	Error   *RPCError       `json:"error"`
}

type AccountAccess struct {
	Address        string          `json:"address"`
	BalanceChanges []BalanceChange `json:"balanceChanges"`
	CodeChanges    []CodeChange    `json:"codeChanges"`
	NonceChanges   []NonceChange   `json:"nonceChanges"`
	StorageChanges []SlotChanges   `json:"storageChanges"`
	StorageReads   []string        `json:"storageReads"`
}

type BalanceChange struct {
	Index string `json:"index"`
	Value string `json:"value"`
}

type CodeChange struct {
	Code  string `json:"code"`
	Index string `json:"index"`
}

type NonceChange struct {
	Index string `json:"index"`
	Value string `json:"value"`
}

type SlotChanges struct {
	Key     string          `json:"key"`
	Changes []StorageChange `json:"changes"`
}

type StorageChange struct {
	Index string `json:"index"`
	Value string `json:"value"`
}

type BlockAccessList = []AccountAccess

type BlockEvent struct {
	Number       string          `json:"number"`
	Hash         string          `json:"hash"`
	GasDelta     int64           `json:"gasDelta"`
	TimeDelta    int64           `json:"timeDelta"`
	TxAdded      int             `json:"txAdded"`
	TxRemoved    int             `json:"txRemoved"`
	NewContracts int             `json:"newContracts"`
	Transactions []Transaction   `json:"txs"`
	BAL          BlockAccessList `json:"bal"`
}
