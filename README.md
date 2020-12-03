# MOVED TO https://github.com/web-ridge/gqlgen-sqlboiler

## sqlboiler-graphql-schema

We want developers to be able to build software faster using modern tools like GraphQL, Golang, React Native without depending on commercial providers like Firebase or AWS Amplify.

This program generates a grapql schema based on the generated sqlboiler structs we do this because we want to support the sqlboiler aliasing in our schema. Generating the schema is a good way too add type safety to queries and filters and prevent too much manual typing.

You can edit your schema like you want later and re-generate if your database changes. This program will create a merge conflict with your existing schema so you can choose to accept/reject changes.

## How to run

`go run github.com/web-ridge/sqlboiler-graphql-schema`

## Before running

- Install prettier globally (https://prettier.io/ `yarn global add prettier`)
- Install git command line (required to do three way merging)

## Other related projects from webRidge

- https://github.com/web-ridge/gqlgen-sqlboiler (Generates converts between your qqlgen scheme and sqlboiler, and it will generate resolvers for the generated schema if you enable it!)
- https://github.com/web-ridge/graphql-schema-react-native-app (Generate a React Native (Web) app based on your GraphQL scheme, WIP.)

## Options

```
NAME:
   sqlboiler-graphql-schema

USAGE:
   sqlboiler-graphql-schema [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value              directory where the sqlboiler models are (default: "models")
   --output value             filepath for schema (default: "schema.graphql")
   --skip-input-fields value  input names which should be skipped: e.g. --skip-input-fields=userId --skip-input-fields=organizationId
   --mutations                generate mutations for models (default: true)
   --batch-update             generate batch update for models (default: true)
   --batch-create             generate batch create for models (default: true)
   --batch-delete             generate batch delete for models (default: true)
   --pagination               generate pagination support for models (default: "")
   --help, -h                 show help (default: false)
```


## Features
- [x] Support for manual updating the schema and re-generating (doing a three way merge https://github.com/charlesvdv/go-three-way-merge)
- [x] Generating basic models
- [x] Generating basic queries
- [x] Generating mutations (Followed best practices https://blog.apollographql.com/designing-graphql-mutations-e09de826ed97)
- [x] Generating filter for array queries (100%)
- [x] Generating filters for relationships (100%)
- [x] Generating input for mutations (100%)
- [x] Generating payload for mutations (100%)
- [x] Generating mutations (100%)
- [x] Generating mutations for array models (0% WIP)
- [x] Generating pagination for array models (20%, offset-based pagination done, TODO: cursor-based paginiation)

## Future roadmap

- [ ] Tests / snapshots
- [ ] Edges / connections
- [ ] Detecting when relationship is many to many
- [ ] Adding node from to many-to-many relationships
- [ ] Removing node from many-to-many relationships
- [ ] Supporting schema per model



Feel free to help and make a PR.

```
go run github.com/webridge-git/sqlboiler-graphql-schema
```

## Example result
https://github.com/web-ridge/gqlgen-sqlboiler-examples/blob/master/social-network/schema.graphql

```graphql

directive @isAuthenticated on FIELD_DEFINITION

type Comment {
  id: ID!
  content: String!
  post: Post!
  user: User!
  commentLikes: [CommentLike]
}

type CommentLike {
  id: ID!
  comment: Comment!
  user: User!
  likeType: String!
  createdAt: Int
}

type Friendship {
  id: ID!
  createdAt: Int
  users: [User]
}

type Image {
  id: ID!
  post: Post!
  views: Int
  originalUrl: String
  imageVariations: [ImageVariation]
}

type ImageVariation {
  id: ID!
  image: Image!
}

type Like {
  id: ID!
  post: Post!
  user: User!
  likeType: String!
  createdAt: Int
}

type Post {
  id: ID!
  content: String!
  user: User!
  comments: [Comment]
  images: [Image]
  likes: [Like]
}

type User {
  id: ID!
  firstName: String!
  lastName: String!
  email: String!
  comments: [Comment]
  commentLikes: [CommentLike]
  likes: [Like]
  posts: [Post]
  friendships: [Friendship]
}

input IDFilter {
  equalTo: ID
  notEqualTo: ID
  in: [ID!]
  notIn: [ID!]
}

input StringFilter {
  equalTo: String
  notEqualTo: String

  in: [String!]
  notIn: [String!]

  startWith: String
  notStartWith: String

  endWith: String
  notEndWith: String

  contain: String
  notContain: String

  startWithStrict: String # Camel sensitive
  notStartWithStrict: String # Camel sensitive
  endWithStrict: String # Camel sensitive
  notEndWithStrict: String # Camel sensitive
  containStrict: String # Camel sensitive
  notContainStrict: String # Camel sensitive
}

input IntFilter {
  equalTo: Int
  notEqualTo: Int
  lessThan: Int
  lessThanOrEqualTo: Int
  moreThan: Int
  moreThanOrEqualTo: Int
  in: [Int!]
  notIn: [Int!]
}

input FloatFilter {
  equalTo: Float
  notEqualTo: Float
  lessThan: Float
  lessThanOrEqualTo: Float
  moreThan: Float
  moreThanOrEqualTo: Float
  in: [Float!]
  notIn: [Float!]
}

input BooleanFilter {
  isTrue: Boolean
  isFalse: Boolean
  isNull: Boolean
}

input CommentFilter {
  search: String
  where: CommentWhere
}

input CommentWhere {
  id: IDFilter
  content: StringFilter
  post: PostWhere
  user: UserWhere
  commentLikes: CommentLikeWhere
  or: CommentWhere
  and: CommentWhere
}

input CommentLikeFilter {
  search: String
  where: CommentLikeWhere
}

input CommentLikeWhere {
  id: IDFilter
  comment: CommentWhere
  user: UserWhere
  likeType: StringFilter
  createdAt: IntFilter
  or: CommentLikeWhere
  and: CommentLikeWhere
}

input FriendshipFilter {
  search: String
  where: FriendshipWhere
}

input FriendshipWhere {
  id: IDFilter
  createdAt: IntFilter
  users: UserWhere
  or: FriendshipWhere
  and: FriendshipWhere
}

input ImageFilter {
  search: String
  where: ImageWhere
}

input ImageWhere {
  id: IDFilter
  post: PostWhere
  views: IntFilter
  originalUrl: StringFilter
  imageVariations: ImageVariationWhere
  or: ImageWhere
  and: ImageWhere
}

input ImageVariationFilter {
  search: String
  where: ImageVariationWhere
}

input ImageVariationWhere {
  id: IDFilter
  image: ImageWhere
  or: ImageVariationWhere
  and: ImageVariationWhere
}

input LikeFilter {
  search: String
  where: LikeWhere
}

input LikeWhere {
  id: IDFilter
  post: PostWhere
  user: UserWhere
  likeType: StringFilter
  createdAt: IntFilter
  or: LikeWhere
  and: LikeWhere
}

input PostFilter {
  search: String
  where: PostWhere
}

input PostWhere {
  id: IDFilter
  content: StringFilter
  user: UserWhere
  comments: CommentWhere
  images: ImageWhere
  likes: LikeWhere
  or: PostWhere
  and: PostWhere
}

input UserFilter {
  search: String
  where: UserWhere
}

input UserWhere {
  id: IDFilter
  firstName: StringFilter
  lastName: StringFilter
  email: StringFilter
  comments: CommentWhere
  commentLikes: CommentLikeWhere
  likes: LikeWhere
  posts: PostWhere
  friendships: FriendshipWhere
  or: UserWhere
  and: UserWhere
}

type Query {
  comment(id: ID!): Comment! @isAuthenticated
  comments(filter: CommentFilter): [Comment!]! @isAuthenticated
  commentLike(id: ID!): CommentLike! @isAuthenticated
  commentLikes(filter: CommentLikeFilter): [CommentLike!]! @isAuthenticated
  friendship(id: ID!): Friendship! @isAuthenticated
  friendships(filter: FriendshipFilter): [Friendship!]! @isAuthenticated
  image(id: ID!): Image! @isAuthenticated
  images(filter: ImageFilter): [Image!]! @isAuthenticated
  imageVariation(id: ID!): ImageVariation! @isAuthenticated
  imageVariations(filter: ImageVariationFilter): [ImageVariation!]!
    @isAuthenticated
  like(id: ID!): Like! @isAuthenticated
  likes(filter: LikeFilter): [Like!]! @isAuthenticated
  post(id: ID!): Post! @isAuthenticated
  posts(filter: PostFilter): [Post!]! @isAuthenticated
  user(id: ID!): User! @isAuthenticated
  users(filter: UserFilter): [User!]! @isAuthenticated
}

input CommentCreateInput {
  content: String!
  postId: ID!
}

input CommentUpdateInput {
  content: String
  postId: ID
}

input CommentsCreateInput {
  comments: [CommentCreateInput!]!
}

type CommentPayload {
  comment: Comment!
}

type CommentDeletePayload {
  id: ID!
}

type CommentsPayload {
  comments: [Comment!]!
}

type CommentsDeletePayload {
  ids: [ID!]!
}

type CommentsUpdatePayload {
  ok: Boolean!
}

input CommentLikeCreateInput {
  commentId: ID!
  likeType: String!
  createdAt: Int
}

input CommentLikeUpdateInput {
  commentId: ID
  likeType: String
  createdAt: Int
}

input CommentLikesCreateInput {
  commentLikes: [CommentLikeCreateInput!]!
}

type CommentLikePayload {
  commentLike: CommentLike!
}

type CommentLikeDeletePayload {
  id: ID!
}

type CommentLikesPayload {
  commentLikes: [CommentLike!]!
}

type CommentLikesDeletePayload {
  ids: [ID!]!
}

type CommentLikesUpdatePayload {
  ok: Boolean!
}

input FriendshipCreateInput {
  createdAt: Int
}

input FriendshipUpdateInput {
  createdAt: Int
}

input FriendshipsCreateInput {
  friendships: [FriendshipCreateInput!]!
}

type FriendshipPayload {
  friendship: Friendship!
}

type FriendshipDeletePayload {
  id: ID!
}

type FriendshipsPayload {
  friendships: [Friendship!]!
}

type FriendshipsDeletePayload {
  ids: [ID!]!
}

type FriendshipsUpdatePayload {
  ok: Boolean!
}

input ImageCreateInput {
  postId: ID!
  views: Int
  originalUrl: String
}

input ImageUpdateInput {
  postId: ID
  views: Int
  originalUrl: String
}

input ImagesCreateInput {
  images: [ImageCreateInput!]!
}

type ImagePayload {
  image: Image!
}

type ImageDeletePayload {
  id: ID!
}

type ImagesPayload {
  images: [Image!]!
}

type ImagesDeletePayload {
  ids: [ID!]!
}

type ImagesUpdatePayload {
  ok: Boolean!
}

input ImageVariationCreateInput {
  imageId: ID!
}

input ImageVariationUpdateInput {
  imageId: ID
}

input ImageVariationsCreateInput {
  imageVariations: [ImageVariationCreateInput!]!
}

type ImageVariationPayload {
  imageVariation: ImageVariation!
}

type ImageVariationDeletePayload {
  id: ID!
}

type ImageVariationsPayload {
  imageVariations: [ImageVariation!]!
}

type ImageVariationsDeletePayload {
  ids: [ID!]!
}

type ImageVariationsUpdatePayload {
  ok: Boolean!
}

input LikeCreateInput {
  postId: ID!
  likeType: String!
  createdAt: Int
}

input LikeUpdateInput {
  postId: ID
  likeType: String
  createdAt: Int
}

input LikesCreateInput {
  likes: [LikeCreateInput!]!
}

type LikePayload {
  like: Like!
}

type LikeDeletePayload {
  id: ID!
}

type LikesPayload {
  likes: [Like!]!
}

type LikesDeletePayload {
  ids: [ID!]!
}

type LikesUpdatePayload {
  ok: Boolean!
}

input PostCreateInput {
  content: String!
}

input PostUpdateInput {
  content: String
}

input PostsCreateInput {
  posts: [PostCreateInput!]!
}

type PostPayload {
  post: Post!
}

type PostDeletePayload {
  id: ID!
}

type PostsPayload {
  posts: [Post!]!
}

type PostsDeletePayload {
  ids: [ID!]!
}

type PostsUpdatePayload {
  ok: Boolean!
}

input UserCreateInput {
  firstName: String!
  lastName: String!
  email: String!
}

input UserUpdateInput {
  firstName: String
  lastName: String
  email: String
}

input UsersCreateInput {
  users: [UserCreateInput!]!
}

type UserPayload {
  user: User!
}

type UserDeletePayload {
  id: ID!
}

type UsersPayload {
  users: [User!]!
}

type UsersDeletePayload {
  ids: [ID!]!
}

type UsersUpdatePayload {
  ok: Boolean!
}

type Mutation {
  createComment(input: CommentCreateInput!): CommentPayload! @isAuthenticated
  createComments(input: CommentsCreateInput!): CommentsPayload! @isAuthenticated
  updateComment(id: ID!, input: CommentUpdateInput!): CommentPayload!
    @isAuthenticated
  updateComments(
    filter: CommentFilter
    input: CommentUpdateInput!
  ): CommentsUpdatePayload! @isAuthenticated
  deleteComment(id: ID!): CommentDeletePayload! @isAuthenticated
  deleteComments(filter: CommentFilter): CommentsDeletePayload! @isAuthenticated
  createCommentLike(input: CommentLikeCreateInput!): CommentLikePayload!
    @isAuthenticated
  createCommentLikes(input: CommentLikesCreateInput!): CommentLikesPayload!
    @isAuthenticated
  updateCommentLike(
    id: ID!
    input: CommentLikeUpdateInput!
  ): CommentLikePayload! @isAuthenticated
  updateCommentLikes(
    filter: CommentLikeFilter
    input: CommentLikeUpdateInput!
  ): CommentLikesUpdatePayload! @isAuthenticated
  deleteCommentLike(id: ID!): CommentLikeDeletePayload! @isAuthenticated
  deleteCommentLikes(filter: CommentLikeFilter): CommentLikesDeletePayload!
    @isAuthenticated
  createFriendship(input: FriendshipCreateInput!): FriendshipPayload!
    @isAuthenticated
  createFriendships(input: FriendshipsCreateInput!): FriendshipsPayload!
    @isAuthenticated
  updateFriendship(id: ID!, input: FriendshipUpdateInput!): FriendshipPayload!
    @isAuthenticated
  updateFriendships(
    filter: FriendshipFilter
    input: FriendshipUpdateInput!
  ): FriendshipsUpdatePayload! @isAuthenticated
  deleteFriendship(id: ID!): FriendshipDeletePayload! @isAuthenticated
  deleteFriendships(filter: FriendshipFilter): FriendshipsDeletePayload!
    @isAuthenticated
  createImage(input: ImageCreateInput!): ImagePayload! @isAuthenticated
  createImages(input: ImagesCreateInput!): ImagesPayload! @isAuthenticated
  updateImage(id: ID!, input: ImageUpdateInput!): ImagePayload! @isAuthenticated
  updateImages(
    filter: ImageFilter
    input: ImageUpdateInput!
  ): ImagesUpdatePayload! @isAuthenticated
  deleteImage(id: ID!): ImageDeletePayload! @isAuthenticated
  deleteImages(filter: ImageFilter): ImagesDeletePayload! @isAuthenticated
  createImageVariation(
    input: ImageVariationCreateInput!
  ): ImageVariationPayload! @isAuthenticated
  createImageVariations(
    input: ImageVariationsCreateInput!
  ): ImageVariationsPayload! @isAuthenticated
  updateImageVariation(
    id: ID!
    input: ImageVariationUpdateInput!
  ): ImageVariationPayload! @isAuthenticated
  updateImageVariations(
    filter: ImageVariationFilter
    input: ImageVariationUpdateInput!
  ): ImageVariationsUpdatePayload! @isAuthenticated
  deleteImageVariation(id: ID!): ImageVariationDeletePayload! @isAuthenticated
  deleteImageVariations(
    filter: ImageVariationFilter
  ): ImageVariationsDeletePayload! @isAuthenticated
  createLike(input: LikeCreateInput!): LikePayload! @isAuthenticated
  createLikes(input: LikesCreateInput!): LikesPayload! @isAuthenticated
  updateLike(id: ID!, input: LikeUpdateInput!): LikePayload! @isAuthenticated
  updateLikes(filter: LikeFilter, input: LikeUpdateInput!): LikesUpdatePayload!
    @isAuthenticated
  deleteLike(id: ID!): LikeDeletePayload! @isAuthenticated
  deleteLikes(filter: LikeFilter): LikesDeletePayload! @isAuthenticated
  createPost(input: PostCreateInput!): PostPayload! @isAuthenticated
  createPosts(input: PostsCreateInput!): PostsPayload! @isAuthenticated
  updatePost(id: ID!, input: PostUpdateInput!): PostPayload! @isAuthenticated
  updatePosts(filter: PostFilter, input: PostUpdateInput!): PostsUpdatePayload!
    @isAuthenticated
  deletePost(id: ID!): PostDeletePayload! @isAuthenticated
  deletePosts(filter: PostFilter): PostsDeletePayload! @isAuthenticated
  createUser(input: UserCreateInput!): UserPayload! @isAuthenticated
  createUsers(input: UsersCreateInput!): UsersPayload! @isAuthenticated
  updateUser(id: ID!, input: UserUpdateInput!): UserPayload! @isAuthenticated
  updateUsers(filter: UserFilter, input: UserUpdateInput!): UsersUpdatePayload!
    @isAuthenticated
  deleteUser(id: ID!): UserDeletePayload! @isAuthenticated
  deleteUsers(filter: UserFilter): UsersDeletePayload! @isAuthenticated
}
```

## Donate

Did we save you a lot of time? Please consider a donation so we can invest more time in this library: [![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donate_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=7B9KKQLXTEW9Q&source=url)
