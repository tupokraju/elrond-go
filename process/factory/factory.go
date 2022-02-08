package factory

const (
	// TransactionTopic is the topic used for sharing transactions
	TransactionTopic = "transactions"
	// UnsignedTransactionTopic is the topic used for sharing unsigned transactions
	UnsignedTransactionTopic = "unsignedTransactions"
	// RewardsTransactionTopic is the topic used for sharing fee transactions
	RewardsTransactionTopic = "rewardsTransactions"
	// ShardBlocksTopic is the topic used for sharing block headers
	ShardBlocksTopic = "shardBlocks"
	// MiniBlocksTopic is the topic used for sharing mini blocks
	MiniBlocksTopic = "txBlockBodies"
	// PeerChBodyTopic is the topic used for sharing peer change block bodies
	PeerChBodyTopic = "peerChangeBlockBodies"
	// MetachainBlocksTopic is the topic used for sharing metachain block headers
	MetachainBlocksTopic = "metachainBlocks"
	// AccountTrieNodesTopic is used for sharing state trie nodes
	AccountTrieNodesTopic = "accountTrieNodes"
	// ValidatorTrieNodesTopic is used for sharding validator state trie nodes
	ValidatorTrieNodesTopic = "validatorTrieNodes"
	// PeerAuthenticationTopic is used for sharing peer authentication messages
	PeerAuthenticationTopic = "peerAuthentication"
	// HeartbeatTopic is used for sharing heartbeat messages
	HeartbeatTopic = "heartbeat"
)

// SystemVirtualMachine is a byte array identifier for the smart contract address created for system VM
var SystemVirtualMachine = []byte{0, 1}

// IELEVirtualMachine is a byte array identifier for the smart contract address created for IELE VM
var IELEVirtualMachine = []byte{1, 0}

// ArwenVirtualMachine is a byte array identifier for the smart contract address created for Arwen VM
var ArwenVirtualMachine = []byte{5, 0}

// InternalTestingVM is a byte array identified for the smart contract address created for the testing VM
var InternalTestingVM = []byte{255, 255}
