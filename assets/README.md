
# todo-graphql-example modified from https://github.com/bahmutov/todo-graphql-example

## App

Start server with `npm start`. You can find GraphQL playground at `http://localhost:3000/playground`

![App in action](images/app.gif)

Example asking for all todos

```
query {
  allTodos {
    id,
    title,
    completed
  }
}
```

Response

```json
{
  "data": {
    "allTodos": [
      {
        "id": "1",
        "title": "do something",
        "completed": false
      },
      {
        "id": "2",
        "title": "another",
        "completed": false
      }
    ]
  }
}
```

Example creating new todo object

```
mutation {
  createTodo(id: 2, title: "another", completed: false) {
    id
  }
}
```

Response

```json
{
  "data": {
    "createTodo": {
      "id": "2"
    }
  }
}
```

Example asking for a single todo (notice `id` argument)

```
query {
  Todo(id: 2) {
    id,
    title,
    completed
  }
}
```

Response

```json
{
  "data": {
    "Todo": {
      "id": "2",
      "title": "another",
      "completed": false
    }
  }
}
```

