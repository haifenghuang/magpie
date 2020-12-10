include "examples/sub_package/calc"

printf("10+2 = %d\n", calc.Add(10, 2))
printf("10-2 = %d\n", calc.Minus(10, 2))

//multiply is not exported, so below line will report runtime error.
//printf("10*2=%d\n", calc.multiply(10, 2))

printf("2*2 = %d\n", calc.SQUARE_OF_TWO)
//printf("4*4 = %d\n", calc.sQUARE_OF_FOUR)
