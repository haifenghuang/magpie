//datetime literal
let month = "01"
let dt1 = dt/2018-{month}-01 12:01:00/
let dt2 = dt/2019-01-01 12:01:00/

println(dt1)
println(dt2)
//two datetime literals could be compared using '>', '>=', '<', '<=' and '=='
println(dt1 <= dt2)


/* 'datetime literal' + string:
     string support 'YMDhms' where
       Y:Year    M:Month    D:Day
       h:hour    m:minute   s:second

*/
let dt3 = dt/2018-01-01 12:01:00/ + "-12M" // -12M: minus 12 months
println(dt3) # result: 2017-01-01 12:01:00

let dt4 = dt/2018-01-01 12:01:00/ + "10h0m12s" // 10h0m12s: add 10 hours and 12 seconds
println(dt4) # result: 2018-01-01 22:01:12

let dt5 = dt/2018-01-01 12:01:00/.add(time.SECOND * 10) //add 10 seconds
println(dt5)

let dt6 = dt/2018-01-01 12:01:00/.addDate(1, 2, 3) //add 1 year, two months and three days.
println(dt6)

let dt71 = dt/2018-01-01 12:01:00/.addDate(1, 2, 3).add(time.SECOND * 10) //add 1 year, two months, three days and 10 seconds
printf("dt71 = %v\n", dt71)

//same result as 'dt71'
let dt72 = dt/2018-01-01 12:01:00/ + "1Y2M3D10s" //add 1 year, two months, three days and 10 seconds
printf("dt72 = %v\n", dt72)
//same resutl as above
//printf("dt72 = %s\n", dt72.toStr()) //use 'toStr()' method to convert datetime to string.

let dt8 = dt/2018-01-01 12:01:00/.unix() //convert datetiem to timestamp
println(dt8)


let dt9 = dt/2018-01-01 12:01:00/.unix() //convert datetiem to timestamp
let timeStr = unixTime(dt9).format() //convert timestamp to datetime string
println(timeStr)

let dt10 = dt/2019-01-01 12:01:00/
//convert to yyyy/mm/dd hh:mm:ss
format = dt10.strftime("%Y/%m/%d %T")
println(dt10.toStr(format))

////////////////////////////////
// time object to timestamp
////////////////////////////////
println(dt10.unix()) //to timestamp(UTC)
println(dt10.unixNano()) //to timestamp(UTC)
println(dt10.unixLocal()) //to timestamp(LOCAL)
println(dt10.unixLocalNano()) //to timestamp(LOCAL)

////////////////////////////////
// timestamp to time object
////////////////////////////////
timestampUTC = dt10.unix()      //to timestamp(UTC)
println(unixTime(timestampUTC)) //timestamp to time

timestampLocal = dt10.unixLocal() //to timestamp(LOCAL)
println(unixTime(timestampLocal)) //timestamp to time