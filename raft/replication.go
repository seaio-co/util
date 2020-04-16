package raft

import (
	"errors"
	"sync"
	"time"
)

const (
	maxFailureScale = 12
	failureWait     = 10 * time.Millisecond
)

var (
	// ErrLogNotFound indicates a given log entry is not available.
	ErrLogNotFound = errors.New("log not found")

	// ErrPipelineReplicationNotSupported can be returned by the transport to
	// signal that pipeline replication is not supported in general, and that
	// no error message should be produced.
	ErrPipelineReplicationNotSupported = errors.New("pipeline replication not supported")
)

// followerReplication is in charge of sending snapshots and log entries from
// this leader during this particular term to a remote follower.
type followerReplication struct {
	// currentTerm and nextIndex must be kept at the top of the struct so
	// they're 64 bit aligned which is a requirement for atomic ops on 32 bit
	// platforms.

	// currentTerm is the term of this leader, to be included in AppendEntries
	// requests.
	currentTerm uint64

	// nextIndex is the index of the next log entry to send to the follower,
	// which may fall past the end of the log.
	nextIndex uint64

	// peer contains the network address and ID of the remote follower.
	peer Server

	// commitment tracks the entries acknowledged by followers so that the
	// leader's commit index can advance. It is updated on successful
	// AppendEntries responses.
	commitment *commitment

	// stopCh is notified/closed when this leader steps down or the follower is
	// removed from the cluster. In the follower removed case, it carries a log
	// index; replication should be attempted with a best effort up through that
	// index, before exiting.
	stopCh chan uint64

	// triggerCh is notified every time new entries are appended to the log.
	triggerCh chan struct{}

	// triggerDeferErrorCh is used to provide a backchannel. By sending a
	// deferErr, the sender can be notifed when the replication is done.
	triggerDeferErrorCh chan *deferError

	// lastContact is updated to the current time whenever any response is
	// received from the follower (successful or not). This is used to check
	// whether the leader should step down (Raft.checkLeaderLease()).
	lastContact time.Time
	// lastContactLock protects 'lastContact'.
	lastContactLock sync.RWMutex

	// failures counts the number of failed RPCs since the last success, which is
	// used to apply backoff.
	failures uint64

	// notifyCh is notified to send out a heartbeat, which is used to check that
	// this server is still leader.
	notifyCh chan struct{}
	// notify is a map of futures to be resolved upon receipt of an
	// acknowledgement, then cleared from this map.
	notify map[*verifyFuture]struct{}
	// notifyLock protects 'notify'.
	notifyLock sync.Mutex

	// stepDown is used to indicate to the leader that we
	// should step down based on information from a follower.
	stepDown chan struct{}

	// allowPipeline is used to determine when to pipeline the AppendEntries RPCs.
	// It is private to this replication goroutine.
	allowPipeline bool
}
