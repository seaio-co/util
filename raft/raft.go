package raft

import (
	"container/list"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	minCheckInterval = 10 * time.Millisecond
)

var (
	keyCurrentTerm  = []byte("CurrentTerm")
	keyLastVoteTerm = []byte("LastVoteTerm")
	keyLastVoteCand = []byte("LastVoteCand")
)

// getRPCHeader returns an initialized RPCHeader struct for the given
// Raft instance. This structure is sent along with RPC requests and
// responses.
func (r *Raft) getRPCHeader() RPCHeader {
	return RPCHeader{
		ProtocolVersion: r.conf.ProtocolVersion,
	}
}

// checkRPCHeader houses logic about whether this instance of Raft can process
// the given RPC message.
func (r *Raft) checkRPCHeader(rpc RPC) error {
	// Get the header off the RPC message.
	wh, ok := rpc.Command.(WithRPCHeader)
	if !ok {
		return fmt.Errorf("RPC does not have a header")
	}
	header := wh.GetRPCHeader()

	// First check is to just make sure the code can understand the
	// protocol at all.
	if header.ProtocolVersion < ProtocolVersionMin ||
		header.ProtocolVersion > ProtocolVersionMax {
		return ErrUnsupportedProtocol
	}

	// Second check is whether we should support this message, given the
	// current protocol we are configured to run. This will drop support
	// for protocol version 0 starting at protocol version 2, which is
	// currently what we want, and in general support one version back. We
	// may need to revisit this policy depending on how future protocol
	// changes evolve.
	if header.ProtocolVersion < r.conf.ProtocolVersion-1 {
		return ErrUnsupportedProtocol
	}

	return nil
}

// getSnapshotVersion returns the snapshot version that should be used when
// creating snapshots, given the protocol version in use.
func getSnapshotVersion(protocolVersion ProtocolVersion) SnapshotVersion {
	// Right now we only have two versions and they are backwards compatible
	// so we don't need to look at the protocol version.
	return 1
}

// commitTuple is used to send an index that was committed,
// with an optional associated future that should be invoked.
type commitTuple struct {
	log    *Log
	future *logFuture
}

// leaderState is state that is used while we are a leader.
type leaderState struct {
	leadershipTransferInProgress int32 // indicates that a leadership transfer is in progress.
	commitCh                     chan struct{}
	commitment                   *commitment
	inflight                     *list.List // list of logFuture in log index order
	replState                    map[ServerID]*followerReplication
	notify                       map[*verifyFuture]struct{}
	stepDown                     chan struct{}
}

// setLeader is used to modify the current leader of the cluster
func (r *Raft) setLeader(leader ServerAddress) {
	r.leaderLock.Lock()
	oldLeader := r.leader
	r.leader = leader
	r.leaderLock.Unlock()
	if oldLeader != leader {
		r.observe(LeaderObservation{Leader: leader})
	}
}

// requestConfigChange is a helper for the above functions that make
// configuration change requests. 'req' describes the change. For timeout,
// see AddVoter.
func (r *Raft) requestConfigChange(req configurationChangeRequest, timeout time.Duration) IndexFuture {
	var timer <-chan time.Time
	if timeout > 0 {
		timer = time.After(timeout)
	}
	future := &configurationChangeFuture{
		req: req,
	}
	future.init()
	select {
	case <-timer:
		return errorFuture{ErrEnqueueTimeout}
	case r.configurationChangeCh <- future:
		return future
	case <-r.shutdownCh:
		return errorFuture{ErrRaftShutdown}
	}
}

// run is a long running goroutine that runs the Raft FSM.
func (r *Raft) run() {
	for {
		// Check if we are doing a shutdown
		select {
		case <-r.shutdownCh:
			// Clear the leader to prevent forwarding
			r.setLeader("")
			return
		default:
		}

		// Enter into a sub-FSM
		switch r.getState() {
		case Follower:
			r.runFollower()
		case Candidate:
			r.runCandidate()
		case Leader:
			r.runLeader()
		}
	}
}

// runFollower runs the FSM for a follower.
func (r *Raft) runFollower() {
	didWarn := false
	r.logger.Info("entering follower state", "follower", r, "leader", r.Leader())
	metrics.IncrCounter([]string{"raft", "state", "follower"}, 1)
	heartbeatTimer := randomTimeout(r.conf.HeartbeatTimeout)

	for r.getState() == Follower {
		select {
		case rpc := <-r.rpcCh:
			r.processRPC(rpc)

		case c := <-r.configurationChangeCh:
			// Reject any operations since we are not the leader
			c.respond(ErrNotLeader)

		case a := <-r.applyCh:
			// Reject any operations since we are not the leader
			a.respond(ErrNotLeader)

		case v := <-r.verifyCh:
			// Reject any operations since we are not the leader
			v.respond(ErrNotLeader)

		case r := <-r.userRestoreCh:
			// Reject any restores since we are not the leader
			r.respond(ErrNotLeader)

		case r := <-r.leadershipTransferCh:
			// Reject any operations since we are not the leader
			r.respond(ErrNotLeader)

		case c := <-r.configurationsCh:
			c.configurations = r.configurations.Clone()
			c.respond(nil)

		case b := <-r.bootstrapCh:
			b.respond(r.liveBootstrap(b.configuration))

		case <-heartbeatTimer:
			// Restart the heartbeat timer
			heartbeatTimer = randomTimeout(r.conf.HeartbeatTimeout)

			// Check if we have had a successful contact
			lastContact := r.LastContact()
			if time.Now().Sub(lastContact) < r.conf.HeartbeatTimeout {
				continue
			}

			// Heartbeat failed! Transition to the candidate state
			lastLeader := r.Leader()
			r.setLeader("")

			if r.configurations.latestIndex == 0 {
				if !didWarn {
					r.logger.Warn("no known peers, aborting election")
					didWarn = true
				}
			} else if r.configurations.latestIndex == r.configurations.committedIndex &&
				!hasVote(r.configurations.latest, r.localID) {
				if !didWarn {
					r.logger.Warn("not part of stable configuration, aborting election")
					didWarn = true
				}
			} else {
				r.logger.Warn("heartbeat timeout reached, starting election", "last-leader", lastLeader)
				metrics.IncrCounter([]string{"raft", "transition", "heartbeat_timeout"}, 1)
				r.setState(Candidate)
				return
			}

		case <-r.shutdownCh:
			return
		}
	}
}

// liveBootstrap attempts to seed an initial configuration for the cluster. See
// the Raft object's member BootstrapCluster for more details. This must only be
// called on the main thread, and only makes sense in the follower state.
func (r *Raft) liveBootstrap(configuration Configuration) error {
	// Use the pre-init API to make the static updates.
	err := BootstrapCluster(&r.conf, r.logs, r.stable, r.snapshots,
		r.trans, configuration)
	if err != nil {
		return err
	}

	// Make the configuration live.
	var entry Log
	if err := r.logs.GetLog(1, &entry); err != nil {
		panic(err)
	}
	r.setCurrentTerm(1)
	r.setLastLog(entry.Index, entry.Term)
	r.processConfigurationLogEntry(&entry)
	return nil
}

// runCandidate runs the FSM for a candidate.
func (r *Raft) runCandidate() {
	r.logger.Info("entering candidate state", "node", r, "term", r.getCurrentTerm()+1)
	metrics.IncrCounter([]string{"raft", "state", "candidate"}, 1)

	// Start vote for us, and set a timeout
	voteCh := r.electSelf()

	// Make sure the leadership transfer flag is reset after each run. Having this
	// flag will set the field LeadershipTransfer in a RequestVoteRequst to true,
	// which will make other servers vote even though they have a leader already.
	// It is important to reset that flag, because this priviledge could be abused
	// otherwise.
	defer func() { r.candidateFromLeadershipTransfer = false }()

	electionTimer := randomTimeout(r.conf.ElectionTimeout)

	// Tally the votes, need a simple majority
	grantedVotes := 0
	votesNeeded := r.quorumSize()
	r.logger.Debug("votes", "needed", votesNeeded)

	for r.getState() == Candidate {
		select {
		case rpc := <-r.rpcCh:
			r.processRPC(rpc)

		case vote := <-voteCh:
			// Check if the term is greater than ours, bail
			if vote.Term > r.getCurrentTerm() {
				r.logger.Debug("newer term discovered, fallback to follower")
				r.setState(Follower)
				r.setCurrentTerm(vote.Term)
				return
			}

			// Check if the vote is granted
			if vote.Granted {
				grantedVotes++
				r.logger.Debug("vote granted", "from", vote.voterID, "term", vote.Term, "tally", grantedVotes)
			}

			// Check if we've become the leader
			if grantedVotes >= votesNeeded {
				r.logger.Info("election won", "tally", grantedVotes)
				r.setState(Leader)
				r.setLeader(r.localAddr)
				return
			}

		case c := <-r.configurationChangeCh:
			// Reject any operations since we are not the leader
			c.respond(ErrNotLeader)

		case a := <-r.applyCh:
			// Reject any operations since we are not the leader
			a.respond(ErrNotLeader)

		case v := <-r.verifyCh:
			// Reject any operations since we are not the leader
			v.respond(ErrNotLeader)

		case r := <-r.userRestoreCh:
			// Reject any restores since we are not the leader
			r.respond(ErrNotLeader)

		case c := <-r.configurationsCh:
			c.configurations = r.configurations.Clone()
			c.respond(nil)

		case b := <-r.bootstrapCh:
			b.respond(ErrCantBootstrap)

		case <-electionTimer:
			// Election failed! Restart the election. We simply return,
			// which will kick us back into runCandidate
			r.logger.Warn("Election timeout reached, restarting election")
			return

		case <-r.shutdownCh:
			return
		}
	}
}

func (r *Raft) setLeadershipTransferInProgress(v bool) {
	if v {
		atomic.StoreInt32(&r.leaderState.leadershipTransferInProgress, 1)
	} else {
		atomic.StoreInt32(&r.leaderState.leadershipTransferInProgress, 0)
	}
}

func (r *Raft) getLeadershipTransferInProgress() bool {
	v := atomic.LoadInt32(&r.leaderState.leadershipTransferInProgress)
	if v == 1 {
		return true
	}
	return false
}

func (r *Raft) setupLeaderState() {
	r.leaderState.commitCh = make(chan struct{}, 1)
	r.leaderState.commitment = newCommitment(r.leaderState.commitCh,
		r.configurations.latest,
		r.getLastIndex()+1 /* first index that may be committed in this term */)
	r.leaderState.inflight = list.New()
	r.leaderState.replState = make(map[ServerID]*followerReplication)
	r.leaderState.notify = make(map[*verifyFuture]struct{})
	r.leaderState.stepDown = make(chan struct{}, 1)
}
