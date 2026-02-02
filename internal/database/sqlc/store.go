package sqlc

import (
	"context"
	"database/sql"
)

// Store provides all database operations. When sqlc is used, run `sqlc generate` and this file can be replaced by generated code.
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Admin
func (s *Store) GetAdminByID(ctx context.Context, id string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx, `SELECT id, email, password, name, created_at, updated_at FROM admins WHERE id = $1`, id).Scan(
		&a.ID, &a.Email, &a.Password, &a.Name, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Store) GetAdminByEmail(ctx context.Context, email string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx, `SELECT id, email, password, name, created_at, updated_at FROM admins WHERE email = $1`, email).Scan(
		&a.ID, &a.Email, &a.Password, &a.Name, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Store) CreateAdmin(ctx context.Context, email, password, name string) (*Admin, error) {
	var a Admin
	err := s.db.QueryRowContext(ctx, `INSERT INTO admins (email, password, name) VALUES ($1, $2, $3) RETURNING id, email, password, name, created_at, updated_at`,
		email, password, name).Scan(&a.ID, &a.Email, &a.Password, &a.Name, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// User
func (s *Store) GetUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `SELECT id, email, password, name, phone, address, created_at, updated_at FROM users WHERE id = $1`, id).Scan(
		&u.ID, &u.Email, &u.Password, &u.Name, &u.Phone, &u.Address, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `SELECT id, email, password, name, phone, address, created_at, updated_at FROM users WHERE email = $1`, email).Scan(
		&u.ID, &u.Email, &u.Password, &u.Name, &u.Phone, &u.Address, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) CreateUser(ctx context.Context, email, password, name string, phone, address *string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `INSERT INTO users (email, password, name, phone, address) VALUES ($1, $2, $3, $4, $5) RETURNING id, email, password, name, phone, address, created_at, updated_at`,
		email, password, name, phone, address).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Phone, &u.Address, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UpdateUser(ctx context.Context, id string, name *string, phone, address *string) (*User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `UPDATE users SET name = COALESCE($2, name), phone = COALESCE($3, phone), address = COALESCE($4, address), updated_at = now() WHERE id = $1 RETURNING id, email, password, name, phone, address, created_at, updated_at`,
		id, name, phone, address).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.Phone, &u.Address, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Category
func (s *Store) GetCategoryByID(ctx context.Context, id string) (*Category, error) {
	var c Category
	err := s.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, updated_at FROM categories WHERE id = $1`, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) GetCategoryByName(ctx context.Context, name string) (*Category, error) {
	var c Category
	err := s.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, updated_at FROM categories WHERE name = $1`, name).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description, created_at, updated_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (s *Store) CreateCategory(ctx context.Context, name string, description *string) (*Category, error) {
	var c Category
	err := s.db.QueryRowContext(ctx, `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at`,
		name, description).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) UpdateCategory(ctx context.Context, id string, name *string, description *string) (*Category, error) {
	var c Category
	err := s.db.QueryRowContext(ctx, `UPDATE categories SET name = COALESCE($2, name), description = $3, updated_at = now() WHERE id = $1 RETURNING id, name, description, created_at, updated_at`,
		id, name, description).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) DeleteCategory(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}

// Brand
func (s *Store) GetBrandByID(ctx context.Context, id string) (*Brand, error) {
	var b Brand
	err := s.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, updated_at FROM brands WHERE id = $1`, id).Scan(
		&b.ID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Store) GetBrandByName(ctx context.Context, name string) (*Brand, error) {
	var b Brand
	err := s.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, updated_at FROM brands WHERE name = $1`, name).Scan(
		&b.ID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Store) ListBrands(ctx context.Context) ([]Brand, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description, created_at, updated_at FROM brands ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Brand
	for rows.Next() {
		var b Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (s *Store) CreateBrand(ctx context.Context, name string, description *string) (*Brand, error) {
	var b Brand
	err := s.db.QueryRowContext(ctx, `INSERT INTO brands (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at`,
		name, description).Scan(&b.ID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Store) UpdateBrand(ctx context.Context, id string, name *string, description *string) (*Brand, error) {
	var b Brand
	err := s.db.QueryRowContext(ctx, `UPDATE brands SET name = COALESCE($2, name), description = $3, updated_at = now() WHERE id = $1 RETURNING id, name, description, created_at, updated_at`,
		id, name, description).Scan(&b.ID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Store) DeleteBrand(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM brands WHERE id = $1`, id)
	return err
}

// Product
func (s *Store) GetProductByID(ctx context.Context, id string) (*Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx, `SELECT id, name, description, price, stock, image_url, category_id, brand_id, created_at, updated_at FROM products WHERE id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageUrl, &p.CategoryID, &p.BrandID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) GetProductWithDetails(ctx context.Context, id string) (*ProductWithDetails, error) {
	var p ProductWithDetails
	err := s.db.QueryRowContext(ctx, `SELECT p.id, p.name, p.description, p.price, p.stock, p.image_url, p.category_id, p.brand_id, p.created_at, p.updated_at,
		c.id as cat_id, c.name as cat_name, c.description as cat_description, c.created_at as cat_created_at, c.updated_at as cat_updated_at,
		b.id as brand_id, b.name as brand_name, b.description as brand_description, b.created_at as brand_created_at, b.updated_at as brand_updated_at
		FROM products p JOIN categories c ON p.category_id = c.id JOIN brands b ON p.brand_id = b.id WHERE p.id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageUrl, &p.CategoryID, &p.BrandID, &p.CreatedAt, &p.UpdatedAt,
		&p.CatID, &p.CatName, &p.CatDescription, &p.CatCreatedAt, &p.CatUpdatedAt,
		&p.BrandIDAlt, &p.BrandName, &p.BrandDescription, &p.BrandCreatedAt, &p.BrandUpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) ListProducts(ctx context.Context, categoryID, brandID *string) ([]ProductWithDetails, error) {
	query := `SELECT p.id, p.name, p.description, p.price, p.stock, p.image_url, p.category_id, p.brand_id, p.created_at, p.updated_at,
		c.id as cat_id, c.name as cat_name, c.description as cat_description, c.created_at as cat_created_at, c.updated_at as cat_updated_at,
		b.id as brand_id, b.name as brand_name, b.description as brand_description, b.created_at as brand_created_at, b.updated_at as brand_updated_at
		FROM products p JOIN categories c ON p.category_id = c.id JOIN brands b ON p.brand_id = b.id
		WHERE ($1::uuid IS NULL OR p.category_id = $1) AND ($2::uuid IS NULL OR p.brand_id = $2) ORDER BY p.name`
	rows, err := s.db.QueryContext(ctx, query, categoryID, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ProductWithDetails
	for rows.Next() {
		var p ProductWithDetails
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageUrl, &p.CategoryID, &p.BrandID, &p.CreatedAt, &p.UpdatedAt,
			&p.CatID, &p.CatName, &p.CatDescription, &p.CatCreatedAt, &p.CatUpdatedAt,
			&p.BrandIDAlt, &p.BrandName, &p.BrandDescription, &p.BrandCreatedAt, &p.BrandUpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (s *Store) CreateProduct(ctx context.Context, name string, description *string, price float64, stock int32, imageURL *string, categoryID, brandID string) (*Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx, `INSERT INTO products (name, description, price, stock, image_url, category_id, brand_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, description, price, stock, image_url, category_id, brand_id, created_at, updated_at`,
		name, description, price, stock, imageURL, categoryID, brandID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageUrl, &p.CategoryID, &p.BrandID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) UpdateProduct(ctx context.Context, id string, name *string, description *string, price *float64, stock *int32, imageURL *string, categoryID, brandID *string) (*Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx, `UPDATE products SET name = COALESCE($2, name), description = $3, price = COALESCE($4, price), stock = COALESCE($5, stock), image_url = $6, category_id = COALESCE($7, category_id), brand_id = COALESCE($8, brand_id), updated_at = now() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, category_id, brand_id, created_at, updated_at`,
		id, name, description, price, stock, imageURL, categoryID, brandID).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageUrl, &p.CategoryID, &p.BrandID, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) UpdateProductStock(ctx context.Context, id string, stock int32) (*Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx, `UPDATE products SET stock = $2, updated_at = now() WHERE id = $1 RETURNING id, name, description, price, stock, image_url, category_id, brand_id, created_at, updated_at`, id, stock).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageUrl, &p.CategoryID, &p.BrandID, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) DeleteProduct(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

// Order
func (s *Store) CreateOrder(ctx context.Context, userID string, totalAmount float64, status OrderStatus) (*Order, error) {
	var o Order
	err := s.db.QueryRowContext(ctx, `INSERT INTO orders (user_id, total_amount, status) VALUES ($1, $2, $3) RETURNING id, user_id, total_amount, status, created_at, updated_at`,
		userID, totalAmount, status).Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Store) GetOrderByID(ctx context.Context, id string) (*Order, error) {
	var o Order
	err := s.db.QueryRowContext(ctx, `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`, id).Scan(
		&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Store) ListOrdersByUserID(ctx context.Context, userID string) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, rows.Err()
}

func (s *Store) ListOrdersAll(ctx context.Context) ([]Order, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, rows.Err()
}

func (s *Store) UpdateOrderStatus(ctx context.Context, id string, status OrderStatus) (*Order, error) {
	var o Order
	err := s.db.QueryRowContext(ctx, `UPDATE orders SET status = $2, updated_at = now() WHERE id = $1 RETURNING id, user_id, total_amount, status, created_at, updated_at`, id, status).Scan(
		&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// OrderItem
func (s *Store) CreateOrderItem(ctx context.Context, orderID, productID string, quantity int32, price float64) (*OrderItem, error) {
	var oi OrderItem
	err := s.db.QueryRowContext(ctx, `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id, order_id, product_id, quantity, price, created_at`,
		orderID, productID, quantity, price).Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity, &oi.Price, &oi.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &oi, nil
}

func (s *Store) ListOrderItemsByOrderID(ctx context.Context, orderID string) ([]OrderItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, order_id, product_id, quantity, price, created_at FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []OrderItem
	for rows.Next() {
		var oi OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Quantity, &oi.Price, &oi.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, oi)
	}
	return list, rows.Err()
}
