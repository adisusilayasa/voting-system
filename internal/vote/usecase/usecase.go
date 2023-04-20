package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"belajar-blockchain/internal/vote/models"
	"belajar-blockchain/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BlockchainUsecase interface {
	LatestBlock() *models.Block
	AddBlock(b *models.Block) error
	AddVote(c *gin.Context, v *models.Vote) error
	GetVoteCountsHandler() (*models.VoteResult, error)
	GetAllBlocks() *models.Blockchain
}

func NewVoteUsecase(blockChain *models.Blockchain) BlockchainUsecase {
	return &Blockchain{
		Blockchain: blockChain,
	}
}

type Blockchain struct {
	Blockchain *models.Blockchain
	// Block      *models.Block
}

func (bc *Blockchain) LatestBlock() *models.Block {
	return bc.Blockchain.Blocks[len(bc.Blockchain.Blocks)-1]
}

func (bc *Blockchain) GetAllBlocks() *models.Blockchain {
	return bc.Blockchain
}

func (bc *Blockchain) AddBlock(b *models.Block) error {
	bc.Blockchain.Mutex.Lock()
	defer bc.Blockchain.Mutex.Unlock()

	if len(bc.Blockchain.Blocks) == 0 {
		// Set previous hash to an empty string for the first block
		b.PreviousHash = ""
	} else {
		// Set previous hash to the hash of the latest block
		b.PreviousHash = bc.LatestBlock().Hash
	}

	b.Index = len(bc.Blockchain.Blocks)
	b.Timestamp = uint64(time.Now().UnixMilli())

	b.Hash = bc.CalculateHash(b)

	// Validate the hash
	if b.Hash != bc.CalculateHash(b) {
		return errors.New("Invalid block hash")
	}

	bc.Blockchain.Blocks = append(bc.Blockchain.Blocks, b)
	return nil
}

func (bc *Blockchain) CalculateHash(b *models.Block) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d%d%s%s%d", b.Index, b.Timestamp, b.Votes, b.PreviousHash, b.Nonce)))
	return hex.EncodeToString(hash[:])
}

func (bc *Blockchain) AddVote(c *gin.Context, v *models.Vote) error {

	// Generate a random ID for the voter
	// id := utils.GenerateRandomID()

	// Set the generated ID for the vote
	// hash := sha256.Sum256([]byte(id))
	v.VoterID = utils.EncryptID(v.VoterID)
	// Validate the vote
	err := bc.isExistingVoter(v.VoterID)
	if err != nil {
		return err
	}

	// Create a new block containing the vote
	b := &models.Block{Votes: []models.Vote{*v}}

	// Add the block to the blockchain
	bc.AddBlock(b)

	return nil
}

func (bc *Blockchain) isExistingVoter(voterID string) error {
	for _, block := range bc.Blockchain.Blocks {
		for _, vote := range block.Votes {
			if vote.VoterID == voterID {
				return errors.New("User already voted")
			}
		}
	}
	return nil
}

func (bc *Blockchain) validateVoter(c *gin.Context, voterID string) error {
	// Check if the voter ID is valid (e.g. length, format, etc.)
	// ...
	log.Println("Jumlah voter", len(bc.Blockchain.Blocks))
	if len(bc.Blockchain.Blocks) == 0 {
		return nil
	}

	// Check if the voter ID already exists in the blockchain
	for _, block := range bc.Blockchain.Blocks {
		for _, vote := range block.Votes {
			log.Println(vote.VoterID)
			if vote.VoterID == voterID {
				return errors.New("Voter already exists")
			}
		}
	}

	// Voter ID is valid and does not exist in the blockchain
	return nil
}

func (bc *Blockchain) GetVoteCountsHandler() (*models.VoteResult, error) {
	// Check if the blockchain has any blocks
	if len(bc.Blockchain.Blocks) == 1 {
		err := errors.New("Blockchain has no blocks")
		return nil, err
	}

	// Create a map to store the vote counts for each candidate
	voteCounts := make(map[string]int)

	// Iterate over the blockchain and count the votes for each candidate
	for _, block := range bc.Blockchain.Blocks {
		if block.Votes != nil {
			for _, vote := range block.Votes {
				candidate := vote.Candidate
				voteCounts[candidate]++
			}
		}
	}

	// Calculate the total number of voters
	totalVoters := len(bc.Blockchain.Blocks) - 1

	// Create a nested JSON object containing the vote count and percentage for each candidate
	resultCandidate := []models.Candidate{}
	for candidate, count := range voteCounts {
		percentage := float64(count) / float64(totalVoters) * 100

		newData, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", percentage), 64)

		result := models.Candidate{
			CandidateName: candidate,
			VoterNum:      count,
			Percentage:    newData,
		}

		resultCandidate = append(resultCandidate, result)
	}

	results := &models.VoteResult{
		TotalVoters: totalVoters,
		Candidates:  resultCandidate,
	}

	return results, nil
}
