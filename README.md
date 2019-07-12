# pidriver
Driver for controlling a tile board with a Raspberry Pi

Example usage:

```go
pi := pidriver.Connect()

pi.Command(pwm, phase, board, 0, 0, sel, 1)

data := pi.Command(0, 0, board, anadr, quad, 0, 0)
fmt.Println(data)
```

Benchmarked on a Pi 3 B+ at ~0.06ms per Command call

TODO: 
```go
pi.Disconnect()
```
