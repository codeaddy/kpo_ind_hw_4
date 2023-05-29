package order_processing

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"hw4/internal/app/pkg/core"
	"hw4/internal/app/pkg/dish"
	"hw4/internal/app/pkg/order"
	"hw4/internal/app/pkg/order_dish"
	"net/http"
	"strconv"
)

type OrderProcessingService struct {
	core *core.CoreService
}

func NewService(core *core.CoreService) *OrderProcessingService {
	return &OrderProcessingService{core: core}
}

type createOrderInput struct {
	UserID          int         `json:"user_id"`
	AmountOfDishes  int         `json:"amount_of_dishes"`
	Dishes          []dish.Dish `json:"dishes"`
	SpecialRequests string      `json:"special_requests"`
	Status          string      `json:"status"`
}

type createOrderResponse struct {
	Success string `json:"success"`
}

// CreateOrder godoc
// @Summary		CreateOrder
// @Description	Creating order
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		createOrderInput	true	"UserID, amountOfDishes, dishes array, specialRequests and status"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/create-order [post]
func (s *OrderProcessingService) CreateOrder(c *gin.Context) {
	var input createOrderInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	_, err := s.core.UserService.GetById(c, input.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "there is no user with such id"})
		return
	}

	if len(input.Dishes) != input.AmountOfDishes {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "length of dishes array is not equal to amountOfDishes"})
		return
	}

	if input.Status != "pending" && input.Status != "processing" && input.Status != "finished" && input.Status != "cancelled" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid status of order"})
		return
	}

	ok := true
	for _, dish := range input.Dishes {
		_, err := s.core.DishService.GetById(c, dish.ID)
		if err != nil {
			ok = false
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "there is no such dish named: " + dish.Name})
			break
		}
	}
	if !ok {
		return
	}

	orderID, err := s.core.OrderService.Create(c, order.Order{
		UserID:          input.UserID,
		Status:          input.Status,
		SpecialRequests: input.SpecialRequests,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok = true
	for _, dish := range input.Dishes {
		_, err := s.core.OrderDishService.Create(c, order_dish.OrderDish{
			OrderID:  orderID,
			DishID:   dish.ID,
			Quantity: dish.Quantity,
			Price:    dish.Price,
		})
		if err != nil {
			ok = false
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			break
		}
	}
	if !ok {
		return
	}

	c.IndentedJSON(http.StatusOK, createOrderResponse{Success: "order successfully created, id = " + strconv.Itoa(orderID)})
}

type updateOrderInput struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

type updateOrderResponse struct {
	Success string `json:"success"`
}

// UpdateOrder godoc
// @Summary		UpdateOrder
// @Description	Updating order
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		updateOrderInput	true	"OrderID and new order status"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/update-order [put]
func (s *OrderProcessingService) UpdateOrder(c *gin.Context) {
	var input updateOrderInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	_, err := s.core.OrderService.GetById(c, input.OrderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "there is no order with such id"})
		return
	}

	if input.Status != "pending" && input.Status != "processing" && input.Status != "finished" && input.Status != "cancelled" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid status of order"})
		return
	}

	_, err = s.core.OrderService.GetById(c, input.OrderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updateOrderResponse{Success: "order successfully updated"})
}

type getOrderInput struct {
	OrderID int `json:"order_id"`
}

// GetOrder godoc
// @Summary		GetOrder
// @Description	Getting order
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		getOrderInput	true	"OrderID to get"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/get-order [post]
func (s *OrderProcessingService) GetOrder(c *gin.Context) {
	var input getOrderInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	order, err := s.core.OrderService.GetById(c, input.OrderID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, order)
}

func (s *OrderProcessingService) checkIfUserIsManagerByToken(c *gin.Context, token string) bool {
	session, err := s.core.SessionService.GetByToken(c, token)
	if err != nil {
		return false
	}

	user, err := s.core.UserService.GetById(c, session.UserID)
	if err != nil {
		return false
	}

	return user.Role == "manager"
}

type getDishInput struct {
	DishID int    `json:"dish_id"`
	Token  string `json:"token"`
}

// GetDish godoc
// @Summary		GetDish
// @Description	Getting dish
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		getDishInput	true	"DishID and user token"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		423		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/get-dish [post]
func (s *OrderProcessingService) GetDish(c *gin.Context) {
	var input getDishInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	if !s.checkIfUserIsManagerByToken(c, input.Token) {
		c.AbortWithStatusJSON(http.StatusLocked, gin.H{"error": "dishes management is accessible for managers only"})
		return
	}

	dish, err := s.core.DishService.GetById(c, input.DishID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, dish)
}

type addDishInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	IsAvailable bool    `json:"is_available"`
	UserToken   string  `json:"user_token"`
}

type addDishResponse struct {
	Id int `json:"id"`
}

// AddDish godoc
// @Summary		AddDish
// @Description	Adding dish
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		addDishInput	true	"Name, description, price, quantity, isAvailable and userToken"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		423		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/add-dish [post]
func (s *OrderProcessingService) AddDish(c *gin.Context) {
	var input addDishInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	if !s.checkIfUserIsManagerByToken(c, input.UserToken) {
		c.AbortWithStatusJSON(http.StatusLocked, gin.H{"error": "dishes management is accessible for managers only"})
		return
	}

	id, err := s.core.DishService.Create(c, dish.Dish{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		IsAvailable: input.IsAvailable,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.IndentedJSON(http.StatusOK, addDishResponse{Id: id})
}

type updateDishInput struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	IsAvailable bool    `json:"is_available"`
	UserToken   string  `json:"user_token"`
}

type updateDishResponse struct {
	Success string `json:"success"`
}

// UpdateDish godoc
// @Summary		UpdateDish
// @Description	Updating dish
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		updateDishInput	true	"id, name, description, price, quantity, isAvailable and userToken"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		423		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/update-dish [put]
func (s *OrderProcessingService) UpdateDish(c *gin.Context) {
	var input updateDishInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	if !s.checkIfUserIsManagerByToken(c, input.UserToken) {
		c.AbortWithStatusJSON(http.StatusLocked, gin.H{"error": "dishes management is accessible for managers only"})
		return
	}

	_, err := s.core.DishService.Update(c, dish.Dish{
		ID:          input.Id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		IsAvailable: input.IsAvailable,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, updateDishResponse{Success: "updated"})
}

type deleteDishInput struct {
	DishID int    `json:"dish_id"`
	Token  string `json:"token"`
}

type deleteDishResponse struct {
	Success string `json:"success"`
}

// DeleteDish godoc
// @Summary		DeleteDish
// @Description	Deleting dish
// @Tags			order-processing
// @Accept			json
// @Produce		json
// @Param			input	body		deleteDishInput	true	"dishID and userToken"
// @Success		200		{object}	int
// @Failure		400		{object}	string
// @Failure		423		{object}	string
// @Failure		500		{object}	string
// @Router			/order-processing/delete-dish [post]
func (s *OrderProcessingService) DeleteDish(c *gin.Context) {
	fmt.Println("check1")
	var input deleteDishInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "error when unmarshalling input"})
		return
	}

	tmp, _ := json.Marshal(input)
	fmt.Println(string(tmp))

	if !s.checkIfUserIsManagerByToken(c, input.Token) {
		c.AbortWithStatusJSON(http.StatusLocked, gin.H{"error": "dishes management is accessible for managers only"})
		return
	}

	_, err := s.core.DishService.Delete(c, input.DishID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.IndentedJSON(http.StatusOK, deleteDishResponse{Success: "deleted"})
}

// GetMenu godoc
// @Summary		GetMenu
// @Description	Getting menu
// @Tags			order-processing
// @Produce		json
// @Success		200	{object}	int
// @Failure		500	{object}	string
// @Router			/order-processing/get-menu [get]
func (s *OrderProcessingService) GetMenu(c *gin.Context) {
	dishes, err := s.core.DishService.GetAll(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, dishes)
}
