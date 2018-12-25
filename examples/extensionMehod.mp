fn float$to_integer() {
   return ( int( self ) );
}

printf("12.5.to_integer()=%d\n", 12.5.to_integer())

fn array$find(item) {
   i = 0;
   length = len(self);

   while (i < length) {
     if (self[i] == item) {
       return i;
     }
     i++;
   }

   // if not found
   return -1;
};

idx = [25,20,38].find(10);
printf("[25,20,38].find(10) = %d\n", idx) // not found, return -1

idx = [25,20,38].find(38);
printf("[25,20,38].find(38) = %d\n", idx) //found, returns 2
