package routes

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/BatJoz21/my-online-shop-go-api/models"
	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func initiatePayment(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("orderID"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	order, err := models.GetOrderForPayment(id)
	if err == sql.ErrNoRows {
		context.JSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if order.Status != "pending" {
		context.JSON(http.StatusConflict, gin.H{"message": "Order is not waiting for payment"})
		return
	}

	pending, err := models.IsExistingPaymentPending(order.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if pending {
		context.JSON(http.StatusConflict, gin.H{"message": "A payment is already in progress for this order"})
		return
	}

	var snapClient snap.Client
	snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order.OrderNumber,
			GrossAmt: order.TotalAmount.IntPart(),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: order.CustomerName,
			Email: order.CustomerEmail,
		},
		Callbacks: &snap.Callbacks{
			Finish: os.Getenv("CI4_BASE_URL") + "orders/" + strconv.FormatInt(id, 10) + "/payment/result",
		},
	}

	snapResponse, midtransErr := snapClient.CreateTransaction(req)
	if midtransErr != nil {
		context.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to initiate payment",
			"detai": midtransErr.GetMessage()})
		return
	}

	rawResponse, _ := json.Marshal(snapResponse)

	payment := models.Payment{
		OrderID:     order.ID,
		Provider:    "midtrans",
		Amount:      order.TotalAmount,
		Status:      "pending",
		RawResponse: rawResponse,
	}

	if err := payment.Save(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"redirect_url": snapResponse.RedirectURL,
		"token":        snapResponse.Token,
	})
}

func handlePaymentWebhook(context *gin.Context) {
	bodyBytes, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read request body", "error": err.Error()})
		return
	}

	var notif models.MidtransNotification
	if err := json.Unmarshal(bodyBytes, &notif); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload", "error": err.Error()})
		return
	}

	if !isValidMidtransNotification(&notif) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid signature"})
		return
	}

	newStatus := mapMidtransStatus(notif.TransactionStatus, notif.FraudStatus)

	err = models.UpdateFromWebhook(notif.OrderID, newStatus, notif.TransactionID, bodyBytes)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update payment", "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Notification processed"})
}

func isValidMidtransNotification(notif *models.MidtransNotification) bool {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	raw := notif.OrderID + notif.StatusCode + notif.GrossAmount + serverKey

	hash := sha512.Sum512([]byte(raw))
	expectedSignature := hex.EncodeToString(hash[:])

	return expectedSignature == notif.SignatureKey
}

func mapMidtransStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "accept" {
			return "success"
		}
		return "pending"
	case "settlement":
		return "success"
	case "pending":
		return "pending"
	case "deny", "cancel":
		return "failed"
	case "expire":
		return "expired"
	default:
		return "pending"
	}
}
