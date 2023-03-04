# MonGear


[![tag](https://img.shields.io/github/tag/cascadiansw/mongear.svg)](https://github.com/samber/lo/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-%23007d9c)
[![Contributors](https://img.shields.io/github/contributors/cascadiansw/mongear)](https://github.com/samber/lo/graphs/contributors)
[![License](https://img.shields.io/github/license/cascadiansw/mongear)](./LICENSE)

MonGear is a light set of convenience functions for using Mongo
in Go. 


## Installation
Use the standard `go get github.com/cascadiansw/mongear`

## Functions
In this version, the functions revolve entirely around creation
of bson documents for filters and aggregation pipelines, using
plain 'ol strings instead of ugly looking bson documents.

### Query

Given a valid Mongo query/filter string, this returns a bson.D
document that can be directly fed to any Mongo Find*, Update, Replace, Delete
function.  For example:

```go
keys := ["Hello", "World"]

query, _ := mongear.Query(fmt.Sprintf(`{
    "key": {"$in": [%s]}
}`, keys))

cursor, _ := collection.Find(context.TODO(), query, nil)
```
Now, admitedly that's not much of a win.  But when you get into something like this,
it definitely is.

```go
query, _ := mongear.Query(`{
    "$and": [
        {
            "$or": [
                {
                    "user.id" : avogadrosothernumber
                },
                {
                    "retweeted_status.user.id": "Yre72kda3901"
                }
            ]
        },
        {
            "$and": [
                {
                    "date_posted": {
                        "$gt": somedate
                    }
                },
                {
                    "date_posted": {
                        "$lt": somedate
                    }
                }
            ]
        }
    ]
}`)
```

### Pipeline

Given a valid Mongo aggregation pipeline string, this returns a bson.D
document that can be directly fed to Mongo's Aggregation
function.  For example:

```go
keys := ["Hello", "World"]

pipeline, _ := mongear.Pipeline(fmt.Sprintf(`[
	{"$match": { "key": {"$in": [%s]}}},
	{"$lookup": {
		"from": "Users",
		"localField": "key",
		"foreignField": "foreignKey",
		"as": "users"
	}}
]`, keys))

cursor, _ := collection.Aggregate(context.TODO(), pipeline, nil)
```

### Stage

Sometimes it's useful to be able to build an aggregation pipeline one
stage at a time.  For this case, you can use the `Stage()` function.

Given a valid Mongo aggregation pipeline string, this returns a bson.D
document that can be directly fed to Mongo's Aggregation
function.  For example:

```go
keys := ["Hello", "World"]

p := mongo.Pipeline{}

p, match, _ = mongear.Stage(p, fmt.Sprintf(`
    {"$match": { "key": {"$in": [%s]}}}`, 
keys))

p, lookup, _ = mongear.Stage(p, `{
    "$lookup": {
        "from": "Users",
        "localField": "key",
        "foreignField": "foreignKey",
        "as": "users"
    }}`)

cursor, _ := collection.Aggregate(context.TODO(), pipeline, nil)
```

This function returns both the updated pipeline, well as a bson.D document
representing the stage itself, so you can do your own pipeline
management.  If you pass `nil` in the pipeline parameter, the pipeline
is ignored and only the stage is returned.