package models

import "sync"

type VoteResult struct {
	Candidates  []Candidate `json:"candidates"`
	TotalVoters int         `json:"total_voters"`
	Percentage  float64     `json:"percentage"`
}

type Vote struct {
	VoterID   string `json:"voter_id"`
	Candidate string `json:"candidate"`
	Timestamp uint64 `json:"timestamp"`
}

type Candidate struct {
	CandidateName string  `json:"candidate_name"`
	VoterNum      int     `json:"voter_num"`
	Percentage    float64 `json:"percentage"`
}

type Block struct {
	Index        int    `json:"index"`
	Timestamp    uint64 `json:"timestamp"`
	Votes        []Vote `json:"votes"`
	Hash         string `json:"hash"`
	PreviousHash string `json:"previous_hash"`
	Nonce        int    `json:"nonce"`
}

type Blockchain struct {
	Blocks []*Block
	Mutex  sync.Mutex
}
