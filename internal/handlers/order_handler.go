package handlers

import (
	"net/http"

	"go-ecommerce/internal/database/sqlc"
	apierrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/logger"
	"go-ecommerce/internal/middleware"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
	"go-ecommerce/internal/validator"
)

type OrderHandler struct {
	store *sqlc.Store
}

func NewOrderHandler(store *sqlc.Store) *OrderHandler {
	return &OrderHandler{store: store}
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
// @Router /api/v1/orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Unauthorized"))
		return
	}

	var req models.CreateOrderRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	productIDs := make([]string, 0, len(req.Items))
	productQty := make(map[string]int32, len(req.Items))
	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductID)
		productQty[item.ProductID] += int32(item.Quantity)
	}

	products, err := h.store.ListProductsWithDetailsByIDs(r.Context(), uniqueStrings(productIDs))
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to batch fetch products for order")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to validate order items"))
		return
	}

	productMap := make(map[string]*sqlc.Product, len(products))
	productDetailsMap := make(map[string]*sqlc.ProductWithDetails, len(products))
	for i := range products {
		p := products[i]
		productMap[p.ID] = &sqlc.Product{
			ID:         p.ID,
			Name:       p.Name,
			Description: p.Description,
			Price:      p.Price,
			Stock:      p.Stock,
			ImageUrl:   p.ImageUrl,
			CategoryID: p.CategoryID,
			BrandID:    p.BrandID,
			CreatedAt:  p.CreatedAt,
			UpdatedAt:  p.UpdatedAt,
		}
		productDetailsMap[p.ID] = &products[i]
	}

	for productID, qty := range productQty {
		product := productMap[productID]
		if product == nil {
			apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeNotFound, "Product not found: "+productID))
			return
		}
		if product.Stock < qty {
			apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInsufficientStock, "Insufficient stock for product: "+product.Name))
			return
		}
	}

	orderItems := make([]sqlc.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		orderItems = append(orderItems, sqlc.OrderItem{
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
	}

	order, err := h.store.CreateOrderWithItems(r.Context(), userID, sqlc.OrderStatusPending, productMap, orderItems)
	if err != nil {
		logger.Log.Error().Err(err).Str("user_id", userID).Msg("Failed to create order")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to create order"))
		return
	}

	user, _ := h.store.GetUserByID(r.Context(), userID)
	items, _ := h.store.ListOrderItemsByOrderID(r.Context(), order.ID)
	itemProducts := make(map[string]*sqlc.ProductWithDetails)
	for _, item := range items {
		if prod, ok := productDetailsMap[item.ProductID]; ok {
			itemProducts[item.ProductID] = prod
		}
	}
	logger.Log.Info().Str("order_id", order.ID).Str("user_id", userID).Float64("total", order.TotalAmount).Msg("Order created")
	resp := models.ToOrderResponse(order, user, items, itemProducts)
	apierrors.RespondWithJSON(w, http.StatusCreated, resp)
}

// GetOrders godoc
// @Summary Get all orders
// @Description Get all orders (User sees own orders, Admin sees all)
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.OrderResponse
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	role := middleware.GetUserRole(r)

	var orders []sqlc.Order
	var err error
	if role == "admin" {
		orders, err = h.store.ListOrdersAll(r.Context())
	} else {
		orders, err = h.store.ListOrdersByUserID(r.Context(), userID)
	}
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch orders")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to fetch orders"))
		return
	}

	responses := make([]models.OrderResponse, 0, len(orders))
	userMap, orderItemsMap, productMap := h.prefetchOrderRelations(r, orders)
	for i := range orders {
		responses = append(responses, models.ToOrderResponse(
			&orders[i],
			userMap[orders[i].UserID],
			orderItemsMap[orders[i].ID],
			productMap,
		))
	}
	apierrors.RespondWithJSON(w, http.StatusOK, responses)
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
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Order ID is required"))
		return
	}

	userID := middleware.GetUserID(r)
	role := middleware.GetUserRole(r)

	order, err := h.store.GetOrderByID(r.Context(), id)
	if err != nil || order == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Order not found"))
		return
	}
	if role != "admin" && order.UserID != userID {
		apierrors.RespondWithError(w, http.StatusForbidden, apierrors.New(apierrors.ErrCodeForbidden, "Access denied"))
		return
	}

	user, _ := h.store.GetUserByID(r.Context(), order.UserID)
	items, _ := h.store.ListOrderItemsByOrderID(r.Context(), order.ID)
	itemProducts := h.loadProductsByOrderItems(r, items)
	resp := models.ToOrderResponse(order, user, items, itemProducts)
	apierrors.RespondWithJSON(w, http.StatusOK, resp)
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
// @Router /api/v1/orders/{id}/status [put]
func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Order ID is required"))
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	order, err := h.store.UpdateOrderStatus(r.Context(), id, req.Status)
	if err != nil || order == nil {
		logger.Log.Error().Err(err).Str("order_id", id).Msg("Failed to update order status")
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Order not found"))
		return
	}

	user, _ := h.store.GetUserByID(r.Context(), order.UserID)
	items, _ := h.store.ListOrderItemsByOrderID(r.Context(), order.ID)
	itemProducts := h.loadProductsByOrderItems(r, items)
	logger.Log.Info().Str("order_id", id).Str("status", string(req.Status)).Msg("Order status updated")
	resp := models.ToOrderResponse(order, user, items, itemProducts)
	apierrors.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) prefetchOrderRelations(r *http.Request, orders []sqlc.Order) (map[string]*sqlc.User, map[string][]sqlc.OrderItem, map[string]*sqlc.ProductWithDetails) {
	userIDs := make([]string, 0, len(orders))
	orderIDs := make([]string, 0, len(orders))
	for i := range orders {
		userIDs = append(userIDs, orders[i].UserID)
		orderIDs = append(orderIDs, orders[i].ID)
	}

	users, _ := h.store.ListUsersByIDs(r.Context(), uniqueStrings(userIDs))
	userMap := make(map[string]*sqlc.User, len(users))
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	items, _ := h.store.ListOrderItemsByOrderIDs(r.Context(), uniqueStrings(orderIDs))
	orderItemsMap := make(map[string][]sqlc.OrderItem)
	for _, item := range items {
		orderItemsMap[item.OrderID] = append(orderItemsMap[item.OrderID], item)
	}

	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}
	products, _ := h.store.ListProductsWithDetailsByIDs(r.Context(), uniqueStrings(productIDs))
	productMap := make(map[string]*sqlc.ProductWithDetails, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	return userMap, orderItemsMap, productMap
}

func (h *OrderHandler) loadProductsByOrderItems(r *http.Request, items []sqlc.OrderItem) map[string]*sqlc.ProductWithDetails {
	productIDs := make([]string, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}
	products, _ := h.store.ListProductsWithDetailsByIDs(r.Context(), uniqueStrings(productIDs))
	productMap := make(map[string]*sqlc.ProductWithDetails, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}
	return productMap
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, s := range input {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
