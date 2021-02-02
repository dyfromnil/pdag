package blkstorage

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	cb "github.com/dyfromnil/pdag/proto-go/common"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	blockfilePrefix = "blockfile_"
)

var (
	blkMgrInfoKey = []byte("blkMgrInfo")
)

type blockfileMgr struct {
	rootDir           string
	conf              *Conf
	blockfilesInfo    *blockfilesInfo
	currentFileWriter *blockfileWriter
	blkfilesInfoCond  *sync.Cond
}

/*
Creates a new manager that will manage the files used for block persistence.
This manager manages the file system FS including
  -- the directory where the files are stored
  -- the individual files where the blocks are stored
  -- the blockfilesInfo which tracks the latest file being persisted to
  -- the index which tracks what block and transaction is in what file
When a new blockfile manager is started (i.e. only on start-up), it checks
if this start-up is the first time the system is coming up or is this a restart
of the system.

The blockfile manager stores blocks of data into a file system.  That file
storage is done by creating sequentially numbered files of a configured size
i.e blockfile_000000, blockfile_000001, etc..

Each transaction in a block is stored with information about the number of
bytes in that transaction
 Adding txLoc [fileSuffixNum=0, offset=3, bytesLength=104] for tx [1:0] to index
 Adding txLoc [fileSuffixNum=0, offset=107, bytesLength=104] for tx [1:1] to index
Each block is stored with the total encoded length of that block as well as the
tx location offsets.

Remember that these steps are only done once at start-up of the system.
At start up a new manager:
  *) Checks if the directory for storing files exists, if not creates the dir
  *) Checks if the key value database exists, if not creates one
       (will create a db dir)
  *) Determines the blockfilesInfo used for storage
		-- Loads from db if exist, if not instantiate a new blockfilesInfo
		-- If blockfilesInfo was loaded from db, compares to FS
		-- If blockfilesInfo and file system are not in sync, syncs blockfilesInfo from FS
  *) Starts a new file writer
		-- truncates file per blockfilesInfo to remove any excess past last block
  *) Determines the index information used to find tx and blocks in
  the file blkstorage
		-- Instantiates a new blockIdxInfo
		-- Loads the index from the db if exists
		-- syncIndex comparing the last block indexed to what is in the FS
		-- If index and file system are not in sync, syncs index from the FS
  *)  Updates blockchain info used by the APIs
*/
func newBlockfileMgr(conf *Conf) (*blockfileMgr, error) {
	fmt.Println("newBlockfileMgr() initializing file-based block storage for ledger")
	rootDir := conf.getLedgerBlockDir()
	_, err := CreateDirIfMissing(rootDir)
	if err != nil {
		panic(fmt.Sprintf("Error creating block storage root dir [%s]: %s", rootDir, err))
	}
	mgr := &blockfileMgr{rootDir: rootDir, conf: conf}

	fmt.Printf(`Getting block information from block storage`)
	var blockfilesInfo *blockfilesInfo
	if blockfilesInfo, err = constructBlockfilesInfo(rootDir); err != nil {
		panic(fmt.Sprintf("Could not build blockfilesInfo info from block files: %s", err))
	}
	fmt.Printf("Info constructed by scanning the blocks dir = %s", spew.Sdump(blockfilesInfo))

	currentFileWriter, err := newBlockfileWriter(deriveBlockfilePath(rootDir, blockfilesInfo.latestFileNumber))
	if err != nil {
		panic(fmt.Sprintf("Could not open writer to current file: %s", err))
	}
	err = currentFileWriter.truncateFile(blockfilesInfo.latestFileSize)
	if err != nil {
		panic(fmt.Sprintf("Could not truncate current file to known size in db: %s", err))
	}

	mgr.blockfilesInfo = blockfilesInfo

	mgr.currentFileWriter = currentFileWriter
	mgr.blkfilesInfoCond = sync.NewCond(&sync.Mutex{})

	return mgr, nil
}

func deriveBlockfilePath(rootDir string, suffixNum int) string {
	return rootDir + "/" + blockfilePrefix + fmt.Sprintf("%06d", suffixNum)
}

func (mgr *blockfileMgr) close() {
	mgr.currentFileWriter.close()
}

func (mgr *blockfileMgr) moveToNextFile() {
	blkfilesInfo := &blockfilesInfo{
		latestFileNumber: mgr.blockfilesInfo.latestFileNumber + 1,
		latestFileSize:   0,
	}

	nextFileWriter, err := newBlockfileWriter(
		deriveBlockfilePath(mgr.rootDir, blkfilesInfo.latestFileNumber))

	if err != nil {
		panic(fmt.Sprintf("Could not open writer to next file: %s", err))
	}
	mgr.currentFileWriter.close()

	mgr.currentFileWriter = nextFileWriter
	mgr.updateBlockfilesInfo(blkfilesInfo)
}

func (mgr *blockfileMgr) addBlock(block *cb.Block) error {
	// Add the previous hash check - Though, not essential but may not be a bad idea to
	// verify the field `block.Header.PreviousHash` present in the block.
	// This check is a simple bytes comparison and hence does not cause any observable performance penalty
	// and may help in detecting a rare scenario if there is any bug in the ordering service.

	blockBytes, info, err := serializeBlock(block)
	if err != nil {
		return errors.WithMessage(err, "error serializing block")
	}
	//Get the location / offset where each transaction starts in the block and where the block ends
	txOffsets := info.txOffsets
	currentOffset := mgr.blockfilesInfo.latestFileSize

	blockBytesLen := len(blockBytes)
	blockBytesEncodedLen := proto.EncodeVarint(uint64(blockBytesLen))
	totalBytesToAppend := blockBytesLen + len(blockBytesEncodedLen)

	//Determine if we need to start a new file since the size of this block
	//exceeds the amount of space left in the current file
	if currentOffset+totalBytesToAppend > mgr.conf.maxBlockfileSize {
		mgr.moveToNextFile()
		currentOffset = 0
	}
	//append blockBytesEncodedLen to the file
	err = mgr.currentFileWriter.append(blockBytesEncodedLen, false)
	if err == nil {
		//append the actual block bytes to the file
		err = mgr.currentFileWriter.append(blockBytes, true)
	}
	if err != nil {
		truncateErr := mgr.currentFileWriter.truncateFile(mgr.blockfilesInfo.latestFileSize)
		if truncateErr != nil {
			panic(fmt.Sprintf("Could not truncate current file to known size after an error during block append: %s", err))
		}
		return errors.WithMessage(err, "error appending block to file")
	}

	//Update the blockfilesInfo with the results of adding the new block
	currentBlkfilesInfo := mgr.blockfilesInfo
	newBlkfilesInfo := &blockfilesInfo{
		latestFileNumber: currentBlkfilesInfo.latestFileNumber,
		latestFileSize:   currentBlkfilesInfo.latestFileSize + totalBytesToAppend,
		noBlockFiles:     false,
	}
	//save the blockfilesInfo in the database

	//Index block file location pointer updated with file suffex and offset for the new block
	blockFLP := &fileLocPointer{fileSuffixNum: newBlkfilesInfo.latestFileNumber}
	blockFLP.offset = currentOffset
	// shift the txoffset because we prepend length of bytes before block bytes
	for _, txOffset := range txOffsets {
		txOffset.loc.offset += len(blockBytesEncodedLen)
	}

	//update the blockfilesInfo (for storage) and the blockchain info (for APIs) in the manager
	mgr.updateBlockfilesInfo(newBlkfilesInfo)
	return nil
}

func (mgr *blockfileMgr) updateBlockfilesInfo(blkfilesInfo *blockfilesInfo) {
	mgr.blkfilesInfoCond.L.Lock()
	defer mgr.blkfilesInfoCond.L.Unlock()
	mgr.blockfilesInfo = blkfilesInfo
	mgr.blkfilesInfoCond.Broadcast()
}

// fileLocPointer
type fileLocPointer struct {
	fileSuffixNum int
	locPointer
}

// blockfilesInfo maintains the summary about the blockfiles
type blockfilesInfo struct {
	latestFileNumber int
	latestFileSize   int
	noBlockFiles     bool
}

// constructBlockfilesInfo scans the last blockfile (if any) and construct the blockfilesInfo
// if the last file contains no block or only a partially written block (potentially because of a crash while writing block to the file),
// this scans the second last file (if any)
func constructBlockfilesInfo(rootDir string) (*blockfilesInfo, error) {
	fmt.Printf("constructing BlockfilesInfo")

	blkfilesInfo := &blockfilesInfo{
		latestFileNumber: 0,
		latestFileSize:   0,
		noBlockFiles:     true,
	}
	return blkfilesInfo, nil
}
