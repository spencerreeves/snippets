# Postgres as a Pub/Sub service

Recently I heard that you can use Postgres as a Pub/Sub service using its built in `pg_notify` function. So I took some
time creating a test that could give some insights into the capabilities using Postgres as a Pub/Sub service.

## Code Recap

I have highlighted the important bits below so that you can create your own Postgres Pub/Sub service.

- First thing we care about is having adding a trigger on the table we are watching. In this case we will trigger
  after every insert, update, or delete.
  ```postgresql
  CREATE TRIGGER job_trigger
    AFTER INSERT OR UPDATE OR DELETE
    ON job
    FOR EACH ROW
  EXECUTE PROCEDURE job_notify();
  ```

- Next we need to notify the channel which can be achieved by executing `pg_notify` with your channel name and payload.
  In this case our channel name is `job_channel` and our payload is a constructed json object.
  ```postgresql
  PERFORM 
  pg_notify('job_channel', 
            json_build_object('job_id', job_id,
                              'notification_id', notification_id,
                              'status', job_status)::TEXT
  );
  ```

- To listen to the channel, we establish a connection to the database
  using [jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx#hdr-Listen_and_Notify)
  then wait for a notification on that connection. This can also be achieved
  with [lib/pq Listen](https://pkg.go.dev/github.com/lib/pq/example/listen)
  function. In this case `s.ChannelName` is `job_channel`as specified in the earlier postgres statement.
  ```go
  if err = conn.Listen(s.ChannelName); err != nil {
      return errors.Wrap(err, "unable to listen to channel")
  }
  ...
  notification, err := conn.WaitForNotification(ctx)
  ```

## Performance

I ran a load test using my local machine and achieved the following performance. As you will notice, the bottleneck
was writing the rows to the database. It is possible the extra step of notifying the channel after a row was inserted
caused the delay, but additional testing is needed to determine that. Overall, we saw a performance of about **6K 
writes/notifies per second**. 

_Side note_: We did see a slight increase (+0.8K) writes/notifies per second when we turned off pretty print logging. There
are likely other enhancements we could do to achieve a higher throughput.

#### Test System

```
Macbook 13in, 2019
OS: 10.15.7 (19H1615)
Processor: 2.8 GHz Quad-Core Intel Core i7
Memory: 16 GB 2133 MHz LPDDR3
Graphics: Intel Iris Plus Graphics 655 1536 MB
```

| Rows Created | Bytes Written |  Time(write/notify)   | Performance (writes/notify) / second |
|:------------:|:-------------:|:---------------------:|:------------------------------------:|
|      1K      |     400KB     |   162.7ms / 503.6ms   |             6146 / 1988              |
|     100K     |    34272KB    | 16604.7ms / 17005.8ms |             6022 / 5880              |
|   10,000K    |    3340MB     |     1637s / 1638s     |             6105 / 6105              |

<details>
<summary>See each runs output below</summary>
<li>For the 1K inserts</li> 
<pre>
10:47PM INF  Busy=161.543383 End=2022-05-15T22:47:48-07:00 Errors=0 Idle=0.698485   Processed=250  Start=2022-05-15T22:47:47-07:00 Total Duration=162.295176
10:47PM INF  Busy=161.934899 End=2022-05-15T22:47:48-07:00 Errors=0 Idle=0.730477   Processed=250  Start=2022-05-15T22:47:47-07:00 Total Duration=162.71577
10:47PM INF  Busy=161.213466 End=2022-05-15T22:47:48-07:00 Errors=0 Idle=0.708023   Processed=250  Start=2022-05-15T22:47:47-07:00 Total Duration=161.988129
10:47PM INF  Busy=161.715501 End=2022-05-15T22:47:48-07:00 Errors=0 Idle=0.642744   Processed=250  Start=2022-05-15T22:47:47-07:00 Total Duration=162.422688
10:47PM INF  Busy=0.231242   End=2022-05-15T22:47:48-07:00 Errors=0 Idle=502.983911 Processed=2000 Start=2022-05-15T22:47:47-07:00 Total Duration=503.660254
</pre>
<li>For the 100K inserts</li> 
<pre>
10:49PM INF  Busy=16524.533874 End=2022-05-15T22:49:02-07:00 Errors=0 Idle=72.28105     Processed=25000  Start=2022-05-15T22:48:46-07:00 Total Duration=16603.092408
10:49PM INF  Busy=16523.869366 End=2022-05-15T22:49:02-07:00 Errors=0 Idle=73.611458    Processed=25000  Start=2022-05-15T22:48:46-07:00 Total Duration=16603.861422
10:49PM INF  Busy=16522.619537 End=2022-05-15T22:49:02-07:00 Errors=0 Idle=73.318314    Processed=25000  Start=2022-05-15T22:48:46-07:00 Total Duration=16602.535577
10:49PM INF  Busy=16524.537404 End=2022-05-15T22:49:02-07:00 Errors=0 Idle=73.884878    Processed=25000  Start=2022-05-15T22:48:46-07:00 Total Duration=16604.702015
10:49PM INF  Busy=24.345814    End=2022-05-15T22:49:03-07:00 Errors=0 Idle=16935.544612 Processed=200000 Start=2022-05-15T22:48:46-07:00 Total Duration=17005.837183
</pre>
<li>For the 10M inserts</li> 
<pre>
11:17PM INF  Busy=1630089.20162  End=2022-05-15T23:17:19-07:00 Errors=0 Idle=7126.715646   Processed=2500000  Start=2022-05-15T22:50:01-07:00 Total Duration=1637816.150764
11:17PM INF  Busy=1629653.814291 End=2022-05-15T23:17:18-07:00 Errors=0 Idle=7115.449761   Processed=2500000  Start=2022-05-15T22:50:01-07:00 Total Duration=1637369.52435
11:17PM INF  Busy=1629528.026562 End=2022-05-15T23:17:18-07:00 Errors=0 Idle=7113.604759   Processed=2500000  Start=2022-05-15T22:50:01-07:00 Total Duration=1637298.563681
11:17PM INF  Busy=1629594.307956 End=2022-05-15T23:17:18-07:00 Errors=0 Idle=7113.918192   Processed=2500000  Start=2022-05-15T22:50:01-07:00 Total Duration=1637307.170204
11:17PM INF  Busy=2339.845426    End=2022-05-15T23:17:19-07:00 Errors=0 Idle=1631231.39163 Processed=20000000 Start=2022-05-15T22:50:01-07:00 Total Duration=1638058.451103
</pre>
</details>

## Local development

#### Environment configs

All configs live in the `.env` file at the base directory.

| Config Name     |               Description                | Example                                                     |
|-----------------|:----------------------------------------:|:------------------------------------------------------------|
| DEBUG           | Enables debug settings like pretty print | true                                                        |
| DATABASE_URL    | Used to establish connection to database | postgres://postgres@127.0.0.1:5432/postgres?sslmode=disable |
| PUB_SUB_CHANNEL |   Specifies the database channel name    | job_channel                                                 |

#### Run locally

- Make sure the `DATABASE_URL` is set, download [dbmate](https://github.com/amacneil/dbmate), and create the necessary
  tables and triggers.
  ```shell
  dbmate up
  ```
- To run notification service in logging mode
  ```shell
  go run main.go
  ```
- To run load test
  ```shell
  go run main.go load
  ```
