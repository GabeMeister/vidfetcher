# YouTube Video Fetcher

Golang program to fetch YouTube video data to store in a Postgres database. This is meant to run as a background process to support the [YouTube Collections Revamp Server.](https://github.com/GabeMeister/Youtube-Collections-Revamp-Server)

## Data Fetching Optimization Techniques:

- Just synchronously calling to youtube api once per youtube id is extremely slow

- Created a go routine for every youtube channel in database (6,000+ channels). Uses many open sockets. 6,000+ "threads" is a little too much

- You can "batch" together requests of about 25-50 channels to make requests for, and wait until they are all completed. Still one channel per request

- Turns out you can make api requests with up to 50 channels. So instead of 1 api request containing 1 youtube id, you can hit 50 birds with 1 stone.

- For channels that don't have any uploads, we can ignore them.

- Originally thought I would have "batch" channel fetches. So for 50 channels at a time, I would run the fetching, then when all 50 are up, I would start a new "batch" of 50 channels to fetch videos for. But [thanks to this StackOverflow answer](http://stackoverflow.com/a/25324090/1751481) I can do a "rolling" go routine fetch technique, where as soon as one go routine ends, another is started.

- Kind of a no-brainer, but instead of checking the count(*) of videos that have a particular channel id, instead just check the VideoCount column in the Channels table

- Before, I was manually copying over all data that I cared about from the api calls. Now I'm just storing a reference to the youtube.Channel in my YoutubeChannel struct, and the attributes that I care about are just function getters.

- Using "bounded parallelism" seems to be a lot simpler and easier to do than the "fan in" strategy that I initially used. It's a simpler concept: start a bunch of go routines that all read off of the same channel, and then push all your stuff "down" the channel. The fan in strategy is a lot more involved, with creating a go routine for each item or "batch" of asynchronous tasks you want to perform. 

- Things seemed to "slow down" after about 10 min running the fetcher. I bumped the number of concurrent go routines down to 10, and created separate database instances for each go routine. Things are running a lot smoother now.

- To fetch the channel ids for all the channel youtube ids, I just do 1 big sql query with an ANY keyword in postgres. This query takes like 5 minutes to run, but I still believe it's faster than doing one at a time or "batches".

- I used to be printing out the video titles of each video I was inserting, but it was slowing down the ubuntu terminal after an extended period of time. I took out these print statements and the program can run for a lot longer time now. (almost a full day without crashes or halts)

## Potential Ideas to Explore:

- Instead of "waiting" to form one big slice of all channel data, just immediately begin fetching videos of channels that are out of date

- Instead of reading every video from the api no matter what, check the video count in the database first, and only if different, begin fetching all videos from api

- Currently, I'm slightly "cheating" my way of having to check video-for-video whenever I fetch new video uploads for a channel. To determine if I need to fetch videos for a channel, I check if the video count for the channel in the database matches the video count found from the api. If it's different, I fetch new videos only up until I find one that I've already seen from the database. This strategy works great for more efficiently grabbing new videos on a channel, but doesn't account for videos that may be deleted off the channel. So currently the vid fetcher script is printing "Inserting 0 new videos for channel xyz", when in reality there's no new videos to fetch, just one that (most likely) got deleted and made the counts not match. This should be solved with more careful checking of the video counts. If the api count is greater than the database count, then most likely there's new videos. If the api count is LESS than the database count, then most likely videos got deleted and we have to do a "hard" search through the api.

## Thoughts Log:

- How to properly "print" data that you want to see, but not have to do extra work. For example, tasks.AreVideosOutOfDate() returns true if a channels videos are out of date, and false otherwise. This function must make a query to a database to check the video count for a channel. I want to be able to have the command line program print the video count in the database, and the video count retrieved through the api. Ideally, I only want to query the database once, because that's technically only what we need. 

- Whether to pass just ids to functions, or pass objects that contain ids to functions. Overall just when to put things into objects that represent it, or just use the raw data.

- How much should the data access layer check before doing actions? Should a SelectVideoCountOfChannel() function check if the passed in Channel ID exists in the database, and if not, throw an exception? Should it just query anyway and just return 0 rows?

- For every task, there seems to be a balance that you have to strike between what data you know "beforehand" vs. what work you do "while you go". For instance, I could calculate the total number of channels that are out of date, or I could instead just iterate through all the channels, and check if each (one at a time) is out of date, and then immediately start fetching videos, and just keep track of the count as you go. One way takes longer initially to start fetching videos, but gives the amount of work that the video fetcher is about to do. The other more quickly starts fetching videos, but you don't get to see how much work it has to do until it has completed it.



