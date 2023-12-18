# Spring Boot 3.2 Example Repository using TestContainer
This repository is structured to demonstrate the use of TestContainer in Java and Spring Boot 3.2 projects.

## Strucure of project
* PostController for RESTful APIs
* PostRepository for manage data in MongoDB

How to run test ?
```
$gradlew clean test
```

How to run Spring Boot server with Test Container ?
```
$gradlew clean test
```
Try to access RESTful APIs
* Get all post with `GET http://localhost:8080/post`
* Create new post with `POST http://localhost:8080/post`