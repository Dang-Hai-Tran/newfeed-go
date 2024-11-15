## System interface definition

- Login: POST v1/sessions {username, password} {msg}
- Sign up: POST v1/users {username, password, email, first_name, last_name, birthday} {msg}
- Edit profile: PUT v1/users {first_name, last_name, birthday, password} {msg}
- Delete user: DELETE v1/users/:user_id {msg}
- Get users: GET v1/users/:user_id {user}

- See follow list: GET v1/friends/:user_id {users}
- Follow user: POST v1/friends/:user_id {user_id} {msg}
- Unfollow user: DELETE v1/friends/:user_id {user_id} {msg}
- See user posts: GET v1/friends/:user_id/posts {posts}
- See post: GET v1/posts/:post_id {post} // post include text, image, comments, likes
- Create post: POST v1/posts {content} {msg} // content include text, image
- Edit post: PUT v1/posts/:post_id {content} {msg} // content include text, image
- Delete post: DELETE v1/posts/:post_id {msg}
- See comments: GET v1/posts/:post_id/comments {comments}
- Create comment: POST v1/posts/:post_id/comments {content} {comment}
- Edit comment: PUT v1/posts/:post_id/comments/:comment_id {content} {comment}
- Delete comment: DELETE v1/posts/:post_id/comments/:comment_id {msg}
- See likes: GET v1/posts/:post_id/likes {likes}
- Like post: POST v1/posts/:post_id/likes {msg}
- Unlike post: DELETE v1/posts/:post_id/likes {msg}

- Newsfeed: GET v1/users/:user_id/newsfeed {posts}
