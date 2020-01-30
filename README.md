# sqlboiler-graphql-schema

This websites generates a grapql schema based on the generated sqlboiler structs.

- Models
- Mutations (Followed best practices https://blog.apollographql.com/designing-graphql-mutations-e09de826ed97)

## First goals

- Generating basic models + query + mutations

## Future roadmap

- Edges / connections
- Detecting when relationship is many to many
- Adding node from to many-to-many relationships
- Removing node from many-to-many relationships
- Merging with existing schema?

Feel free to help and make a PR.

```
go run github.com/webridge-git/sqlboiler-graphql-schema
```

## How to detect many to many (notes to myself)

User {
UserOrganizations []UserOrganization
Posts []Post
}

Post {
User User
Message string
}

UserOrganization {
User User
Organization Organization
}

Organization {
Users []User
}

## Possible detection

if table has only 2 relationship and table contains both tables
