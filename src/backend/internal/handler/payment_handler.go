package handler

import (
	"net/http"

	"unihub-workshop/internal/model"
	"unihub-workshop/internal/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(ps *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: ps}
}

// InitiatePayment starts payment flow for a paid workshop registration
func (h *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	registrationID := getURLParam(r, "registrationId")

	payment, checkoutURL, err := h.paymentService.InitiatePayment(r.Context(), registrationID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"payment":      payment,
			"checkout_url": checkoutURL,
		},
	})
}

// Webhook handles payment gateway callbacks
func (h *PaymentHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	var req model.PaymentWebhookRequest
	if err := decodeJSON(r, &req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid webhook payload")
		return
	}

	if err := h.paymentService.HandleWebhook(r.Context(), &req); err != nil {
		if err.Error() == "invalid signature" {
			errorResponse(w, http.StatusUnauthorized, "Invalid signature")
			return
		}
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model.APIResponse{Success: true, Message: "Webhook processed"})
}

// GetCircuitBreakerStatus returns the current state of the payment circuit breaker
func (h *PaymentHandler) GetCircuitBreakerStatus(w http.ResponseWriter, r *http.Request) {
	state := h.paymentService.GetCircuitBreakerState()
	writeJSON(w, http.StatusOK, model.APIResponse{
		Success: true,
		Data:    map[string]string{"circuit_breaker_state": state},
	})
}
