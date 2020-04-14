package raft

import "io"

// SnapshotMeta is for metadata of a snapshot.
type SnapshotMeta struct {
	// Version is the version number of the snapshot metadata. This does not cover
	// the application's data in the snapshot, that should be versioned
	// separately.
	Version SnapshotVersion

	// ID is opaque to the store, and is used for opening.
	ID string

	// Index and Term store when the snapshot was taken.
	Index uint64
	Term  uint64

	// Peers is deprecated and used to support version 0 snapshots, but will
	// be populated in version 1 snapshots as well to help with upgrades.
	Peers []byte

	// Configuration and ConfigurationIndex are present in version 1
	// snapshots and later.
	Configuration      Configuration
	ConfigurationIndex uint64

	// Size is the size of the snapshot in bytes.
	Size int64
}

// SnapshotStore interface is used to allow for flexible implementations
// of snapshot storage and retrieval. For example, a client could implement
// a shared state store such as S3, allowing new nodes to restore snapshots
// without streaming from the leader.
type SnapshotStore interface {
	// Create is used to begin a snapshot at a given index and term, and with
	// the given committed configuration. The version parameter controls
	// which snapshot version to create.
	Create(version SnapshotVersion, index, term uint64, configuration Configuration,
		configurationIndex uint64, trans Transport) (SnapshotSink, error)

	// List is used to list the available snapshots in the store.
	// It should return then in descending order, with the highest index first.
	List() ([]*SnapshotMeta, error)

	// Open takes a snapshot ID and provides a ReadCloser. Once close is
	// called it is assumed the snapshot is no longer needed.
	Open(id string) (*SnapshotMeta, io.ReadCloser, error)
}

type SnapshotSink interface {
	io.WriteCloser
	ID() string
	Cancel() error
}

// runSnapshots is a long running goroutine used to manage taking
// new snapshots of the FSM. It runs in parallel to the FSM and
// main goroutines, so that snapshots do not block normal operation.
func (r *Raft) runSnapshots() {
	for {
		select {
		case <-randomTimeout(r.conf.SnapshotInterval):
			// Check if we should snapshot
			if !r.shouldSnapshot() {
				continue
			}

			// Trigger a snapshot
			if _, err := r.takeSnapshot(); err != nil {
				r.logger.Error("failed to take snapshot", "error", err)
			}

		case future := <-r.userSnapshotCh:
			// User-triggered, run immediately
			id, err := r.takeSnapshot()
			if err != nil {
				r.logger.Error("failed to take snapshot", "error", err)
			} else {
				future.opener = func() (*SnapshotMeta, io.ReadCloser, error) {
					return r.snapshots.Open(id)
				}
			}
			future.respond(err)

		case <-r.shutdownCh:
			return
		}
	}
}

// shouldSnapshot checks if we meet the conditions to take
// a new snapshot.
func (r *Raft) shouldSnapshot() bool {
	// Check the last snapshot index
	lastSnap, _ := r.getLastSnapshot()

	// Check the last log index
	lastIdx, err := r.logs.LastIndex()
	if err != nil {
		r.logger.Error("failed to get last log index", "error", err)
		return false
	}

	// Compare the delta to the threshold
	delta := lastIdx - lastSnap
	return delta >= r.conf.SnapshotThreshold
}
