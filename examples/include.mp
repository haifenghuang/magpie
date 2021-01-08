include "examples/sub_package/calc"

a, b = 10, 2
printf("10 + 2 = %d\n", calc.Add(a, b))
printf("10 - 2 = %d\n", calc.Minus(a, b))

//multiply is not exported, so below line will report runtime error.
//printf("10*2=%d\n", calc.multiply(a, b))

printf("2 * 2 = %d\n", calc.SQUARE_OF_TWO)
//printf("4*4 = %d\n", calc.sQUARE_OF_FOUR)
