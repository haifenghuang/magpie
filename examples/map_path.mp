// Demonstrate json path handling
doc = {
	"one": {
		"two": {
			"three": [1, 2, 3,],
			"six":(1,2,3)
		},
		"four": {
			"five": [ 11, 22, 33,]
		},
	},
}

printf("original: doc=[%v]\n", doc)

doc["one"]["two"]["three"][2] = 44
printf("doc[one][two][three][2]=%v\n", doc["one"]["two"]["three"][2])
//same as above
//doc.one.two.three[2] = 55
//printf("doc.one.two.three[2]=%v\n", doc.one.two.three[2])

doc["one"]["four"]["five"] = 4
printf("doc[one][four][five]=%v\n", doc["one"]["four"]["five"])

//runtime error
//doc["one"]["two"]["six"][2] = 44
//printf("doc[one][two][six][2]=%v\n", doc["one"]["two"]["six"][2])

doc["one"]["two"]["six"] = 78
printf("doc[one][two][six]=%v\n", doc["one"]["two"]["six"])

printf("after: doc=[%v]\n", doc)


prop = doc.getPath("one.two.three[2]")
println(prop) // result => 3

prop= doc.getPath("one.two.three")
println(prop) // result => [1, 2, 3]

prop= doc.getPath("one.two.three2[2]")
println(prop.message) // result => Property three2 does not exist


//Another example
out = {
  "one": {
    "two": {
      "three": [1, 2, 3,]
    },
    "four": {
      "five": [11, 22, 33,]
    },
  },
  "lala": {
    "two": {
      "three": {
        "four": {
          "five": {
            "six": {
              "seven":[
                nil,
                nil,
                nil,
                nil,
                nil,
                {
                  "eight": {
                    "nine": "test",
                  },
                },
              ],
            },
          },
        },
      },
    },
  },
}

prop = out.getPath("lala.two.three.four.five.six.seven[5].eight.nine")
println(prop) // result => "test",

out["lala"]["two"]["three"]["four"]["five"]["six"]["seven"][5]["eight"]["nine"]= "testing"
println(out["lala"]["two"]["three"]["four"]["five"]["six"]["seven"][5]["eight"]["nine"]) 
println(out)

out.lala.two.three.four.five.six.seven[5].eight.nine = "testing2"
println(out.lala.two.three.four.five.six.seven[5].eight.nine)
println(out)