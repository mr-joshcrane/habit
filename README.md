## Example Import

```golang
import "github.com/mr-joshcrane/habit_tracker"
```

## Example of Library Usage
### Local Tracker
```golang
// Error handling omitted for brevity
func main() {
	// File based persistance for our habits
    store, err := pbfilestore.Open("store")
	// A habit tracker backed by our store
    tracker := habit.NewTracker(store)
	habit.RunCLI(tracker)
```


### Client/Server Tracker
```golang
// Error handling omitted for brevity
func main() {
	uri := "yourDynamodbUriConnectionString"
	tablename := "your dynamodb table name"
	// A dynamodb backed persistance store for our habits
	store, err := dynamodbstore.Open(uri, tablename)
	// A habit tracker backed by our store
	tracker := habit.NewTracker(store)
    // Creation of a GRPC Habit server that performs habit logic
	server, err := habit.NewServer(tracker)
    // Our CLI runs using our generated client
	habit.RunCLI(server.Client())
}
```


## Example of CLI usage
#### Perform a new habit
```bash
go run cmd/habit.go violin
```

#### Create a new challenge code from a habit
```bash
go run cmd/habit.go -c violin
```
#### Join a new challenge with a code
```bash
go run cmd/habit.go -v="CCODE" violin
```
## What is it and why would I use it?

Habits are hard to create, so why not leverage the power of social contracts to create them? Available in a single user mode (local users only) or a multi user over the network (some assembly required!)