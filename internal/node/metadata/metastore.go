package metadata

import (
	"HorizonFS/internal/node/types"
	log "HorizonFS/pkg/logger"
	"context"
	"fmt"
	json "github.com/bytedance/sonic"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"sync"
)

type Metastore struct {
	metadataStore map[string][]ChunkInfo
	proposeC      chan<- string
	mu            sync.RWMutex
	snapshotter   *snap.Snapshotter
}

func NewMetaSore(snapshotter *snap.Snapshotter, proposeC chan<- string, commitC <-chan *types.Commit, errorC <-chan error) *Metastore {
	s := &Metastore{proposeC: proposeC, metadataStore: make(map[string][]ChunkInfo), snapshotter: snapshotter}
	snapshot, err := s.loadSnapshot()
	if err != nil {
		log.Panic(err)
	}
	if snapshot != nil {
		log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
		if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
			log.Panic(err)
		}
	}
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	return s
}

func (s *Metastore) Lookup(key string) ([]ChunkInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.metadataStore[key]
	return v, ok
}

func (s *Metastore) Propose(ctx context.Context, action string, k string, v string) error {
	kv := types.KV{Action: action, Key: k, Val: v}
	data, err := json.Marshal(kv)
	if err != nil {
		log.Error(err)
		return err
	}
	select {
	case s.proposeC <- string(data):
		return nil // 发送成功
	case <-ctx.Done():
		return fmt.Errorf("Context cancled")
	}
}

func (s *Metastore) GetSnapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.metadataStore)
}

func (s *Metastore) readCommits(commitC <-chan *types.Commit, errorC <-chan error) {
	for commit := range commitC {
		if commit == nil {
			// signaled to load snapshot
			snapshot, err := s.loadSnapshot()
			if err != nil {
				log.Panic(err)
			}
			if snapshot != nil {
				log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
				if err := s.recoverFromSnapshot(snapshot.Data); err != nil {
					log.Panic(err)
				}
			}
			continue
		}

		for _, data := range commit.Data {
			var dataKv types.KV
			if err := json.Unmarshal([]byte(data), &dataKv); err != nil {
				log.Fatalf("metadata: could not decode message (%v)", err)
			}

			log.Println("commit")

			switch dataKv.Action {
			case "SET":
				val, err := deserializeChunkInfo(dataKv.Val)
				if err != nil {
					log.Println(dataKv.Val)
					log.Fatal(err, "Parse error!")
				}
				s.mu.Lock()
				s.metadataStore[dataKv.Key] = val
				s.mu.Unlock()
				log.Printf("v: %+v\n", val)

			case "DEL":
				s.mu.Lock()
				delete(s.metadataStore, dataKv.Key)
				s.mu.Unlock()

			default:
				log.Printf("not valid action %s", dataKv.Action)
			}

		}
		close(commit.ApplyDoneC)
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (s *Metastore) loadSnapshot() (*raftpb.Snapshot, error) {
	snapshot, err := s.snapshotter.Load()
	if err == snap.ErrNoSnapshot {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (s *Metastore) recoverFromSnapshot(snapshot []byte) error {
	var store map[string][]ChunkInfo
	if err := json.Unmarshal(snapshot, &store); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.metadataStore = store
	return nil
}
