// Demonstrate json path handling
doc = {
	one: {
		"two": {
			"three": [1, 2, 3,]
		},
		"four": {
			"five": [ 11, 22, 33,]
		},
	},
}

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