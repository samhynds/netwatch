Notes: 
- if the crawl queue is full, and the link extractor tries to add more the crawl worker will be blocked until it can. could be an issue - if so, make link extraction its own goroutine so that the crawler can continue working through the queue and then the links are added when they can. alternatively a large buffer on the queue channel will also solve it

To Do:
    - Roam ✅
    - Extract Links ✅
    - Deduplicate Links ✅
    - Rate limiting (from config & 429s)
    - Extract Content
    - TransportQueue
    - TransportManager
    - TransportWorker
    - Save to DB, External Queue (AMQP), WebSocket 

Later:
    - Auto Content Extraction
    - Auto Categorisation
    - CrawlQueue inQueue and processed will grow forever, need a way of trimming them down over time