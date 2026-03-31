package models

import (
	"time"

	"go-ecommerce/internal/database/sqlc"
)

// User models
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	Address   *string   `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Name     string  `json:"name" validate:"required"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Category models
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Brand models
type BrandResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateBrandRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type UpdateBrandRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Product models
type ProductResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description *string           `json:"description,omitempty"`
	Price       float64           `json:"price"`
	Stock       int               `json:"stock"`
	ImageURL    *string           `json:"image_url,omitempty"`
	CategoryID  string            `json:"category_id"`
	BrandID     string            `json:"brand_id"`
	Category    *CategoryResponse `json:"category,omitempty"`
	Brand       *BrandResponse    `json:"brand,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	ImageURL    *string `json:"image_url,omitempty"`
	CategoryID  string  `json:"category_id" validate:"required"`
	BrandID     string  `json:"brand_id" validate:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	CategoryID  *string  `json:"category_id,omitempty"`
	BrandID     *string  `json:"brand_id,omitempty"`
}

// Order models - use sqlc.OrderStatus for API
type OrderResponse struct {
	ID          string              `json:"id"`
	UserID      string              `json:"user_id"`
	TotalAmount float64             `json:"total_amount"`
	Status      sqlc.OrderStatus    `json:"status"`
	Items       []OrderItemResponse `json:"items,omitempty"`
	User        *UserResponse       `json:"user,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        string           `json:"id"`
	OrderID   string           `json:"order_id"`
	ProductID string           `json:"product_id"`
	Quantity  int              `json:"quantity"`
	Price     float64          `json:"price"`
	Product   *ProductResponse `json:"product,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status sqlc.OrderStatus `json:"status" validate:"required"`
}

// Helpers: sqlc types -> API response
func ToUserResponse(u *sqlc.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Phone:     u.Phone,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToCategoryResponse(c *sqlc.Category) CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func ToBrandResponse(b *sqlc.Brand) BrandResponse {
	return BrandResponse{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func ToProductResponseFromProduct(p *sqlc.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int(p.Stock),
		ImageURL:    p.ImageUrl,
		CategoryID:  p.CategoryID,
		BrandID:     p.BrandID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToProductResponseFromDetails(p *sqlc.ProductWithDetails) ProductResponse {
	resp := ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int(p.Stock),
		ImageURL:    p.ImageUrl,
		CategoryID:  p.CategoryID,
		BrandID:     p.BrandID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
	resp.Category = &CategoryResponse{
		ID: p.CatID, Name: p.CatName, Description: p.CatDescription,
		CreatedAt: p.CatCreatedAt, UpdatedAt: p.CatUpdatedAt,
	}
	resp.Brand = &BrandResponse{
		ID: p.BrandIDAlt, Name: p.BrandName, Description: p.BrandDescription,
		CreatedAt: p.BrandCreatedAt, UpdatedAt: p.BrandUpdatedAt,
	}
	return resp
}

func ToProductResponse(p *sqlc.Product) ProductResponse {
	return ToProductResponseFromProduct(p)
}

func ToOrderResponse(o *sqlc.Order, user *sqlc.User, items []sqlc.OrderItem, itemProducts map[string]*sqlc.ProductWithDetails) OrderResponse {
	resp := OrderResponse{
		ID:          o.ID,
		UserID:      o.UserID,
		TotalAmount: o.TotalAmount,
		Status:      o.Status,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
	if user != nil {
		u := ToUserResponse(user)
		resp.User = &u
	}
	resp.Items = make([]OrderItemResponse, len(items))
	for i, item := range items {
		resp.Items[i] = OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		}
		if itemProducts != nil {
			if prod, ok := itemProducts[item.ProductID]; ok {
				p := ToProductResponseFromDetails(prod)
				resp.Items[i].Product = &p
			}
		}
	}
	return resp
}
