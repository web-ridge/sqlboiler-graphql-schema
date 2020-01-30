# sqlboiler-graphql-schema

This websites generates a grapql schema based on the generated sqlboiler structs we do this because we want to support the sqlboiler aliasing in our schema.

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

### Possible detection

if table has only 2 relationship and table contains both tables

## Filtering lists (WIP)

```grapql
enum FilterTypeID {
IN
NOT_IN
EQUAL
NOT_EQUAL
}
enum FilterTypeBool {
TRUE
FALSE
NULL
}
enum FilterTypeInteger {
EQUAL
NOT_EQUAL
LOWER_THAN
LOWER_THAN_OR_EQUAL
IN
NOT_IN
}
enum FilterTypeString {
EQUAL
NOT_EQUAL
LOWER_THAN
LOWER_THAN_OR_EQUAL
IN
NOT_IN
}

enum FlowBlockColumn {
ID

}

type FlowBlockExpression {
Field:
Expressions: [FlowBlockExpression!]
}
type FlowBlockWhere {
Search: String
Expressions: [FlowBlockExpression!]!
}
```

```graphql
query(filter: {
    where: {
        id:{
            in: [1, 2, 3, 4]
            equalTo: 1
        },
        name: {
            like: "joe"
        },
        organizationId: {
            equalTo: 1,
        }
        or: {
            id:{
                in: []
                equalTo: 1
            },
             name: {
                like: "joe"
            },
        }
        and: {

        }
    }

})
```
