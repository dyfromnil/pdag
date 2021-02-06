package globleconfig

const (
	//PreferredMaxBytes is
	PreferredMaxBytes = 1024 * 1024 * 3
	//MaxMessageCount is
	MaxMessageCount = 100
	//BatchTimeOut is
	BatchTimeOut = 2

	//PostReference for
	PostReference = 3
	//PbftNumOfGoroutine for
	PbftNumOfGoroutine = 1

	//BlockStorageDir is
	BlockStorageDir = "./"
	// ChainsDir is the name of the directory containing the channel ledgers.
	ChainsDir = "store"
	//DefaultMaxBlockfileSize is
	DefaultMaxBlockfileSize = 1024 * 64 // bytes

	//ClientAddr for
	ClientAddr = "client:8688"

	//LeaderListenEnvelopeAddr for
	LeaderListenEnvelopeAddr = "server0:8689"

	//LeaderNodeID for
	LeaderNodeID = "N0"
)

//NodeTable for
var NodeTable = map[string]string{
	"N0": "server0:8600",
	"N1": "server1:8601",
	"N2": "server2:8602",
	"N3": "server3:8603",
}
