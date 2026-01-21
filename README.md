# ecommerce-w-go
Clean Architecture

## Models
### Product
- id
- name
- price
- quantity
- created_at
- updated_at


### Customer
- id
- first_name
- last_name
- created_at
- updated_at

### Order
- id
- status
- created_at
- updated_at


### Order items
- id
- product_id
- order_id
- quantity
- price
- created_at
- updated_at

## APIS
- prefix: /v1
- GET /v1/products?name&limit=20&offset=0
- POST /v1/orders

