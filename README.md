# Batch Logger

A simple web-server made using `golang` and `chi`. The server writes logs in batchs in intervals to a post endpoint. You need to configure `Batch Size`, `Interval Time` and `Post endpoint` via environemnt variables.

## About the codebase

There are two functioning versions of said application in the commit history. The [older](https://github.com/utsav0209/batch-logger/tree/cf6e9b0080c4e29faa4710f6bd0e1c4cd54b9ef7) version was using shared memory for triggering the post endpoint syncing which could end up in race condition when throughput of input is very high and is not concurrently scalable. The newer version which is current [head](https://github.com/utsav0209/batch-logger) on master branch uses go channels for supplying logs and is taking advantage of concurrency. The insipiration for which came from [this talk](https://www.youtube.com/watch?v=oV9rvDllKEg) by Rob pyke.

## Steps to run

- run: `docker build . -t batch-logger`
- create an environment file as described in sample.env or pass env variables from command line
- run: `docker run -it --rm -p 3000:3000 --env-file .env batch-logger`

## Endpoints

- `GET: /healthz` - Returns health of the server
- `POST: /log` - Send your logs here
