# Job Queue

This service is used to receive instructions, or jobs, which are then added to a queue and executed.

Jobs are received through an gRPC service, which are then added to a database and executed by RPC to the desired destination.

## Variables

-   `WORKER_LIMIT` - is an environment variable used to limit how many jobs can be executed concurrently. This is not a required variable and the default value is 3. The value of this variable must be a valid int32.
