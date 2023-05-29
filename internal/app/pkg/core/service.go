package core

import (
	"hw4/internal/app/pkg/dish"
	"hw4/internal/app/pkg/order"
	"hw4/internal/app/pkg/order_dish"
	"hw4/internal/app/pkg/session"
	"hw4/internal/app/pkg/user"
)

type CoreService struct {
	DishService      *dish.Service
	OrderService     *order.Service
	OrderDishService *order_dish.Service
	SessionService   *session.Service
	UserService      *user.Service
}

func NewService(dishService *dish.Service, orderService *order.Service, orderDishService *order_dish.Service, sessionService *session.Service, userService *user.Service) *CoreService {
	return &CoreService{
		DishService:      dishService,
		OrderService:     orderService,
		OrderDishService: orderDishService,
		SessionService:   sessionService,
		UserService:      userService,
	}
}
