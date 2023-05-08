

### Create a Todo:
```
# I Specify that it's a mutation
mutation {
  # I invoke the createTodo mutation and pass it the input
  createTodo(input: {
    text: "This is a new todo"
  })
  # Specify which properties I want from the return value
  {
    id
    text
    done
  }
}
```
### Update an existing Todo:
```
# I Specify that it's a mutation
mutation {
    # I invoke the updateTodo mutation and pass it the input
updateTodo(input:{
  id:"43f2d64e-65f8-48f6-ac16-66c59cb68fa8",
  text: "This is a new todo",
  done: true
})
# Specify which properties I want from the return value
  {
    id
    text
    done
  }  
}
```
### Query Single Todo
```
# I Specify that it's a query
query {
    # I invoke the getTodo query and pass it the input
getTodo(todoId:"43f2d64e-65f8-48f6-ac16-66c59cb68fa8")
# Specify which properties I want from the return value
  {
    id
    text
    done
  }  
}
```
### Get All Todos
```
# I Specify that it's a query
query {
    # I invoke the getTodo query and pass it the input
getTodos
# Specify which properties I want from the return value
  {
    id
    text
    done
  }
}
```
### Delete a Single Todo
```
# I Specify that it's a mutation
mutation DeleteTodo{
    # I invoke the deleteTodo query and pass it the input
deleteTodo(todoId:"43f2d64e-65f8-48f6-ac16-66c59cb68fa8") 
}
```



// To get single ToDo item by ID
curl -g 'http://localhost:8081/query?query={todo(id:"1"){id,text,done}}'

// To create a ToDo item
curl -g 'http://localhost:8081/query?query=mutation+_{createTodo(text:"My+new+todo"){id,text,done}}'

// To get a list of ToDo items
curl -g 'http://localhost:8081/query?query={todoList{id,text,done}}'

// To update a ToDo
curl -g 'http://localhost:8081/query?query=mutation+_{updateTodo(id:"1",changes:{text:"My+new+todo+updated",done:true}){id,text,done}}'