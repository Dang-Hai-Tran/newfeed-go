# NewsFeed Go

A high-performance social media backend system built with Go, featuring advanced caching and performance optimization strategies.

## Features

- Clean Architecture (Domain, Repository, Usecase, Delivery layers)
- PostgreSQL database with GORM
- Redis caching with intelligent invalidation
- JWT-based authentication
- RESTful API with Gin framework
- Rate limiting
- CORS support
- Graceful shutdown
- Comprehensive logging
- Input validation
- Pagination support

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Redis

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Dang-Hai-Tran/newfeed-go.git
cd newfeed-go
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up your configuration in `config/config.yaml`

4. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### User Management
- `POST /v1/sessions` - Login
- `POST /v1/users` - Register
- `GET /v1/users/:user_id` - Get Profile
- `PUT /v1/users` - Update Profile
- `DELETE /v1/users` - Delete Profile
- `GET /v1/friends/:user_id` - Get Followers
- `POST /v1/friends/:user_id` - Follow User
- `DELETE /v1/friends/:user_id` - Unfollow User

### Post Management
- `GET /v1/posts/:post_id` - Get Post
- `POST /v1/posts` - Create Post
- `PUT /v1/posts/:post_id` - Update Post
- `DELETE /v1/posts/:post_id` - Delete Post
- `GET /v1/friends/:user_id/posts` - Get User Posts
- `GET /v1/users/:user_id/newsfeed` - Get Newsfeed

### Comment Management
- `GET /v1/posts/:post_id/comments` - Get Comments
- `POST /v1/posts/:post_id/comments` - Create Comment
- `PUT /v1/posts/:post_id/comments/:comment_id` - Update Comment
- `DELETE /v1/posts/:post_id/comments/:comment_id` - Delete Comment

### Like Management
- `GET /v1/posts/:post_id/likes` - Get Likes
- `POST /v1/posts/:post_id/likes` - Like Post
- `DELETE /v1/posts/:post_id/likes` - Unlike Post

## Architecture

The project follows Clean Architecture principles with the following layers:

1. Domain Layer - Core business logic and interfaces
2. Repository Layer - Data persistence and caching
3. Usecase Layer - Business logic implementation
4. Delivery Layer - HTTP handlers and middleware

## Performance Features

- Connection pooling (PostgreSQL and Redis)
- Intelligent caching layer
- Efficient database queries
- Pagination support
- Minimal database round trips
- Rate limiting
- Request/response optimization

## Security Features

- JWT token-based authentication
- Password hashing with bcrypt
- Rate limiting protection
- Input validation
- Resource authorization
- CORS protection

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
