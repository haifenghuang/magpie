##############################################
# Tesing for LINQ's Query expression style
##############################################

//Create the data source
let scores = [ 97, 92, 79, 60, 88, 95, 84];

// Create the query.
let queryHighScores = from score in scores where score > 80 && score % 2 == 0 select score;

// Execute the query.
println()
for i in queryHighScores {
    printf("%d ", i)
}

printf("\n================================================================\n\n")
let ingredients = [
    {Name: "Sugar",  Calories: 500},
    {Name: "Egg",    Calories: 100},
    {Name: "Milk",   Calories: 150},
    {Name: "Flour",  Calories: 50},
    {Name: "Butter", Calories: 200},
]

ingredient = from i in ingredients where i.Calories >= 150 orderby i.Name select i
for item in ingredient {
    println(item)

    // uncomment for formatted output
    //itemStr = json.marshal(item)
    //println(json.indent(itemStr, "  "))
}

printf("\n================================================================\n\n")

let students = [
       { First: "Svetlana", Last:"Omelchenko", ID: 111, Scores: 77.5 },
       { First: "Claire",   Last:"O'Donnell",  ID: 112, Scores: 72.25 },
       { First: "Sven",     Last:"Mortensen",  ID: 113, Scores: 93.5 },
       { First: "Cesar",    Last:"Garcia",     ID: 114, Scores: 75.5 },
       { First: "Debra",    Last:"Garcia",     ID: 115, Scores: 88.25 } 
]
booleanGroupQuery = from student in students group student by student.Scores >= 80; //pass or fail!

// Execute the query and access items in each group
for studentGroup in booleanGroupQuery { //convert it to an array
    println(studentGroup.key == true ? "High averages" : "Low averages");

    for student in studentGroup.value {
        printf("   %s, %s: %v\n", student.Last, student.First, student.Scores)
    }
}


printf("\n================================================================\n\n")
//Create the data source
let students = [
    { First: "Svetlana", Last: "Omelchenko", ID: 111 },
    { First: "Claire",   Last: "O'Donnell",  ID: 112 },
    { First: "Sven",     Last: "Mortensen",  ID: 113 },
    { First: "Cesar",    Last: "Garcia",     ID: 114 },
    { First: "Debra",    Last: "Garcia",     ID: 115 },
]

// Create the query.
let sortedStudents =
            from student in students
            orderby student.Last ascending, student.First ascending
            select student

// Execute the query.
println("sortedStudents:")
for student in sortedStudents {
    printf("\t%s %s\n", student.Last, student.First)
}

printf("\n================================================================\n\n")
file = newFile("./examples/linqSample.csv", "r")
result = from line in file where int(line[1]) > 300000 select line[0]
for item in result
{
    printf("item = %s\n", item)
}


//Magpie can parse below linq expression
//without any problem, but the evaluation wil fail.

//stringList = [
//            "A penny saved is a penny earned.",
//            "The early bird catches the worm.",
//            "The pen is mightier than the sword." 
//];
//earlyBirdQuery =
//    from sentence in stringList
//    let words = sentence.split(" ")
//    from word in words
//    let w = word.lower()
//    where w[0] == "a" || w[0] == "e" ||
//          w[0] == "i" || w[0] == "o" ||
//          w[0] == "u"
//    select word