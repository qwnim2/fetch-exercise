# Receipt Processor - Fetch exercise

This application is a receipt processing service. Submit a receipt using POST request, and it returns an ID for the receipt. Retrieve the points awarded to that receipt using the ID.

## Endpoints

1. **POST `/receipts/process`**  
   Submits a receipt for processing.
   
   **Request Body Example (JSON):**
   ```json
   {
     "retailer": "Target",
     "purchaseDate": "2022-01-01",
     "purchaseTime": "13:01",
     "total": "6.49",
     "items": [
       {
         "shortDescription": "Mountain Dew 12PK",
         "price": "6.49"
       }
     ]
   }
   ```

   **Response Example (JSON)**
   ```json
   {
    "id": "c1bc44bc-3d6e-423c-a27c-6bda660a921f"
   }
   ```

2. **GET `/receipts/{id}/points`**
   
   **Path Example: **
   ```
   The unique ID returned from the POST /processReceipt call
   localhost:8801/receipts/c1bc44bc-3d6e-423c-a27c-6bda660a921f/points
   ```

   **Response Example (JSON)**
   ```json
   {
    "points":37
   }
   ```
