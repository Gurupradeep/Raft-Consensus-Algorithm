package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import "labrpc"
import "math"
import "math/rand"
import "time"

// import "bytes"
// import "encoding/gob"



//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

type LogEntry struct {
	Term int
	Command interface{}
}
//
// A Go object implementing a single Raft peer.
//
type AppendEntriesArgs struct {
	Term int
	LeaderId int
	PrevLogIndex int
	PrevLogTerm int
	Entries []*LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term int
	Success bool
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {

	if(rf.isCandidate = true){
		rf.isFollower = true
		rf.isCandidate = false
	}
	if(args.Term<rf.currentTerm){
		reply.Success=false
	}
	else if(rf.*log[args.PrevLogIndex]==nil || rf.*logs[args.PrevLogIndex].Term<args.PrevLogTerm){
		reply.Success=false
	}
	else
	{
		if(len(args.*Entries)==0) {
			reply.Success=true //heartbeat
		}
		else if(rf.*log[args.PrevLogIndex+1]!=nil && rf.*log[args.PrevLogIndex+1].Term<(args.*Entries[0]).Term) {
			for i:=args.PrevLogIndex+1; i <len(rf.*log);i++ {
				rf.*log[i]=nil
			}
			rf.*log[rf.lastLogIndex+1]=args.*Entries[0]
			rf.lastLogIndex=rf.lastLogIndex+1
			rf.lastLogTerm=args.*Entries[0].Term
			reply.Success=true
		}
		else {
			rf.*log[rf.lastLogIndex+1]=args.*Entries[0]
			rf.lastLogIndex=rf.lastLogIndex+1
			rf.lastLogTerm=args.*Entries[0].Term
			reply.Success=true
		}
	}
	if(args.LeaderCommit>rf.commitIndex) {
		rf.commitIndex=Min(args.LeaderCommit,rf.lastLogIndex)
		if(rf.commitIndex>rf.lastApplied){
			rf.lastApplied=rf.lastApplied+1
		}

	}

}
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	currentTerm int
	votedFor	int
	log 		[]*LogEntry
	commitIndex int
	lastApplied int
	nextIndex  []int
	matchIndex []int
	lastLogTerm int
	lastLogIndex int
	electiontimeout time.Duration
	heartbeattimeout time.Duration
	isLeader bool
	isCandidate bool
	isFollower bool
	temp Ticker chan
	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.	
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	term= rf.currentTerm
	if(rf.isLeader == true){
		isleader = true
	}
	else {
		isleader = false
	}
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.xxx)
	// d.Decode(&rf.yyy)
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
}


	

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term int
	CandidateId int
	LastLogIndex int
	LastLogTerm int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term int
	VoteGranted bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	if(args.Term )
	if (args.Term<rf.currentTerm) {
		reply.VoteGranted=false
		reply.Term=rf.currentTerm
	}
	else if((rf.votedFor==nil || rf.votedFor==args.CandidateId)&&(rf.lastLogIndex<args.LastLogIndex || (rf.lastLogTerm==args.LastLogTerm && rf.lastLogIndex<=args.LastLogIndex))){
		rf.votedFor = args.CandidateId
		reply.VoteGranted=true
		rf.currentTerm=Max(rf.currentTerm,args.Term)
		reply.Term=rf.currentTerm
	}
	else{
		reply.VoteGranted=false
		reply.Term=rf.currentTerm
	}

}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}
func (rf *Raft) sendAppendEntries(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).


	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.currentTerm = 0 
	rf.votedFor=nil
	//rf.*log[]=nil
	rf.commitIndex=0
	rf.lastApplied=0
	rf.lastLogTerm=0
	rf.lastLogIndex=0
	rf.electiontimeout=(rand.Intn(250)+500)
	rf.heartbeattimeout=100
	rf.isleader=false
	rf.isFollower = true
	rf.isCandidate = false
	// Your initialization code here (2A, 2B, 2C).

	election_tick := time.NewTicker(time.Millisecond* rf.electiontimeout)
	heartbeat_tick := time.NewTicker(time.Millisecond* rf.heartbeattimeout)
    

    for {
    	select{
    		case <-election_tick.C :{
    			var temp time.Duration
    			temp = (rand.Intn(250)+500)
    			election_tick = time.NewTicker(time.Millisecond*rf.temp)
    			if(rf.isFollower == true || rf.isCandidate == true){
    				rf.isFollower = false
    				rf.isCandidate = true
    				rf.currentTerm = rf.currentTerm + 1
    				var Arguments RequestVoteArgs
    				Arguments.Term = rf.currentTerm
    				Arguments.CandidateId = rf.me
    				Arguments.lastLogIndex = rf.lastLogIndex
    				Arguments.lastLogTerm = rf.lastLogTerm
    				var Reply RequestVoteReply
    				var vote_count = 0
    				rf.votedFor = rf.me
    				vote_count++
    				var no_of_peers = len(rf.peers)
    				for i:=0; i<no_of_peers;i++{
    					if(i!=rf.me){
    						if(sendRequestVote(i,&Arguments,&Reply){
    							if(Reply.Term > rf.currentTerm){
    								rf.currentTerm = Reply.term
    								//check once
    								rf.isFollower=true
    								rf.isCandidate=false
    								rf.isLeader=false
    							}
    							if(Reply.VoteGranted == true){
    								vote_count++;
    							}
    						}
    					}
    				}
    				if(vote_count > no_of_peers/2){
    					rf.isLeader = true
    				}
    				
    			}
    			
    		}   
    		case <-heartbeat_tick.C : {
    			if(rf.isLeader == true){
    				var Arguments AppendEntriesArgs
    				Arguments.Term=rf.currentTerm
    				Arguments.LeaderId=rf.me
    				Arguments.PrevLogIndex=rf.lastLogIndex
    				Arguments.LeaderCommit=rf.commitIndex
    				var Reply AppendEntriesReply
    				for i:=0; i<no_of_peers;i++{
    				if(i!=rf.me){
    				for ok:=false;!ok; {
    				ok=sendAppendEntries(i,&Arguments,&Reply)
    				}
    				if(Reply.Term>rf.currentTerm){
    					rf.isLeader=false
    					rf.isFollower=true
    				}
    			}
    		} 				
    	}
    }
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())


	return rf
}
