/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blkstorage

import (
	"github.com/dyfromnil/pdag/globleconfig"
	"path/filepath"
)

// Conf encapsulates all the configurations for `BlockStore`
type Conf struct {
	blockStorageDir  string
	maxBlockfileSize int
}

// NewConf constructs new `Conf`.
// blockStorageDir is the top level folder under which `BlockStore` manages its data
func NewConf(blockStorageDir string, maxBlockfileSize int) *Conf {
	if maxBlockfileSize <= 0 {
		maxBlockfileSize = globleconfig.DefaultMaxBlockfileSize
	}
	if blockStorageDir == "" {
		blockStorageDir = globleconfig.BlockStorageDir
	}
	return &Conf{blockStorageDir, maxBlockfileSize}
}

func (conf *Conf) getChainsDir() string {
	return filepath.Join(conf.blockStorageDir, globleconfig.ChainsDir)
}

func (conf *Conf) getLedgerBlockDir(ledgerid string) string {
	return filepath.Join(conf.getChainsDir(), ledgerid)
}
