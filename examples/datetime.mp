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

let dt7 = dt/2018-01-01 12:01:00/.addDate(1, 2, 3).add(time.SECOND * 10) //add 1 year, two months, three days and 10 seconds
println(dt7)

let dt8 = dt/2018-01-01 12:01:00/.unix() //convert datetiem to timestamp
println(dt8)


let dt9 = dt/2018-01-01 12:01:00/.unix() //convert datetiem to timestamp
let timeStr = unixTime(dt9).format() //convert timestamp to datetime string
println(timeStr)
