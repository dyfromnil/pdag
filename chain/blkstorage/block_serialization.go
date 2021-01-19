/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blkstorage

import (
	"fmt"
	cb "github.com/dyfromnil/pdag/proto-go/common"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type serializedBlockInfo struct {
	blockHeader *cb.BlockHeader
	txOffsets   []*txindexInfo
}

//The order of the transactions must be maintained for history
type txindexInfo struct {
	txID string
	loc  *locPointer
}

type locPointer struct {
	offset      int
	bytesLength int
}

func (lp *locPointer) String() string {
	return fmt.Sprintf("offset=%d, bytesLength=%d",
		lp.offset, lp.bytesLength)
}

func serializeBlock(block *cb.Block) ([]byte, *serializedBlockInfo, error) {
	buf := proto.NewBuffer(nil)
	var err error
	info := &serializedBlockInfo{}
	info.blockHeader = block.Header
	if err = addHeaderBytes(block.Header, buf); err != nil {
		return nil, nil, err
	}
	if info.txOffsets, err = addDataBytesAndConstructTxIndexInfo(block.Data, buf); err != nil {
		return nil, nil, err
	}

	return buf.Bytes(), info, nil
}

func addHeaderBytes(blockHeader *cb.BlockHeader, buf *proto.Buffer) error {
	if err := buf.EncodeRawBytes(blockHeader.DataHash); err != nil {
		return errors.Wrapf(err, "error encoding the data hash [%v]", blockHeader.DataHash)
	}
	for _, hash := range blockHeader.PreviousHash {
		if err := buf.EncodeRawBytes(hash); err != nil {
			return errors.Wrapf(err, "error encoding the previous hash [%v]", hash)
		}
	}
	return nil
}

func addDataBytesAndConstructTxIndexInfo(blockData *cb.BlockData, buf *proto.Buffer) ([]*txindexInfo, error) {
	var txOffsets []*txindexInfo

	if err := buf.EncodeVarint(uint64(len(blockData.Data))); err != nil {
		return nil, errors.Wrap(err, "error encoding the length of block data")
	}
	for i, txEnvelopeBytes := range blockData.Data {
		offset := len(buf.Bytes())
		txid := fmt.Sprint("txid", i)

		if err := buf.EncodeRawBytes(txEnvelopeBytes); err != nil {
			return nil, errors.Wrap(err, "error encoding the transaction envelope")
		}
		idxInfo := &txindexInfo{txID: txid, loc: &locPointer{offset, len(buf.Bytes()) - offset}}
		txOffsets = append(txOffsets, idxInfo)
	}
	return txOffsets, nil
}
