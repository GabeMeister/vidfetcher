# YouTube Video Fetcher

Golang program to fetch YouTube video data to store in a Postgres database. This is meant to run as a background process to support the [YouTube Collections Revamp Server.](https://github.com/GabeMeister/Youtube-Collections-Revamp-Server)

## Data Fetching Optimization Techniques:

- Just synchronously calling to youtube api once per youtube id is extremely slow

- Created a go routine for every youtube channel in database (6,000+ channels). Uses many open sockets. 6,000+ "threads" is a little too much

- You can "batch" together requests of about 25-50 channels to make requests for, and wait until they are all completed. Still one channel per request

- Turns out you can make api requests with up to 50 channels. So instead of 1 api request containing 1 youtube id, you can hit 50 birds with 1 stone.

- For channels that don't have any uploads, we can ignore them.

## Potential Ideas to Explore:

- Instead of "waiting" to form one big slice of all channel data, just begin fetching videos of channels that are out of date

- Instead of reading every video from the api no matter what, check the video count in the database first, and only if different, begin fetching all videos from api

## Open Questions:

- How to properly "print" data that you want to see, but not have to do extra work. For example, tasks.AreVideosOutOfDate() returns true if a channels videos are out of date, and false otherwise. This function must make a query to a database to check the video count for a channel. I want to be able to have the command line program print the video count in the database, and the video count retrieved through the api. Ideally, I only want to query the database once, because that's technically only what we need. 

- Whether to pass just ids to functions, or pass objects that contain ids to functions.