# MyOnlineShop API

A REST API for a practice e-commerce project, built with **Go** and **Gin**, backed by **MariaDB**. Handles product catalog, cart, orders, reviews, and payments via the **Midtrans Snap** payment gateway. Designed to be consumed by a separate CodeIgniter 4 frontend.

## Tech Stack

- **Go** + [Gin](https://github.com/gin-gonic/gin) — HTTP routing and middleware
- **MariaDB** — via `database/sql` and `go-sql-driver/mysql` (no ORM)
- **JWT** — stateless auth with refresh tokens
- **Midtrans Snap** — sandbox payment integration ([midtrans-go](https://github.com/midtrans/midtrans-go) SDK)
- **Docker** — multi-stage build, orchestrated via Docker Compose

## Features

- Public product catalog with category, search, price, and sort filters
- Customer auth (register/login/refresh/logout) with JWT + refresh tokens
- Cart management
- Order creation from cart, with stock deduction and human-readable order numbers
- Payment initiation and **signature-verified, idempotent** webhook handling for Midtrans
- Product reviews, gated on a completed order containing the product
- Merchant-role product, variant, and order management
- Role-based middleware (`Authenticate`, `MerchantMiddleware`)

## Project Structure

```
.
├── database/       # DB connection + table creation/seeding
├── middlewares/     # Auth and role-based middleware
├── models/          # DB queries, organized by domain (MVC-style)
├── routes/          # Route registration + handlers
├── main.go
├── go.mod
└── Dockerfile
```

## Environment Variables

| Variable | Description | Example |
|---|---|---|
| `DB_HOST` | Database hostname | `mysql` (Docker) / `127.0.0.1` (local) |
| `DB_PORT` | Database port | `3306` |
| `DB_USER` | Database user | |
| `DB_PASSWORD` | Database password | |
| `DB_NAME` | Database name | `myonlineshop` |
| `MIDTRANS_SERVER_KEY` | Midtrans **sandbox** server key | `SB-Mid-server-...` |
| `CI4_BASE_URL` | Public base URL of the frontend, used for Midtrans's post-payment redirect | `http://myonlineshop.localhost` |

For local (non-Docker) development, create a `.env` file in the project root with these keys — the app falls back to real environment variables if `.env` isn't found, so both approaches work.

## Running Locally (without Docker)

```bash
go mod download
go run main.go
```

Requires a running MariaDB instance reachable at the configured `DB_HOST`/`DB_PORT`.

## Running with Docker

This service is designed to run as part of the full stack via Docker Compose (see the deployment repo/directory, which also includes the MariaDB and CI4 web services). To build just this image standalone:

```bash
docker build -t mos-api .
docker run -p 8080:8080 --env-file .env mos-api
```

## Payment Integration Notes

- Uses **Midtrans Snap** in sandbox mode. Payment initiation (`POST /orders/:orderID/payment`) creates a `pending` payment row and returns a `redirect_url` for the customer.
- The webhook endpoint (`POST /payments/webhook`) is public but protected by **SHA512 signature verification** against the Midtrans server key — never trust this route based on network exposure alone.
- Webhook handling is **idempotent**: once an order is marked `paid`, later out-of-order or duplicate notifications (e.g. a stale `expired` webhook arriving after a `success` one) are safely ignored rather than overwriting a confirmed payment.
- For local webhook testing, expose port 8080 via a tunnel tool (e.g. [ngrok](https://ngrok.com)) and register the resulting public URL + `/payments/webhook` as the Payment Notification URL in the Midtrans sandbox dashboard.

## API Routes

Base path: none (routes are registered at root, e.g. `/products`, `/orders`).

### Auth
```
POST   /register
POST   /login
POST   /logout
POST   /refresh
```

### Public
```
GET    /products
GET    /products/:id
GET    /products/:id/image
GET    /products/:id/variants
GET    /products/:id/reviews
GET    /categories
```

### Customer (auth required)
```
GET    /cart
GET    /cart/total
POST   /addToCart
GET    /carts
PUT    /carts/:id
DELETE /carts/:id
DELETE /carts

GET    /orders
POST   /orders
GET    /orders/:orderID
PUT    /orders/:orderID/populate
PUT    /orders/:orderID/complete
DELETE /orders/:orderID
POST   /orders/:orderID/items
GET    /orders/:orderID/items
DELETE /orders/:orderID/items/:orderItemID

POST   /orders/:orderID/payment
GET    /orders/:orderID/payment

POST   /products/:id/reviews
```

### Merchant (auth + merchant role required)
```
GET    /products/all
POST   /products
PUT    /products/:id
PUT    /products/:id/restore
DELETE /products/:id
DELETE /products/:id/delete

POST   /products/:id/variants
GET    /products/:id/variants/:variant_id
PUT    /products/:id/variants/:variant_id
PUT    /products/:id/variants/:variant_id/stock
DELETE /products/:id/variants/:variant_id

GET    /merchant/orders
GET    /merchant/orders/:orderID
PUT    /merchant/orders/:orderID
```

### Webhook (public, signature-verified)
```
POST   /payments/webhook
```

## Notes

This is a learning/practice project. The Midtrans integration runs against the **sandbox** environment only — no real payments are processed.
