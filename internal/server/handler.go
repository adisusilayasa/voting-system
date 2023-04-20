package server

import (
	voteHttp "belajar-blockchain/internal/vote/delivery/http"
	"belajar-blockchain/internal/vote/models"
	voteUsecase "belajar-blockchain/internal/vote/usecase"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func (s *Server) MapHandlers(e *gin.Engine) error {

	e.Use(gin.Recovery())
	e.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "Options"},
		AllowedHeaders:   []string{"*"},
		Debug:            false,
		MaxAge:           300,
	}))

	apiGroup := e.Group("/api")

	bc := &models.Block{
		Index:     0,
		Timestamp: uint64(time.Now().UnixMilli()),
		// Data:          "Genesis Block",
		PreviousHash: "",
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%d%d%s%s%d", bc.Index, bc.Timestamp, bc.Votes, bc.PreviousHash, bc.Nonce)))
	data := hex.EncodeToString(hash[:])

	bc.Hash = data

	// Store the genesis block in the blockchain
	blockchain := &models.Blockchain{
		Blocks: []*models.Block{bc},
	}

	voteUC := voteUsecase.NewVoteUsecase(blockchain)
	voteHttp.MapGetVoteRoutes(apiGroup, voteUC)
	// e.GET("/vote-count")

	return nil
	// Debug:            true,
}
