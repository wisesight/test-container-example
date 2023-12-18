package com.example.demo_mongo.post;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.test.annotation.DirtiesContext;
import org.testcontainers.containers.MongoDBContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.springframework.boot.test.context.SpringBootTest.WebEnvironment;


@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@DirtiesContext
@Testcontainers
class PostControllerWithTestContainerTest {

    @Autowired
    private TestRestTemplate restTemplate;

    @Autowired
    private PostRepository postRepository;

    @Container
    @ServiceConnection
    static MongoDBContainer mongoDbContainer = new MongoDBContainer("mongo:6");

    @AfterEach
    public void clearData() {
        postRepository.deleteAll();
    }

    @Test
    @DisplayName("Success case with create a new post")
    void createNewPost() {
        // Arrange
        Post newPost = new Post("demo");
        // Act
        Post result = restTemplate.postForObject("/post", newPost, Post.class);
        // Assert
        assertNotNull(result.id());
        assertEquals("demo", result.title());
    }

    @Test
    @DisplayName("Success case with get all posts = 2 posts")
    void getAllPosts() {
        // Arrange
        Post newPost1 = new Post("demo 1");
        Post newPost2 = new Post("demo 2");
        restTemplate.postForObject("/post", newPost1, Post.class);
        restTemplate.postForObject("/post", newPost2, Post.class);

        // Act
        List<Post> posts = List.of(restTemplate.getForObject("/post", Post[].class));
        // Assert
        assertEquals(2, posts.size());
    }
}