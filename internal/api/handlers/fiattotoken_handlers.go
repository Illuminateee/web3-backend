// fiat to token handlers
package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// FiatToTokenHandler calculates how many CIFO tokens can be purchased with a given fiat amount
func (h *Handler) FiatToTokenHandler(c *gin.Context) {
    // Get fiat amount and currency from query
    fiatAmountStr := c.Query("amount")
    currency := strings.ToLower(c.DefaultQuery("currency", "idr")) // Default to IDR

    if fiatAmountStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Fiat amount is required"})
        return
    }

    fiatAmount, err := strconv.ParseFloat(fiatAmountStr, 64)
    if err != nil || fiatAmount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fiat amount"})
        return
    }

    // Validate currency
    if currency != "idr" && currency != "usd" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Currency must be either 'idr' or 'usd'"})
        return
    }

    // Get ETH prices
    ethPriceUSD, ethPriceIDR, err := h.GetEthPrices(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ETH prices"})
        return
    }

    // Convert fiat amount to ETH
    var ethAmount float64
    if currency == "idr" {
        ethAmount = fiatAmount / ethPriceIDR
    } else {
        ethAmount = fiatAmount / ethPriceUSD
    }

    // Get CIFO amount for the calculated ETH
    cifoAmount, err := h.getCifoAmount(ethAmount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch CIFO conversion rate"})
        return
    }

    // Parse cifoAmount as float for calculations
    cifoAmountFloat, err := strconv.ParseFloat(cifoAmount, 64)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CIFO amount"})
        return
    }

    // Format the response
    response := gin.H{
        "fiat_currency":  strings.ToUpper(currency),
        "fiat_amount":    fiatAmount,
        "eth_amount":     ethAmount,
        "cifo_amount":    cifoAmount,
        "eth_price_usd":  ethPriceUSD,
        "eth_price_idr":  ethPriceIDR,
        "exchange_rate":  map[string]interface{}{
            "cifo_per_eth":      cifoAmountFloat / ethAmount,
            "eth_per_cifo":      ethAmount / cifoAmountFloat,
            "cifo_per_usd":      cifoAmountFloat / (ethAmount * ethPriceUSD),
            "cifo_per_idr":      cifoAmountFloat / (ethAmount * ethPriceIDR),
            "usd_per_cifo":      (ethAmount * ethPriceUSD) / cifoAmountFloat,
            "idr_per_cifo":      (ethAmount * ethPriceIDR) / cifoAmountFloat,
        },
    }

    c.JSON(http.StatusOK, response)
}