package com.example.demo_mongo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.context.annotation.Bean;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.testcontainers.containers.MongoDBContainer;
import org.testcontainers.utility.DockerImageName;

@TestConfiguration(proxyBeanMethods = false)
public class TestDemoMongoApplication {

	@Bean
	@ServiceConnection
	MongoDBContainer mongoDbContainer(DynamicPropertyRegistry properties) {
		MongoDBContainer container = new MongoDBContainer("mongo:6");
		properties.add("spring.data.mongodb.host", container::getHost);
		properties.add("spring.data.mongodb.port", container::getFirstMappedPort);
		return container;
	}

	public static void main(String[] args) {
		SpringApplication.from(DemoMongoApplication::main)
				.with(TestDemoMongoApplication.class)
				.run(args);
	}

}
