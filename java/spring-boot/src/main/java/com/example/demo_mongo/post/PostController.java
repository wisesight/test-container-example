package com.example.demo_mongo.post;

import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("post")
public class PostController {

    private final PostRepository postRepository;

    public PostController(PostRepository postRepository) {
        this.postRepository = postRepository;
    }

    @PostMapping
    public Post createNewPost(@RequestBody Post newPost) {
        return postRepository.save(newPost);
    }

    @GetMapping
    public List<Post> getAll() {
        return postRepository.findAll();
    }

}
