# Decimal

 In magpie, decimal is Arbitrary-precision fixed-point decimal numbers. And the code mainly based on [decimal](https://github.com/shopspring/decimal).

 Please see below examples:

```swift
d1 = decimal.fromString("123.45678901234567")  //create decimal from string
d2 = decimal.fromFloat(3)  //create decimal from float

//set decimal division precision.
//Note: this will affect all other code that follows
decimal.setDivisionPrecision(50)

fmt.println("123.45678901234567/3 = ", d1.div(d2))  //print d1/d2
fmt.println(d1.div(d2)) //same as above

fmt.println(decimal.fromString("123.456").trunc(2)) //truncate decimal

//convert string to decimal
d3=decimal("123.45678901234567")
fmt.println(d3)
fmt.println("123.45678901234567/3 = ", d3.div(d2))
```

