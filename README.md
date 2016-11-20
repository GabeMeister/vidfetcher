# YouTube Video Fetcher

Golang program to fetch YouTube video data and store it in a database. This is meant to run as a background process to support the [YouTube Collections Revamp Server.](https://github.com/GabeMeister/Youtube-Collections-Revamp-Server)

## Progression of Data Fetching Strategies:

Strategy 1: Just synchronously calling to youtube api once per youtube id. Yawnnn

Strategy 2: Created a go routine for every youtube channel in databas (6,000+ channels). Yeah, too many open sockets

Strategy 3: "Batch" together requests of about 25-50 channels to make requests for, and wait until they are all completed. Still one channel per request

Strategy 4: Same as strategy 3, but include 50 youtube ids per api call. 
