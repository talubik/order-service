package generator

import (
	"time"
	"strconv"
	"myapp/src/internal/order_model"
	"github.com/brianvoe/gofakeit/v7"
)

// GenerateOrder создает заполненный Order со всеми вложенными структурами
func GenerateOrder(trackNumber string) model.Order {
	order := model.Order{
		OrderUID:          gofakeit.UUID(),
		TrackNumber:       trackNumber,
		Entry:             gofakeit.RandomString([]string{"WBIL", "TEST", "ENTRY"}),
		Delivery:          GenerateDelivery(),
		Payment:           GeneratePayment(),
		Items:             GenerateItems(1 + gofakeit.IntRange(1, 5), trackNumber), // 1-5 товаров
		Locale:            gofakeit.RandomString([]string{"en", "ru", "de"}),
		InternalSignature: gofakeit.LetterN(10),
		CustomerID:        gofakeit.LetterN(8),
		DeliveryService:   gofakeit.RandomString([]string{"meest", "russianpost", "dhl"}),
		Shardkey:          gofakeit.DigitN(2),
		SmID:              gofakeit.IntRange(1, 100),
		DateCreated:       gofakeit.Date().Format(time.RFC3339),
		OofShard:          gofakeit.DigitN(2),
	}

	return order
}

// GenerateDelivery создает заполненную структуру Delivery
func GenerateDelivery() model.Delivery {
	return model.Delivery{
		Name:    gofakeit.Name(),
		Phone:   gofakeit.Phone(),
		Zip:     gofakeit.Zip(),
		City:    gofakeit.City(),
		Address: gofakeit.Street() + ", " + gofakeit.DigitN(2),
		Region:  gofakeit.State(),
		Email:   gofakeit.Email(),
	}
}

// GeneratePayment создает заполненную структуру Payment
func GeneratePayment() model.Payment {
	amount := gofakeit.IntRange(1000, 100000)
	return model.Payment{
		Transaction:  gofakeit.UUID(),
		RequestID:    gofakeit.LetterN(6),
		Currency:     gofakeit.CurrencyShort(),
		Provider:     gofakeit.RandomString([]string{"wbpay", "sberpay", "tinkoff"}),
		Amount:       amount,
		PaymentDt:    gofakeit.Date().Unix(),
		Bank:         gofakeit.RandomString([]string{"sberbank", "tinkoff", "alfa"}),
		DeliveryCost: gofakeit.IntRange(100, 2000),
		GoodsTotal:   amount - gofakeit.IntRange(100, 500),
		CustomFee:    gofakeit.IntRange(0, 100),
	}
}

// GenerateItems создает слот Items
func GenerateItems(count int, trackNumber string) []model.Item {
	items := make([]model.Item, count)
	for i := 0; i < count; i++ {
		price := gofakeit.IntRange(100, 5000)
		quantity := gofakeit.IntRange(1, 5)
		
		items[i] = model.Item{
			ChrtID:     gofakeit.IntRange(1000000, 9999999),
			TrackNum:   trackNumber,
			Price:      price,
			Rid:        gofakeit.UUID(),
			Name:       gofakeit.ProductName(),
			Sale:       gofakeit.IntRange(0, 50),
			Size:       gofakeit.RandomString([]string{"S", "M", "L", "XL", "36", "38", "40"}),
			TotalPrice: price * quantity,
			NmID:       gofakeit.IntRange(100000, 999999),
			Brand:      gofakeit.Company(),
			Status:     gofakeit.IntRange(1, 5),
		}
	}
	return items
}

// GenerateOrders создает несколько заказов
func GenerateOrders(count int) []model.Order {
	orders := make([]model.Order, count)
	for i := 0; i < count; i++ {
		tn := "WB" + strconv.Itoa(i)
		orders[i] = GenerateOrder(tn)
	}
	return orders
}

