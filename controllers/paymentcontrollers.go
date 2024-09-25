package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harmannkibue/golang-mpesa-sdk/pkg/daraja"
)

const (
	mpesaApiKey         = "B3x15APtzLXSAAmjIMC93Dj0JKb4AhOpGAHpKE8imUBH6OcP"
	mpesaConsumerSecret = "XfvpXnoXtdbLNjkkbbAAbUhGYt1EtodNbNWUb1AOj9X6e9A1uzC1a3eiamlGN1im"
	mpesaPassKey        = "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
)

// STKPushRequest is the struct to handle incoming request body for STK push
type STKPushRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

// AccountBalanceRequest is the struct to handle incoming request body for balance inquiry
type AccountBalanceRequest struct {
	IdentifierType int    `json:"identifier_type" binding:"required"` // IdentifierType: 1 for MSISDN, 2 for Till Number, 4 for ShortCode
	PartyA         string `json:"party_a" binding:"required"`
}

// ReversalRequest is the struct to handle incoming request body for reversal
type ReversalRequest struct {
	TransactionID  string `json:"transaction_id" binding:"required"`
	Amount         int    `json:"amount" binding:"required"`
	ReceiverParty  int    `json:"receiver_party" binding:"required"`
	Remarks        string `json:"remarks"`
	Occassion      string `json:"occassion"`
}

// B2CRequest is the struct to handle incoming request body for B2C payment
type B2CRequest struct {
	Amount      int    `json:"amount" binding:"required"`
	PartyA      int    `json:"party_a" binding:"required"`
	PartyB      int64    `json:"party_b" binding:"required"`
	Remarks     string `json:"remarks"`
	Occassion   string `json:"occasion"`
}

// Callback Structs
type STKCallbackBody struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

type B2CCallbackBody struct {
	Result struct {
		ResultType   int    `json:"ResultType"`
		ResultCode   int    `json:"ResultCode"`
		ResultDesc   string `json:"ResultDesc"`
		Originator   string `json:"OriginatorConversationID"`
		Conversation string `json:"ConversationID"`
		Transaction  string `json:"TransactionID"`
	} `json:"Result"`
}

func GenerateSTKPush(c *gin.Context) {
	var reqBody STKPushRequest

	// Bind incoming JSON to reqBody
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	darajaService, err := daraja.New(mpesaApiKey, mpesaConsumerSecret, mpesaPassKey, daraja.SANDBOX)
	if err != nil {
		log.Println("Failed to initialize Safaricom Daraja client: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Daraja client"})
		return
	}

	// Initiate STK push with the phone number from the request body
	stkRes, err := darajaService.InitiateStkPush(daraja.STKPushBody{
		BusinessShortCode: "174379",
		TransactionType:   "CustomerBuyGoodsOnline",
		Amount:            reqBody.Amount,
		PartyA:            reqBody.PhoneNumber,
		PartyB:            "174379",
		PhoneNumber:       reqBody.PhoneNumber,
		CallBackURL:       "https://7936-197-248-208-191.ngrok-free.app/your-callback-endpoint", // Updated
		AccountReference:  "999200200",
		TransactionDesc:   "Daraja SDK testing STK push",
	})

	if err != nil {
		log.Println("STK Push Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate STK push"})
		return
	}

	log.Printf("STK push response: %+v \n", stkRes)
	c.JSON(http.StatusOK, gin.H{"message": "STK push initiated successfully", "response": stkRes})
}

func GetAccountBalance(c *gin.Context) {
	var reqBody AccountBalanceRequest

	// Bind incoming JSON to reqBody
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Convert PartyA to integer
	partyAInt, err := strconv.Atoi(reqBody.PartyA)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid party_a value"})
		return
	}

	darajaService, err := daraja.New(mpesaApiKey, mpesaConsumerSecret, mpesaPassKey, daraja.SANDBOX)
	if err != nil {
		log.Println("Failed to initialize Safaricom Daraja client: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Daraja client"})
		return
	}

	// Make account balance inquiry
	balanceResponse, err := darajaService.QueryAccountBalance(daraja.AccountBalanceRequestBody{
		Initiator:          "testapi",
		SecurityCredential: "UKCrm4IVKWEoW640M3pUHS4hZ2ynDpz+LT6c+acBK28TOMULxVhMP0YM2FNCh2QXx+m6HR8iLNsR0bfbIB1kpvNhciKUrn7Glp4f7UNPF8mHXgNsa/09+i7X8+JUy7tQLEOoPE/xCWBOh2ofBq8N+lX77RUAxDp9HC8Nj6nN6kH07Ygmz7NnRd/dlayqcFKV4UNP/nQAV8lum2HSh9xRBnlexcziYipt/d293qrSSvXtAfz+lmgzzbzwML02zlCQxXS2YQjTluQWzRgxkl+9aCCs51a5BWppTE6iYd8qcMlX/+hMZvl2D9LjQKwisSKJsWP2MtxFxG86DRpwI41I4A==",
		CommandID:          "AccountBalance",
		PartyA:             partyAInt,
		IdentifierType:     reqBody.IdentifierType,
		Remarks:            "Balance",
		QueueTimeOutURL:    "https://webhook.site/bbca16b1-fc3b-4a9f-9a91-14c08972657e",
		ResultURL:          "https://webhook.site/bbca16b1-fc3b-4a9f-9a91-14c08972657e",
	})

	if err != nil {
		log.Println("Balance Inquiry Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inquire balance"})
		return
	}

	log.Printf("Balance inquiry response: %+v \n", balanceResponse)
	c.JSON(http.StatusOK, gin.H{"message": "Account balance inquiry successful", "response": balanceResponse})
}



func ReverseTransaction(c *gin.Context) {
	var reqBody ReversalRequest

	// Bind incoming JSON to reqBody
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	darajaService, err := daraja.New(mpesaApiKey, mpesaConsumerSecret, mpesaPassKey, daraja.SANDBOX)
	if err != nil {
		log.Println("Failed to initialize Safaricom Daraja client: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Daraja client"})
		return
	}

	// Perform transaction reversal
	reversalResponse, err := darajaService.C2BTransactionReversal(daraja.TransactionReversalRequestBody{
		Initiator:              "testapi",
		SecurityCredential:     "oDx3GjKUpc3LJyPMdjiy2Qy64b+Smfyc8xyPTjYQfpGhVngg8OATaXYla0YazHGtM8rqqlRwGiW30NDTezm81YBpEwCvIWTaR1YN3RmiPPvN+kF03BgX8eCJXVzV/99758nSsEKmleudOMmkegHaTrMOlfjQlcVSiS94u2ZvJejS0X5xpp2dPkplITmpLBh/EpMsB0fJLh7fcrtc8v0V/NJG6Zd6W3d2uB3S6zfJPbc4Iby52iYhAWwFOAbJhrTMVDHKLLCzFXZUZufPpntWcElNAtgEb7AA1Os2FbNyJcrCwT22ATQaU/VMJTjMgMB3Cgdw7Xyw+gMilJ+er/kJzA==", // This should be generated dynamically
		CommandID:              "TransactionReversal",
		TransactionID:          reqBody.TransactionID,
		Amount:                 reqBody.Amount,
		ReceiverParty:          reqBody.ReceiverParty,
		ReceiverIdentifierType: 11, // Identifier type, e.g., MSISDN
		ResultURL:              "https://webhook.site/7da5ccfd-3a90-4038-b822-273887b3de7f",
		QueueTimeOutURL:        "https://webhook.site/7da5ccfd-3a90-4038-b822-273887b3de7f",
		Remarks:                reqBody.Remarks,
		Occassion:              reqBody.Occassion,
	})

	if err != nil {
		log.Println("Reversal Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reverse transaction"})
		return
	}

	log.Printf("Reversal response: %+v \n", reversalResponse)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction reversed successfully", "response": reversalResponse})
}


func TransferB2C(c *gin.Context) {
	var reqBody B2CRequest

	// Bind incoming JSON to reqBody
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	darajaService, err := daraja.New(mpesaApiKey, mpesaConsumerSecret, mpesaPassKey, daraja.SANDBOX)
	if err != nil {
		log.Println("Failed to initialize Safaricom Daraja client: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Daraja client"})
		return
	}

	// Perform B2C Payment
	b2cResponse, err := darajaService.B2CPayment(daraja.B2CRequestBody{
		InitiatorName:      "testapi",
		SecurityCredential: "UKCrm4IVKWEoW640M3pUHS4hZ2ynDpz+LT6c+acBK28TOMULxVhMP0YM2FNCh2QXx+m6HR8iLNsR0bfbIB1kpvNhciKUrn7Glp4f7UNPF8mHXgNsa/09+i7X8+JUy7tQLEOoPE/xCWBOh2ofBq8N+lX77RUAxDp9HC8Nj6nN6kH07Ygmz7NnRd/dlayqcFKV4UNP/nQAV8lum2HSh9xRBnlexcziYipt/d293qrSSvXtAfz+lmgzzbzwML02zlCQxXS2YQjTluQWzRgxkl+9aCCs51a5BWppTE6iYd8qcMlX/+hMZvl2D9LjQKwisSKJsWP2MtxFxG86DRpwI41I4A==", // Generate this securely
		CommandID:          "SalaryPayment",
		Amount:             reqBody.Amount,
		PartyA:             reqBody.PartyA,
		PartyB:             int64(reqBody.PartyB),
		Remarks:            reqBody.Remarks,
		QueueTimeOutURL:    "https://webhook.site/7da5ccfd-3a90-4038-b822-273887b3de7f", // Replace with actual URL
		ResultURL:          "https://webhook.site/7da5ccfd-3a90-4038-b822-273887b3de7f", // Replace with actual URL
		Occassion:          reqBody.Occassion,
	})

	if err != nil {
		log.Println("B2C Payment Error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process B2C payment"})
		return
	}

	log.Printf("B2C response: %+v \n", b2cResponse)
	c.JSON(http.StatusOK, gin.H{"message": "B2C payment successful", "response": b2cResponse})
}

// Handle STK Push Callback
func HandleSTKPushCallback(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}

	var stkCallback STKCallbackBody
	if err := json.Unmarshal(bodyBytes, &stkCallback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	log.Printf("STK Callback received: %+v\n", stkCallback)

	if stkCallback.Body.StkCallback.ResultCode == 0 {
		log.Println("STK Push was successful!")
	} else {
		log.Printf("STK Push failed: %s", stkCallback.Body.StkCallback.ResultDesc)
	}

	c.JSON(http.StatusOK, gin.H{"status": "Callback received"})
}

// Handle B2C Callback
func HandleB2CCallback(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read request body"})
		return
	}

	var b2cCallback B2CCallbackBody
	if err := json.Unmarshal(bodyBytes, &b2cCallback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	log.Printf("B2C Callback received: %+v\n", b2cCallback)

	if b2cCallback.Result.ResultCode == 0 {
		log.Println("B2C Payment was successful!")
	} else {
		log.Printf("B2C Payment failed: %s", b2cCallback.Result.ResultDesc)
	}

	c.JSON(http.StatusOK, gin.H{"status": "Callback received"})
}