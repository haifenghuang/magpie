# Annotations

 Magpie also has very simple annotation support like java：

* Only method and property of class can have annotations\(not class itself, or other simple functions\)
* In the body of Annotation class, only support property, do not support methods.
* When use annotations, you must create an object.

 You could use `class @annotationName {}` to declare an annotation class. Magpie also include some builtin annotations:

* @Override annotation\(just like java's @Override\).
* @NotNull
* @NotEmpty

 Please see below example：

```swift
//Declare annotation class
//Note: In the body, you must use property, not method.
class @MinMaxValidator {
  property MinLength
  property MaxLength default 10 //Same as 'property MaxLength = 10'
}

//This is a marker annotation
class @NoSpaceValidator {}

class @DepartmentValidator {
  property Department
}

//The 'Request' class
class Request {
  @MinMaxValidator(MinLength=1)
  property FirstName; // getter and setter is implicit

  @NoSpaceValidator
  property LastName;

  @DepartmentValidator(Department=["Department of Education", "Department of Labors", "Department of Justice"])
  property Dept;
}

//This class is responsible for processing the annotation.
class RequestHandler {
  static fn handle(o) {
    props = o.getProperties()
    for p in props {
      annos = p.getAnnotations()
      for anno in annos {
        if anno.instanceOf(MinMaxValidator) {
          //p.value is the property real value.
          if len(p.value) > anno.MaxLength || len(p.value) < anno.MinLength {
            printf("Property '%s' is not valid!\n", p.name)
          }
        } elseif anno.instanceOf(NoSpaceValidator) {
          for c in p.value {
            if c == " " || c == "\t" {
              printf("Property '%s' is not valid!\n", p.name)
              break
            }
          }
        } elseif anno.instanceOf(DepartmentValidator) {
          found = false
          for d in anno.Department {
            if p.value == d {
              found = true
            }
          }
          if !found {
            printf("Property '%s' is not valid!\n", p.name)
          }
        }
      }
    }
  }
}

class RequestMain {
  static fn main() {
    request = new Request()
    request.FirstName = "Haifeng123456789"
    request.LastName = "Huang     "
    request.Dept = "Department of Justice"
    RequestHandler.handle(request)
  }
}

RequestMain.main()
```

 Below is the result：

```text
Property 'FirstName' is not valid!
Property 'LastName' is not valid!
```

