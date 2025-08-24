package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"github.com/space-wanderer/microservices/notification/internal/client/http"
	"github.com/space-wanderer/microservices/notification/internal/model"
)

const chatID = 236673056

//go:embed templates/paid_notification.tmpl
//go:embed templates/assembled_notification.tmpl
var templates embed.FS

type orderPaidTemplateData struct {
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

type shipAssembledTemplateData struct {
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}

var notificaitonTemplate = template.Must(template.ParseFS(
	templates,
	"templates/paid_notification.tmpl",
	"templates/assembled_notification.tmpl",
))

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, uuid string, event model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(uuid, event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) buildOrderPaidMessage(uuid string, event model.OrderPaidEvent) (string, error) {
	data := orderPaidTemplateData{
		OrderUUID:       uuid,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}

	var buf bytes.Buffer
	err := notificaitonTemplate.ExecuteTemplate(&buf, "paid_notification", data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendShipAssembledNotification(ctx context.Context, uuid string, event model.ShipAssembledEvent) error {
	message, err := s.buildShipAssembledMessage(uuid, event)
	if err != nil {
		return err
	}

	return s.telegramClient.SendMessage(ctx, chatID, message)
}

func (s *service) buildShipAssembledMessage(uuid string, event model.ShipAssembledEvent) (string, error) {
	data := shipAssembledTemplateData{
		OrderUUID:    uuid,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := notificaitonTemplate.ExecuteTemplate(&buf, "assembled_notification", data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
