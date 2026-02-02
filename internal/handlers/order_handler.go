package handlers

import (
	"net/http"

	"go-ecommerce/internal/database"
	"go-ecommerce/internal/database/sqlc"
	"go-ecommerce/internal/middleware"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order (User only)
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param order body models.CreateOrderRequest true "Order data"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Router /api/orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateOrderRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var totalAmount float64
	for _, item := range req.Items {
		product, err := database.Queries.GetProductByID(r.Context(), item.ProductID)
		if err != nil || product == nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Product not found: "+item.ProductID)
			return
		}
		if product.Stock < int32(item.Quantity) {
			utils.RespondWithError(w, http.StatusBadRequest, "Insufficient stock for product: "+product.Name)
			return
		}
		totalAmount += product.Price * float64(item.Quantity)
		_, err = database.Queries.UpdateProductStock(r.Context(), item.ProductID, product.Stock-int32(item.Quantity))
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product stock")
			return
		}
	}

	order, err := database.Queries.CreateOrder(r.Context(), userID, totalAmount, sqlc.OrderStatusPending)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	for _, item := range req.Items {
		product, _ := database.Queries.GetProductByID(r.Context(), item.ProductID)
		if product == nil {
			continue
		}
		_, err = database.Queries.CreateOrderItem(r.Context(), order.ID, item.ProductID, int32(item.Quantity), product.Price)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order item")
			return
		}
	}

	// Build response: order + user + items with product details
	user, _ := database.Queries.GetUserByID(r.Context(), userID)
	items, _ := database.Queries.ListOrderItemsByOrderID(r.Context(), order.ID)
	itemProducts := make(map[string]*sqlc.ProductWithDetails)
	for _, item := range items {
		prod, _ := database.Queries.GetProductWithDetails(r.Context(), item.ProductID)
		if prod != nil {
			itemProducts[item.ProductID] = prod
		}
	}
	resp := models.ToOrderResponse(order, user, items, itemProducts)
	utils.RespondWithJSON(w, http.StatusCreated, resp)
}

// GetOrders godoc
// @Summary Get all orders
// @Description Get all orders (User sees own orders, Admin sees all)
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.OrderResponse
// @Router /api/orders [get]
func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	role := middleware.GetUserRole(r)

	var orders []sqlc.Order
	var err error
	if role == "admin" {
		orders, err = database.Queries.ListOrdersAll(r.Context())
	} else {
		orders, err = database.Queries.ListOrdersByUserID(r.Context(), userID)
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch orders")
		return
	}

	responses := make([]models.OrderResponse, 0, len(orders))
	for i := range orders {
		user, _ := database.Queries.GetUserByID(r.Context(), orders[i].UserID)
		items, _ := database.Queries.ListOrderItemsByOrderID(r.Context(), orders[i].ID)
		itemProducts := make(map[string]*sqlc.ProductWithDetails)
		for _, item := range items {
			prod, _ := database.Queries.GetProductWithDetails(r.Context(), item.ProductID)
			if prod != nil {
				itemProducts[item.ProductID] = prod
			}
		}
		responses = append(responses, models.ToOrderResponse(&orders[i], user, items, itemProducts))
	}
	utils.RespondWithJSON(w, http.StatusOK, responses)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Get a specific order by ID
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.OrderResponse
// @Failure 404 {object} map[string]string
// @Router /api/orders/{id} [get]
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	userID := middleware.GetUserID(r)
	role := middleware.GetUserRole(r)

	order, err := database.Queries.GetOrderByID(r.Context(), id)
	if err != nil || order == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}
	if role != "admin" && order.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	user, _ := database.Queries.GetUserByID(r.Context(), order.UserID)
	items, _ := database.Queries.ListOrderItemsByOrderID(r.Context(), order.ID)
	itemProducts := make(map[string]*sqlc.ProductWithDetails)
	for _, item := range items {
		prod, _ := database.Queries.GetProductWithDetails(r.Context(), item.ProductID)
		if prod != nil {
			itemProducts[item.ProductID] = prod
		}
	}
	resp := models.ToOrderResponse(order, user, items, itemProducts)
	utils.RespondWithJSON(w, http.StatusOK, resp)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update order status (Admin only)
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "Order status"
// @Success 200 {object} models.OrderResponse
// @Failure 404 {object} map[string]string
// @Router /api/orders/{id}/status [put]
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order, err := database.Queries.UpdateOrderStatus(r.Context(), id, req.Status)
	if err != nil || order == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	user, _ := database.Queries.GetUserByID(r.Context(), order.UserID)
	items, _ := database.Queries.ListOrderItemsByOrderID(r.Context(), order.ID)
	itemProducts := make(map[string]*sqlc.ProductWithDetails)
	for _, item := range items {
		prod, _ := database.Queries.GetProductWithDetails(r.Context(), item.ProductID)
		if prod != nil {
			itemProducts[item.ProductID] = prod
		}
	}
	resp := models.ToOrderResponse(order, user, items, itemProducts)
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
