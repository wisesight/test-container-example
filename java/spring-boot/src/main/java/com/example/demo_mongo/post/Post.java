package com.example.demo_mongo.post;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.Date;
import java.util.UUID;

@Document("post")
public record Post(
        @Id
        String id,
        String title,
        String body,
        Date dateTime
) {
	public Post(String title) {
		this(UUID.randomUUID().toString(), title, null, null);
	}
}
