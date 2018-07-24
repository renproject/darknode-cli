# Cogo

A package for high-level concurrency patterns, in Go.

## Usage

### Parallelism

Cogo offers two types of parallelism: implicit and explicit. Implicit parallelism means that Cogo will control how many goroutines are used to introduce parallelism, in contrast to explicit parallelism which gives control to the user.

**Implicit Parallelism**

In the following example, we use the `co.ForAll` function to loop over different iterators. Iterators are any value that makes sense to loop over: arrays, slices, maps, and integers. Cogo will use one goroutine per CPU core available and we cannot make any assumptions about which iteration will run on which goroutine. Calling `co.ForAll` will block until all iterations have finished running.

```go
// Fill an array of integers with random values
xs := [10]int{}
co.ForAll(xs, func(i int) {
    xs[i] = rand.Intn(10)
})

// Map those random values to booleans
ys := [10]bool{}
co.ForAll(10, func(i int) {
    ys[i] = xs[i] > 5
})
```

In the following example, we use the `co.Begin` function to run distinct tasks. As before, Cogo will use one goroutine per CPU core available and map the different tasks over these goroutines. Calling `co.Begin` will block until all tasks have finished running.

```go
co.Begin(
    func() {
        log.Info("[task 1] when will this print?")
    },
    func() {
        log.Info("[task 2] who knows?")
    },
    func() {
        log.Info("[task 3] implicit parallelism is great!")
    })
```

**Explicit Parallelism**

In the following example, we use the `co.ParForAll` function to loop over different iterators. As with implicitly parallel loops, iterators are any value that makes sense to loop over: arrays, slices, maps, and integers. Unlike implicitly parallel loops, Cogo will use one goroutine per iteration, regardless of the number of CPU cores available. Calling `co.ParForAll` will block until all iterations have finished running.

```go
// Fill an array of integers with random values
xs := [10]int{}
co.ParForAll(xs, func(i int) {
    xs[i] = rand.Intn(10)
})

// Map those random values to booleans
ys := [10]bool{}
co.ParForAll(10, func(i int) {
    ys[i] = xs[i] > 5
})
```

In the following example, we use the `co.ParBegin` function to run distinct tasks. Similar to the `co.ParForAll` function, Cogo will use one goroutine per task. Calling `co.ParBegin` will block until all tasks have finished running.

```go
co.ParBegin(
    func() {
        log.Info("[task 1] when will this print?")
    },
    func() {
        log.Info("[task 2] who knows?")
    },
    func() {
        log.Info("[task 3] explicit parallelism is great!")
    })
```
