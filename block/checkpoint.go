package block

import (
	"sync"
)

const MINIBLOCK_LENGTH = 9
const CP_MAX_DIFF = 3

type Checkpoint struct {
	Checkpoints map[MiniBlockKey][]MiniBlock
	sync.RWMutex
}

// create a checkpoint list
func CreateCheckpoint() *Checkpoint {
	return &Checkpoint{Checkpoints: map[MiniBlockKey][]MiniBlock{}}
}

// get all miniblocks
func (d *Checkpoint) GetAllMiniBlocksFromCheckpoint(key MiniBlockKey) (mbls []MiniBlock) {
	d.RLock()
	defer d.RUnlock()

	for _, mbl := range d.Checkpoints[key] {
		mbls = append(mbls, mbl)
	}
	return
}

// purge all heights less than this height
func (d *Checkpoint) PurgeHeight(height int64) (purge_count int) {
	if height < 0 {
		return
	}
	d.Lock()
	defer d.Unlock()

	for k, _ := range d.Checkpoints {
		if k.Height <= uint64(height) {
			purge_count++
			delete(d.Checkpoints, k)
		}
	}
	return purge_count
}

func (d *Checkpoint) Exists(key MiniBlockKey) bool {
	d.RLock()
	defer d.RUnlock()

	if mbls := d.GetAllMiniBlocksFromCheckpoint(key); mbls != nil {
		return true
	}

	return false
}

func (d *Checkpoint) VerifyCheckpoint(key MiniBlockKey, mbl MiniBlock) bool {
	d.RLock()
	defer d.RUnlock()

	mbls_checkpoint := d.GetAllMiniBlocksFromCheckpoint(key)

	for i := range mbls_checkpoint {
		if mbls_checkpoint[i].Height != mbl.Height {
			return false
		}
		if mbls_checkpoint[i].Past != mbl.Past {
			return false
		}
	}

	return true
}
