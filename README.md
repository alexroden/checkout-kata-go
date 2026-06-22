# Checkout Kata Go

## Task Overview

The aim of this task was to build a simple checkout system that can calcuate the total price of scanned items.

### Expected Behaviour

- The checkout system accepts products identified by a SKU string (e.g. "A", "B", "C", "D").
- Each SKU has a fixed unit price that is used for pricing individual items.
- Items can be scanned in any order, and the final price must be independent of scan order.
- The system accumulates scanned items internally until GetTotalPrice() is called.
- The total price is calculated based on:
    - Applying any applicable special pricing rules first (where available).
    - Charging remaining items at their normal unit price.
- Special pricing rules are applied optimally (e.g. bulk discounts):
    - SKU A: 3 for 130 instead of 3 × 50 = 150
    - SKU B: 2 for 45 instead of 2 × 30 = 60
- If the quantity does not fully match a special offer, remaining items are charged at unit price.
    - e.g. 4 × A = 130 + 50
- SKUs without special pricing (C, D) are always charged at unit price only.
- Calling ScanItem(sku) adds a single unit of that SKU to the current basket.
- Calling GetTotalPrice() returns the current total without mutating the scanned items.
- The system should be able to handle multiple scans of the same SKU.

## Project Structure

The API entry point is located in `cmd/main.go`. This file configures the following endpoints:

- * `POST /api/v1/start-session` – Starts a new checkout session.
- * `POST /api/v1/scan-item/{sku}` – Scans an item and adds it to the current session.
- * `GET /api/v1/total` – Retrieves the current checkout total.

The controllers for these endpoints are located in the `pkg/controllers` package. Each controller depends on the `Checkout` service, which provides the implementation of the interface defined for this task and contains the core checkout business logic.

For data persistence, the application uses DynamoDB. The DynamoDB implementation can be found in the `pkg/dynamodb` package.

### Response Models

All service-side errors are returned to the frontend using a standardised error response format:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "basket not found"
  }
}
```

All successful response models returned to the frontend are defined in the `pkg/models` package.

## How to Run

This project is containerised using Docker, so you will need Docker installed on your machine.

To get started, run:

```bash
make setup
```

This command will set up the database, apply migrations, and seed initial user data required to interact with the API.

Once setup is complete, start the application with:

```bash
make start
```

The API will then be available locally at:

```
http://127.0.0.1:8080
```

You can use the provided Postman collection to interact with the endpoints.

When you are finished, clean up the environment by running:

```bash
make cleanup
```

This will stop and remove all running containers created by the project.

## Endpoints

### Start Session

To interact with the core functionality of this API, you must first create a session ID. This is done using the following endpoint:

`POST /api/v1/start-session`

The response will include a session ID, which should be provided in the `Session-Id` header for all subsequent requests to the other endpoints.

**Response:**

```json
{
  "data": {
    "id": "1B5B2D92A6CF4B87A20014E1205BDBCC"
  }
}
```

### Scan Item

Once a session has been created, items can be added to the basket using this endpoint:

`POST /api/v1/scan-item/{sku}`

If the item is successfully added to the basket, the endpoint will return a `204 No Content` response.

### Get Total

After all items have been scanned, the checkout total can be retrieved using the following endpoint:

`GET /api/v1/total`

**Response:**

```json id="9xq2a1"
{
  "data": {
    "total": 175
  }
}
```

This endpoint calculates the total price of all items currently in the basket, including any applicable promotions defined in the `SpecialsTable`.

## Challenges and Future Considerations

- One of the main challenges in this task was deciding on the initial implementation approach. My first consideration was to build a command-line tool that used in-memory storage to persist data between operations.
However, I opted for a more extensible approach using a service-based API architecture. This design better future-proofs the implementation by making it easier to interact with, extend, and integrate into other systems if required.
- For this task, the use of DynamoDB means that products and specials could be managed through additional CRUD endpoints for each entity.
- Because DynamoDB was used, additional logic was required to handle the creation and seeding of tables. This logic can be found in the `internal/dynamo-init` folder.
- Additional considerations for a production-ready version of this system would include handling concurrency when updating basket state in DynamoDB, particularly under high request volumes. Session lifecycle management, such as expiry and cleanup of abandoned sessions, would also be important to prevent unnecessary data growth. As the number of promotional rules increases, a more structured pricing or rules engine may be required to manage complexity and ensure consistent evaluation order.





