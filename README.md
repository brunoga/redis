# redis

This package implements a simple redis-based RW lock.

- Only supports single Redis instance.
- Designed for read-heavy use cases with a simple write starvation-prevention strategy.

This is implemented by making use of 2 keys. One to hold a read counter and one to hold a writer counter (for consistency reasons. We can only have 1 or 0 writers).

Write lock:

If there is an active writer, we do not acquire the lock.
If there is no active writer, we check if there are any active readers (by checking if the writer counter exists). If there are readers, we record our intent to write (writer count set to 0) but also do not acquire the lock.
If there are no active readers, we acquire the lock by setting the writer count to 1.

Write Unlock:

Delete the writer count key for the lock.

Read lock:

We check if there is an active writer or an intent to write (this is to prevent write starvation). If so, we do not acquire the lock.
If there is no active writer or intent to write, we acquire the read lock by incrementing the reader count.

Read unlock:

Decrement the reader count. If it reaches 0, we delete the reader count key.


