## Block balance
Shows address, which balance has changed (in any direction) more than the others over the last 100(by default) blocks

### How to run
Clone repository. Set your API key in docker compose file(optionally you can set the amount of blocks to scan). And run docker compose.

Or.

You can run it locally, but you need to have go installed. Set your API key in .env file and run `go run main.go`

# Configuration Options

The following table describes the configuration options for the application, their types, how they are sourced from environment variables, and their default values.

| Option      | Type          | Environment Variable | Description                                  | Default Value |
| ----------- | ------------- | -------------------- | -------------------------------------------- | ------------- |
| BlockAmount | int           | BLOCK_AMOUNT         | The amount of blocks to process.             | 100           |
| APIKey      | string        | API_KEY              | The API key for accessing GetBlock services. | None          |
| RPS         | int           | RPS                  | The rate of requests per second.             | 10            |
| NumWorkers  | int           | NUM_WORKERS          | The number of workers to process tasks.      | 10            |
| Timeout     | time.Duration | TIMEOUT              | The timeout duration for tasks.              | None          |

## Example Usage

To configure the application, set the following environment variables:

```bash
export BLOCK_AMOUNT=100
export API_KEY=your_api_key_here
export RPS=10
export NUM_WORKERS=5
export TIMEOUT=30s
```