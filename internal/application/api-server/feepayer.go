package apiserver

import (
	"fmt"
	"net/http"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"github.com/gin-gonic/gin"
)

func (app *Application) postFeePayer(c *gin.Context) {
	var req FeePayerRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	var requestTxn aptos.RawTransaction
	if err := bcs.Deserialize(&requestTxn, req.Transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid transaction",
			"details": err.Error(),
		})
		return
	}

	var authenticator crypto.AccountAuthenticator
	if err := bcs.Deserialize(&authenticator, req.Signature); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid transaction",
			"details": err.Error(),
		})
		return
	}

	rawTxn, err := app.aptos.BuildTransactionMultiAgent(
		requestTxn.Sender,
		requestTxn.Payload,
		aptos.FeePayer(&app.sponsor.Address),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to build transaction",
			"details": err.Error(),
		})
		return
	}

	sponsorAuth, err := rawTxn.Sign(app.sponsor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to sign transaction",
			"details": err.Error(),
		})
		return
	}

	signedFeePayerTxn, ok := rawTxn.ToFeePayerSignedTransaction(
		&authenticator,
		sponsorAuth,
		[]crypto.AccountAuthenticator{},
	)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to build fee payer",
			"details": "Failed to build fee payer signed transaction",
		})
		return
	}

	submitResult, err := app.aptos.SubmitTransaction(signedFeePayerTxn)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to submit transaction",
			"details": err.Error(),
		})
		return
	}

	userTxn, err := app.aptos.WaitForTransaction(submitResult.Hash)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to wait for transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": userTxn})
}
