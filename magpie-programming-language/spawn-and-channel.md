# Spawn and channel

 You can use `spawn` to create a new thread, and `chan` to communicate with the thread.

```swift
let aChan = chan()
spawn fn() {
    let message = aChan.recv()
    println('channel received message=<{message}>')
}()

//send message to thread
aChan.send("Hello Channel!")
```

 You could use channel and spawn togeter to support lazy evaluation:

```swift
// XRange is an iterator over all the numbers from 0 to the limit.
fn XRange(limit) {
    ch = chan()
    spawn fn() {
        //for (i = 0; i <= limit; i++)  // Warning: Never use this kind of for loop, or else you will get weird results.
        for i in 0..limit {
            ch.send(i)
        }

        // Ensure that at the end of the loop we close the channel!
        ch.close()
    }()
    return ch
}

for i in XRange(10) {
    fmt.println(i)
}
```

