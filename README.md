Notes: 
- if the crawl queue is full, and the link extractor tries to add more the crawl worker will be blocked until it can. could be an issue - if so, make link extraction its own goroutine so that the crawler can continue working through the queue and then the links are added when they can. alternatively a large buffer on the queue channel will also solve it

To Do:
    - Roam ✅
    - Extract Links ✅
    - Deduplicate Links ✅
    - Rate limiting from config ✅
    - Extract Content ✅
    - TransportQueue ✅
    - TransportManager ✅
    - TransportWorker ✅
    - RecrawlQueue
    - Save to DB, External Queue (AMQP), WebSocket 

Later:
    - HTTP 429
    - Auto Content Extraction
    - Auto Categorisation
    - CrawlQueue inQueue and processed will grow forever, need a way of trimming them down over time


Transport
    - AMQP
    - DB (postgres)
    - WebSocket
    - Data:
        - URL
        - timestamp
        - content (could be any key, value is string)
        - links
        - full body
        - headers

DB connect: export NW_DB_CONN_URL="postgres://netwatch:test@localhost:5432/netwatch"