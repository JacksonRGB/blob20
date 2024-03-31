package server

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
)

func getMints(c *gin.Context) {
	owner := c.Query("owner")

	if len(owner) != 0 {
		if _, err := hexutil.Decode(owner); err != nil || len(owner) != 42 {
			c.String(400, "owner should be a valid address")
			return
		}
	}

	pageS := c.DefaultQuery("page", "1")
	pageSizeS := c.DefaultQuery("pageSize", "50")

	page, _ := strconv.Atoi(pageS)
	pageSize, _ := strconv.Atoi(pageSizeS)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 200 {
		pageSize = 200
	}

	resp, err := srv.GetMints(owner, page, pageSize)
	if err != nil {
		c.String(500, "internal error")
	}

	c.JSON(200, resp)
}

func getTotalCount(c *gin.Context) {
	ct := srv.GetTotalCount()
	c.JSON(200, gin.H{"count": ct})
}

func getTopAccount(c *gin.Context) {
	rsp := srv.GetTopUsers()
	c.JSON(200, rsp)
}
