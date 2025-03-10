## Setup

- Follow https://cloud.google.com/bigquery/docs/reference/libraries to setup `GOOGLE_APPLICATION_CREDENTIALS` environment variable
- install needed packages e.g go get -u github.com/jinzhu/gorm
- Increase ulimit value on command prompt e.g. ulimit -n5000 (This will fix `too many open files` issue, since script is launching many goroutines for fetching each transaction status
- Optionally run the script Recursively by updating main() method with ` time.Sleep(time.Second * 5) main() `
