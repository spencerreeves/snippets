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

- To listen to the channel, we establish a connection to the database using [jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx#hdr-Listen_and_Notify)
then wait for a notification on that connection. This can also be achieved with [lib/pq Listen](https://pkg.go.dev/github.com/lib/pq/example/listen)
function. In this case `s.ChannelName` is `job_channel`as specified in the earlier postgres statement. 
  ```go
  if err = conn.Listen(s.ChannelName); err != nil {
      return errors.Wrap(err, "unable to listen to channel")
  }
  ...
  notification, err := conn.WaitForNotification(ctx)
  ```

## Performance
| Tables   |      Are      |  Cool |
|----------|:-------------:|------:|
| col 1 is |  left-aligned | $1600 |
| col 2 is |    centered   |   $12 |
| col 3 is | right-aligned |    $1 |


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
