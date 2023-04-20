package delivery

import (
	"belajar-blockchain/internal/vote/models"
	"belajar-blockchain/internal/vote/usecase"
	"belajar-blockchain/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	blockchain usecase.BlockchainUsecase
}

func MapGetVoteRoutes(getVoteGroup *gin.RouterGroup, voteUsecase usecase.BlockchainUsecase) {
	h := &VoteHandler{
		blockchain: voteUsecase,
	}

	getVoteGroup.GET("/vote-count", h.GetVoteCountsHandler)
	getVoteGroup.POST("/vote", h.AddVoteHandler)
	getVoteGroup.GET("/blocks", h.GetAllCurrentBlocks)

}

func (h *VoteHandler) GetAllCurrentBlocks(c *gin.Context) {
	c.JSON(http.StatusCreated, utils.SuccessResponse(h.blockchain.GetAllBlocks()))

}
func (h *VoteHandler) AddVoteHandler(c *gin.Context) {
	v := &models.Vote{}
	if err := utils.ReadRequest(c, v); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.blockchain.AddVote(c, v); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(nil))
}

func (h *VoteHandler) GetVoteCountsHandler(c *gin.Context) {
	result, err := h.blockchain.GetVoteCountsHandler()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(result))
}
