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
       { First: "Svetlana", Last:"Omelchenko", ID: 111, Scores: [97, 72, 81, 60] },
       { First: "Claire",   Last:"O'Donnell",  ID: 112, Scores: [75, 84, 91, 39] },
       { First: "Sven",     Last:"Mortensen",  ID: 113, Scores: [99, 89, 91, 95] },
       { First: "Cesar",    Last:"Garcia",     ID: 114, Scores: [72, 81, 65, 84] },
       { First: "Debra",    Last:"Garcia",     ID: 115, Scores: [97, 89, 85, 82] } 
]
booleanGroupQuery = from student in students group student by student.Scores.average() >= 80; //pass or fail!
// Execute the query and access items in each group
for studentGroup in booleanGroupQuery { //convert it to an array
    println(studentGroup.key == true ? "High averages" : "Low averages");

    for student in studentGroup.value {
        printf("   %s, %s: %v\n", student.Last, student.First, student.Scores.average())
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
file.close() //do not forget to clsoe the file

printf("\n================================================================\n\n")
stringList = [
            "A penny saved is a penny earned.",
            "The early bird catches the worm.",
            "The pen is mightier than the sword." 
]

earlyBirdQuery =
    from sentence in stringList
    let words = sentence.split(" ")
    from word in words
    let w = word.lower()
    where w[0] == "a" || w[0] == "e" ||
          w[0] == "i" || w[0] == "o" ||
          w[0] == "u"
    select word

for v in earlyBirdQuery =>printf("'%s' starts with a vowel\n", v)

printf("\n================================================================\n\n")
ingredients =[
    { Name : "Sugar",  Calories: 500 },
    { Name : "Egg",    Calories: 100 },
    { Name : "Milk",   Calories: 150 },
    { Name : "Flour",  Calories: 50  },
    { Name : "Butter", Calories: 200 }
]

highCalDairyQuery =
    from i in ingredients
    let isDairy = i.Name == "Milk" || i.Name == "Butter"
    where i.Calories >= 150 && isDairy
    select i

for ingredient in highCalDairyQuery => println(ingredient.Name)

printf("\n================================================================\n\n")
//Data Source
words = [ "blueberry", "chimpanzee", "abacus", "banana", "apple", "cheese" ];
// Create the query.
wordGroups =
    from w in words
    group w by w[0];

// Execute the query.
for wordGroup in wordGroups 
{
    printf("Words that start with the letter '%s':\n", wordGroup.key);
    for word in wordGroup.value => printf("\t%s\n", word)
    println()
}
