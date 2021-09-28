package globleconfig

const (
	//NumOfClient for
	NumOfClient = 13

	//PreferredMaxBytes is
	PreferredMaxBytes = 1024 * 1024 * 3
	//MaxMessageCount is
	MaxMessageCount = 100
	//BatchTimeOut is
	BatchTimeOut = 2

	//PostReference for
	PostReference = 3
	//PreReference for
	PreReference = 6
	//NumOfConsensusGoroutine for
	NumOfConsensusGoroutine = 8

	//Rate for
	Rate = 0.85

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
	"N4": "server4:8604",
	"N5": "server5:8605",
	"N6": "server6:8606",
	"N7": "server7:8607",

	"N8":  "server8:8608",
	"N9":  "server9:8609",
	"N10": "server10:8610",
	"N11": "server11:8611",
	"N12": "server12:8612",
	"N13": "server13:8613",
	"N14": "server14:8614",
	"N15": "server15:8615",

	"N16": "server16:8616",
	"N17": "server17:8617",
	"N18": "server18:8618",
	"N19": "server19:8619",
	"N20": "server20:8620",
	"N21": "server21:8621",
	"N22": "server22:8622",
	"N23": "server23:8623",

	"N24": "server24:8624",
	"N25": "server25:8625",
	"N26": "server26:8626",
	"N27": "server27:8627",
	"N28": "server28:8628",
	"N29": "server29:8629",
	"N30": "server30:8630",
	"N31": "server31:8631",
}
