package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import "fmt"

const BlockStorageSize = 4096

func NewBlockStorage() *BlockStorage {
	return &BlockStorage{
		Blocks: make([]uint16, BlockStorageSize),
	}
}

type BlockStorage struct {
	Palettes []*BlockState
	Blocks   []uint16
}

func (BlockStorage) At(x, y, z int) int {
	return y<<8 | z<<4 | x
}

// Vaild vailds blockstorage coordinates
func (BlockStorage) Vaild(x, y, z int) error {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return fmt.Errorf("invail coordinate")
	}

	return nil
}

// GetBlock returns the BlockState at blockstorage coordinates
func (storage *BlockStorage) GetBlock(x, y, z int) (*BlockState, error) {
	err := storage.Vaild(x, y, z)
	if err != nil {
		return nil, err
	}

	index := storage.At(x, y, z)

	if index >= len(storage.Blocks) {
		return nil, fmt.Errorf("uninitialized BlockStorage")
	}

	id := storage.Blocks[index]

	if int(id) >= len(storage.Palettes) {
		return nil, fmt.Errorf("couldn't find a palette for the block")
	}

	return storage.Palettes[id], nil
}

// Finalization show the status of a chunk
// It's introduced in mcpe v1.1
type Finalization int

const (
	// Unsupported is unsupported finalization by the chunk format
	Unsupported Finalization = iota

	// NotGenerated is not generated a chunk if it's set
	NotGenerated

	// NotSpawnMobs is not spawned mobs if it's set
	NotSpawnMobs

	// Generated is generated a chunk if it's set
	Generated
)

func GetFinalization(id int) (Finalization, bool) {
	switch id {
	case 0:
		return NotGenerated, true
	case 1:
		return NotSpawnMobs, true
	case 2:
		return Generated, true
	}

	return Unsupported, false
}

func NewSubChunk(y byte) *SubChunk {
	return &SubChunk{
		Y:            y,
		Finalization: NotGenerated,
	}
}

type SubChunk struct {
	Y byte

	Storages []*BlockStorage

	Finalization Finalization
}

// At returns index from subchunk coordinates
// xyz need to be more 0 and less 15
func (SubChunk) At(x, y, z int) int {
	return y<<8 | z<<4 | x
}

// Vaild vailds subchunk coordinates
func (SubChunk) Vaild(x, y, z int) error {
	if x < 0 || x > 15 || y < 0 || y > 15 || z < 0 || z > 15 {
		return fmt.Errorf("invail coordinate")
	}

	return nil
}

// GetBlockStorage returns BlockStorage which subchunk contained with index
func (sub *SubChunk) GetBlockStorage(index int) (*BlockStorage, bool) {
	if len(sub.Storages) >= index || index < 0 {
		return nil, false
	}

	return sub.Storages[index], true
}

// AtBlock returns BlockState at the subchunk coordinates
func (sub *SubChunk) AtBlock(x, y, z, index int) (*BlockState, error) {
	storage, ok := sub.GetBlockStorage(index)
	if !ok {
		return nil, fmt.Errorf("invaild storage index")
	}

	return storage.GetBlock(x, y, z)
}

type SubChunkFormat interface {
	Version() byte
	Read(y byte, b []byte) (*SubChunk, error)
}
