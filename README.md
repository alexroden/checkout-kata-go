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

## Project structure


