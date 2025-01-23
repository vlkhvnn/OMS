package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

type PaymentHTTPHandler struct {
	channel *amqp.Channel
}

func NewPaymentHTTPHandler(channel *amqp.Channel) *PaymentHTTPHandler {
	return &PaymentHTTPHandler{channel}
}

func (h *PaymentHTTPHandler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)
}

func (h *PaymentHTTPHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointStripeSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	if event.Type == stripe.EventTypeCheckoutSessionCompleted || event.Type == stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded {
		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if cs.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("Payment for Checkout Session %v succeeded!", cs.ID)
		}
	}

	w.WriteHeader(http.StatusOK)
}
