package types

// Event types
const (
	EventTypeWalletCreated         = "wallet_created"
	EventTypeWalletRecovered       = "wallet_recovered"
	EventTypeKeyRotationInitiated  = "key_rotation_initiated"
	EventTypeKeyRotationCompleted  = "key_rotation_completed"
	EventTypeTransactionSigned     = "transaction_signed"
	EventTypeTransactionBatchSigned = "transaction_batch_signed"
	EventTypeBatchSignRequested    = "batch_sign_requested"
)

// Attribute keys
const (
	AttributeKeyWalletAddress = "wallet_address"
	AttributeKeyNewPubKey     = "new_pub_key"
	AttributeKeyVersion       = "version"
	AttributeKeyCreator       = "creator"
	AttributeKeyChainId       = "chain_id"
	AttributeKeyTxHash       = "tx_hash"
	AttributeKeyBatchSize    = "batch_size"
	AttributeKeyBatchStatus  = "batch_status"
	AttributeKeyStatus       = "status"
)
