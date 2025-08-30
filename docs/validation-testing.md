# TC-4 Validation Testing

## Introduction

This document describes how to perform validation testing on the TC-2 API. Validation testing is the process of ensuring that the API meets the requirements specified in the [Software Requirements Specification](./tc2-spec.pdf) document. This document will guide you through the process of testing the API to ensure that it meets the requirements.

After following the steps in the [Readme](../README.md) file, you should have the API running on your local machine.

> [!IMPORTANT]
> Use `http://localhost:8080` as the base URL for the API if you are running it locally via Docker Compose (`make compose-up`).  
> Alternatively, you can use `http://localhost` if you are running it locally via Kubernetes (`make k8s-apply`).  

## Test Cases

The following test cases will be used to validate the TC-2 API:

### 1. Verify API health

```bash
curl --location --request GET 'http://localhost:8080/api/v1/health'
```

### 5. Get all products

```bash
curl --location 'http://localhost:8080/api/v1/products
```

### 6. Create a new product

> TC-1 2.b.iii

```bash
curl --location 'http://localhost:8080/api/v1/products' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Product X",
    "description": "Product X description",
    "price": 13,
    "category_id": 1,
    "active": true
}'
```

### 7. Get a product by id

```bash
curl --location 'http://localhost:8080/api/v1/products/6'
```

### 8. Update a product

> TC-1 2.b.iii

```bash
curl --location --request PUT 'http://localhost:8080/api/v1/products/6' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Product X UPDATED",
    "description": "Product X description UPDATED",
    "price": 12.11,
    "category_id": 1,
    "active": true
}'
```

### 9. Delete a product

> TC-1 2.b.iii

```bash
curl --location --request DELETE 'http://localhost:8080/api/v1/products/6'
```

### 10. Get all products by category

> TC-1 2.b.iv

```bash
curl --location 'http://localhost:8080/api/v1/products/?category_id=1'
```

### 11. Create a new order

```bash
curl --location 'http://localhost:8080/api/v1/orders' \
--header 'accept: application/json' \
--header 'Content-Type: application/json' \
--data '{
  "customer_id": 6
}'
```

### 12. Get an order by id

```bash
curl --location 'http://localhost:8080/api/v1/orders/15' \
--header 'accept: application/json'
```

> The order status should be `OPEN`.

### 13. Add a product to an order

```bash
curl --location 'http://localhost:8080/api/v1/orders/products/15/2' \
--header 'Content-Type: application/json' \
--data '{
    "quantity": 4
}'
```

### 17. Get order histories

```bash
curl --request GET \
--url 'http://localhost/api/v1/orders/histories/?order_id=15'
``` 
