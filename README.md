# Let's Get Trivial

**Get a random trivia question**

```curl -X GET 'localhost:8080/trivium'```

**Add a new trivia question**

```curl -X POST -d '{"prompt":"How many books are in the Harry Potter series?","answer":"7"}' localhost:8080/trivium```

with data:

-`prompt`

-`answer`

-`answer_details`

-`attribution`

**Get a random trivia question from a specific themed trivia group ("channel")**

```curl -X GET 'localhost:8080/trivium?channel_name=computers'```

**Add a new trivia question to a specific themed trivia group ("channel")**

```curl -X POST -d '{"channel_name":"computers","prompt":"This moon of Jupiter shares its name with a computer networking term","answer":"Io"}' localhost:8080/trivium```
