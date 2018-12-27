# Static members/methods/properties

You could use `static` keyword to declare static members/methods/properties.

```swift
class Test
{
   static let x = 0;
   static let y = 5;

   static fn Main()
   {
      println(Test.x);
      println(Test.y);

      Test.x = 99;
      println(Test.x);
   }
}

Test.Main()
```

{% hint style="info" %}
 Non-static variable/method/property could access static variable/method/property. On the other hand, static variable/method/property cannot access Non-static variable/method/property.
{% endhint %}

