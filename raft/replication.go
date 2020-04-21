package raft

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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

// notifyAll is used to notify all the waiting verify futures
// if the follower believes we are still the leader.
func (s *followerReplication) notifyAll(leader bool) {
	// Clear the waiting notifies minimizing lock time
	s.notifyLock.Lock()
	n := s.notify
	s.notify = make(map[*verifyFuture]struct{})
	s.notifyLock.Unlock()

	// Submit our votes
	for v := range n {
		v.vote(leader)
	}
}

// cleanNotify is used to delete notify, .
func (s *followerReplication) cleanNotify(v *verifyFuture) {
	s.notifyLock.Lock()
	delete(s.notify, v)
	s.notifyLock.Unlock()
}

// LastContact returns the time of last contact.
func (s *followerReplication) LastContact() time.Time {
	s.lastContactLock.RLock()
	last := s.lastContact
	s.lastContactLock.RUnlock()
	return last
}

// setLastContact sets the last contact to the current time.
func (s *followerReplication) setLastContact() {
	s.lastContactLock.Lock()
	s.lastContact = time.Now()
	s.lastContactLock.Unlock()
}

// replicate is a long running routine that replicates log entries to a single
// follower.
func (r *Raft) replicate(s *followerReplication) {
	// Start an async heartbeating routing
	stopHeartbeat := make(chan struct{})
	defer close(stopHeartbeat)
	r.goFunc(func() { r.heartbeat(s, stopHeartbeat) })

RPC:
	shouldStop := false
	for !shouldStop {
		select {
		case maxIndex := <-s.stopCh:
			// Make a best effort to replicate up to this index
			if maxIndex > 0 {
				r.replicateTo(s, maxIndex)
			}
			return
		case deferErr := <-s.triggerDeferErrorCh:
			lastLogIdx, _ := r.getLastLog()
			shouldStop = r.replicateTo(s, lastLogIdx)
			if !shouldStop {
				deferErr.respond(nil)
			} else {
				deferErr.respond(fmt.Errorf("replication failed"))
			}
		case <-s.triggerCh:
			lastLogIdx, _ := r.getLastLog()
			shouldStop = r.replicateTo(s, lastLogIdx)
		// This is _not_ our heartbeat mechanism but is to ensure
		// followers quickly learn the leader's commit index when
		// raft commits stop flowing naturally. The actual heartbeats
		// can't do this to keep them unblocked by disk IO on the
		// follower. See https://github.com/hashicorp/raft/issues/282.
		case <-randomTimeout(r.conf.CommitTimeout):
			lastLogIdx, _ := r.getLastLog()
			shouldStop = r.replicateTo(s, lastLogIdx)
		}

		// If things looks healthy, switch to pipeline mode
		if !shouldStop && s.allowPipeline {
			goto PIPELINE
		}
	}
	return

PIPELINE:
	// Disable until re-enabled
	s.allowPipeline = false

	// Replicates using a pipeline for high performance. This method
	// is not able to gracefully recover from errors, and so we fall back
	// to standard mode on failure.
	if err := r.pipelineReplicate(s); err != nil {
		if err != ErrPipelineReplicationNotSupported {
			r.logger.Error("failed to start pipeline replication to", "peer", s.peer, "error", err)
		}
	}
	goto RPC
}

// replicateTo is a helper to replicate(), used to replicate the logs up to a
// given last index.
// If the follower log is behind, we take care to bring them up to date.
func (r *Raft) replicateTo(s *followerReplication, lastIndex uint64) (shouldStop bool) {
	// Create the base request
	var req AppendEntriesRequest
	var resp AppendEntriesResponse
	var start time.Time
START:
	// Prevent an excessive retry rate on errors
	if s.failures > 0 {
		select {
		case <-time.After(backoff(failureWait, s.failures, maxFailureScale)):
		case <-r.shutdownCh:
		}
	}

	// Setup the request
	if err := r.setupAppendEntries(s, &req, atomic.LoadUint64(&s.nextIndex), lastIndex); err == ErrLogNotFound {
		goto SEND_SNAP
	} else if err != nil {
		return
	}

	// Make the RPC call
	start = time.Now()
	if err := r.trans.AppendEntries(s.peer.ID, s.peer.Address, &req, &resp); err != nil {
		r.logger.Error("failed to appendEntries to", "peer", s.peer, "error", err)
		s.failures++
		return
	}
	appendStats(string(s.peer.ID), start, float32(len(req.Entries)))

	// Check for a newer term, stop running
	if resp.Term > req.Term {
		r.handleStaleTerm(s)
		return true
	}

	// Update the last contact
	s.setLastContact()

	// Update s based on success
	if resp.Success {
		// Update our replication state
		updateLastAppended(s, &req)

		// Clear any failures, allow pipelining
		s.failures = 0
		s.allowPipeline = true
	} else {
		atomic.StoreUint64(&s.nextIndex, max(min(s.nextIndex-1, resp.LastLog+1), 1))
		if resp.NoRetryBackoff {
			s.failures = 0
		} else {
			s.failures++
		}
		r.logger.Warn("appendEntries rejected, sending older logs", "peer", s.peer, "next", atomic.LoadUint64(&s.nextIndex))
	}

CHECK_MORE:
	// Poll the stop channel here in case we are looping and have been asked
	// to stop, or have stepped down as leader. Even for the best effort case
	// where we are asked to replicate to a given index and then shutdown,
	// it's better to not loop in here to send lots of entries to a straggler
	// that's leaving the cluster anyways.
	select {
	case <-s.stopCh:
		return true
	default:
	}

	// Check if there are more logs to replicate
	if atomic.LoadUint64(&s.nextIndex) <= lastIndex {
		goto START
	}
	return

	// SEND_SNAP is used when we fail to get a log, usually because the follower
	// is too far behind, and we must ship a snapshot down instead
SEND_SNAP:
	if stop, err := r.sendLatestSnapshot(s); stop {
		return true
	} else if err != nil {
		r.logger.Error("failed to send snapshot to", "peer", s.peer, "error", err)
		return
	}

	// Check if there is more to replicate
	goto CHECK_MORE
}
