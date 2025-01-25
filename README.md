# NetWatch

> [!NOTE]
> NetWatch is under early active development

NetWatch is a web crawler which crawls sites provided by a config file. It can discover links that match a provided pattern, or can discover all links in "roam" mode. These links (and links of those links and so on) can then be crawled in a recursive manner.

You can extract content from pages in a structured JSON object or simply save the entire response. NetWatch supports HTML pages, JSON responses and image files which are matched against a post processor using their mime-type.

After post-processing, items can be saved to a Postgres database and sent to a Kafka queue for ingestion by other applications.

## Quick Start

## Configuration

## Content Extraction

## Advanced Usage

## Crawl Mechanism

## Queue System

NetWatch contains three queues, the Crawl Queue, Cooldown Queue, and Recrawl Queue. When the application starts, the initial sites are loaded from the provided NetWatch config file and inserted into the Crawl Queue. If the `recrawl` option is enabled, the sites are also inserted into the Recrawl Queue.

When an item is picked up from a queue by the crawl worker, it checks the two rate limits.

1. If the global rate limit (`requests.maxTotal` in the config) has been reached, then the crawl worker will sleep until it's allowed to continue crawling.
2. If the per-host rate limit (`requests.maxPerHost` in the config) has been reached, the item will be placed into the Cooldown Queue with a time at which it can be crawled again without breaching the limit.

Items in the Recrawl Queue ignore any rate limits and are picked up by the Crawl Worker as a priority.

If enabled (with `sites[].links.crawl`), links that match the pattern will be added to the Crawl Queue.

### How are items picked up from the queues?

When the application is first started, the sites defined in the config are loaded into the Crawl Queue and then processed.

When the worker picks up further items, the following priority is used:
1. Recrawl Queue
2. Cooldown Queue
3. Crawl Queue

This ensures recrawled items are requested regularly as defined in the config, and that previously found items that were rate limited are requested before newer items in the Crawl Queue.

## Post Processors
### HTML
### JSON
### Images

## Transporters
### Database
### Queue

---

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
    - Save to DB ✅
    - Check if site already exists in db before crawling it
    - RecrawlQueue
    - External Queue (transport to kafka)
    - per-host rate limiting (not just global)
    - per-host rate limiting doesn't block crawl queue
    - Load balance requests between hosts
        - e.g. I made 5 requests to example.com and have been rate limited.  

Later:
    - HTTP 429
    - Auto Content Extraction
    - Auto Categorisation
    - CrawlQueue inQueue and processed will grow forever, need a way of trimming them down over time
    - check mimetype and parse according to type (e.g. html,json,img etc)


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


---

Cooldown Queue

If a host hits the rate limit, it's moved into the cooldown queue. 