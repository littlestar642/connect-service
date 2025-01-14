## Design Considerations

- Following `Accept Interface and Return Struct` pattern.
- Using SOLID principles for separation of concerns and making the code modular.

- We are using gin for building the RESTful APIs, redis as the cache and kafka as the message queue(streaming service).
- Using a retry mechanism to ensure connection with kafka at the application layer, since observed that at times docker compose takes time to start the service.
- Handling graceful shutdown to avoid memory leaks and finish ongoing tickers. 

## Data flow Overview

- GET /api/verve/accept
    - Added validation on ID.
    - Accepting url escaped endpoint as query param.
    - Added Redis as a distributed cache for ensuring uniqueness of ID in every minute window.
    - Using redis to store the count of requests in every minute window.
    - Using a Ticker to regularly fetch the count and push on to kafka.

## Assumptions and TradeOffs

- For the sake of simplicity, we have not added any authentication or authorization layer. This can be added in the future.
- Not using a DB, since from the problem statement writing to log file/pushing onto queue seems to be the objective of the service.