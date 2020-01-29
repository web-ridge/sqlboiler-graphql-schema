# sqlboiler-graphql-schema

This websites generates a grapql schema based on the generated sqlboiler structs.

- Models
- Mutations (Followed best practices https://blog.apollographql.com/designing-graphql-mutations-e09de826ed97)

## First goals

- Generating basic models + query + mutations

## Future roadmap

- Edges / connections
- Adding node from to many-to-many relationships
- Removing node from many-to-many relationships
- Merging with existing schema?

```
go run github.com/webridge-git/sqlboiler-graphql-schema
```
