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
)

//NodeTable for
var NodeTable = map[string]string{
	"n0": "127.0.0.1:8600",
	"n1": "127.0.0.1:8601",
	"n2": "127.0.0.1:8602",
	"n3": "127.0.0.1:8603",
}
