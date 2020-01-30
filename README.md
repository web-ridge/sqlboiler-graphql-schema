# sqlboiler-graphql-schema

This websites generates a grapql schema based on the generated sqlboiler structs we do this because we want to support the sqlboiler aliasing in our schema. Generating the schema is a good way too add type safety to queries and filters and prevent too much manual typing.

You can edit your schema like you want later and re-generate if your database changes. This program will create a merge conflict with your existing schema so you can choose to accept/reject changes.

## How to run

`go run github.com/web-ridge/sqlboiler-graphql-schema -h`

## Options

```
USAGE:
    [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --input value   directory where the sqlboiler models are (default: "models")
   --output value  filepath for schema (default: "schema.graphql")
   --mutations     generate mutations for models
   --batch-update  generate batch update for models
   --batch-create  generate batch create for models
   --batch-delete  generate batch delete for models
   --help, -h      show help
   --version, -v   print the version
```

- Models
- Mutations (Followed best practices https://blog.apollographql.com/designing-graphql-mutations-e09de826ed97)

## Features

- Generating basic models (100%)
- Generating basic queries (100%)
- Generating filter for array queries (100%)
- Generating filters for relationships (100%)
- Generating input for mutations (100%)
- Generating payload for mutations (100%)
- Generating mutations (100%)
- Generating mutations for array models (0% WIP)
- Generating pagination for array models (0% WIP)

## Future roadmap

- Three way diff merge https://github.com/charlesvdv/go-three-way-merge to support manual updating the schema and re-generating
- Edges / connections
- Detecting when relationship is many to many
- Adding node from to many-to-many relationships
- Removing node from many-to-many relationships

Feel free to help and make a PR.

```
go run github.com/webridge-git/sqlboiler-graphql-schema
```

## Filtering lists (WIP)

This program generates type safe filters you can use in your frontend

### Search

```graphql
query(filter: {
    search: 'jan'
})
```

### This or that

```graphql
query(filter: {
    where: {
        name:{
            equalTo: 'Jan'
        }
        or: {
            name: {
                equalTo: 'Jannes',
            }
        }
    }
})

where: {
    id: {
        equalTo: 1
    }
    name: {
        startsWith: 'J',
    }
    where:{

    }
}
```

### (() or ())

````graphql
query(filter: {
    where: {
        or:{
            id: {
                equalTo: 1
            }
            name: {
                startsWith: 'J',
            }
            or:{
                id: {
                    equalTo: 2
                }
                name: {
                    startsWith: 'R',
                }
            }
        }
    }})

### Filter
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
````

## How to detect many to many (notes to myself)

```golang
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
```

### Possible detection

if table has only 2 relationship and table contains both tables
