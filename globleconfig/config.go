package globleconfig

const (
	//PreferredMaxBytes is
	PreferredMaxBytes = 1024
	//MaxMessageCount is
	MaxMessageCount = 30
	//BatchTimeOut is
	BatchTimeOut = 2

	//BlockStorageDir is
	BlockStorageDir = "./"
	// ChainsDir is the name of the directory containing the channel ledgers.
	ChainsDir = "store"
	//DefaultMaxBlockfileSize is
	DefaultMaxBlockfileSize = 1024 * 64 // bytes

	//ClientAddr for
	ClientAddr = "127.0.0.1:8688"

	//LeaderListenEnvelopeAddr for
	LeaderListenEnvelopeAddr = "127.0.0.1:8689"
)

//NodeTable for
var NodeTable = map[string]string{
	"N0": "127.0.0.1:8600",
	"N1": "127.0.0.1:8601",
	"N2": "127.0.0.1:8602",
	"N3": "127.0.0.1:8603",
}
