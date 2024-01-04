# Migrations
### Make sure migrate CLI is installed in the system.
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
### Create new table migration. Change create_users_table for your table name - create_name_table 
```
migrate create -ext sql -dir internal/database/migrations -seq create_users_table
```
This will generate and up and down file with the given name


### Run generated migration 
```
migrate -database "database_url" -path internal/database/migrations/  up
```

### Reverse the generated migration 
```
migrate -database "database_url" -path internal/database/migrations/  down 
```

### Force a specific migration
```
migrate -database "database_url" -path internal/database/migrations/  force 20231002034516
```

# Queues
### Create a Function with the queue functionality defined. This is for each item passed to the queue
```
processFunc := func(item queue.QueueItem) error {
    // Simulate processing
    fmt.Printf("Processing item %d: %s\n", item.ID, item.Data)
    time.Sleep(time.Second)

    // Simulate failure for demonstration purposes
    if item.ID%3 == 0 {
        return fmt.Errorf("failed processing item %d", item.ID)
    }

    return nil
}
```
### Create a New Queue and pass the defined function
```
queueService := queue.NewQueueService(3, 2, processFunc)
```

### Start the queue
```
go queueService.StartProcessor(ctx)

for index, post := range dbPosts {
    item := queue.QueueItem{ID: index, Data: post}
    queueService.Enqueue(item)
}
```

