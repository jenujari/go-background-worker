# Go Background job worker pattern
Job scheduler pattern to make background worker that can handle following triggers

- System Interrupt and System Termination signal
- Context closing event.
- Error case in any one of the background worker.

All workers will run in the async way in a parallel goroutine.