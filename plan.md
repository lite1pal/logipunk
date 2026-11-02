# Logipunk

A tool to analyze log files and report useful analytics.

For example:

1. Which endpoints are the slowest
2. How many error logs
3. Repeated IP addresses with failed requests

## Example usage

```bash
2026-10-30T08:12:01Z INFO  192.168.1.10 GET /api/products 200 32ms
2026-10-30T08:12:03Z ERROR 192.168.1.12 POST /api/checkout 500 120ms
2026-10-30T08:12:04Z INFO  192.168.1.10 GET /api/products 200 28ms
2026-10-30T08:12:05Z WARN  192.168.1.15 GET /api/login 401 15ms
2026-10-30T08:12:08Z ERROR 192.168.1.12 POST /api/checkout 500 145ms
2026-10-30T08:12:11Z INFO  192.168.1.11 GET /api/users 200 45ms
2026-10-30T08:12:15Z ERROR 192.168.1.20 GET /api/orders 503 300ms
```

You run:

```
  logipunk app.log
```

Output:

```
LOG ANALYSIS
============

File: app.log
Total entries: 7

LEVELS
------
INFO       3
WARN       1
ERROR      3

HTTP STATUS
-----------
200        3
401        1
500        2
503        1

TOP ENDPOINTS
-------------
/api/products       2
/api/checkout       2
/api/login          1
/api/users          1
/api/orders         1

SLOWEST REQUESTS
----------------
GET /api/orders     300ms
POST /api/checkout  145ms
POST /api/checkout  120ms

UNIQUE IP ADDRESSES: 4
```
