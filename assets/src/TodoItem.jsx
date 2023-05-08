import React from 'react'
import { useMutation, gql } from '@apollo/client'
import { ALL_TODOS } from './Todos'
import classNames from 'classnames'

const DELETE_TODO = gql`
  mutation DeleteTodo($id: ID!) {
    removeTodo(id: $id) {
      id
    }
  }
`

const TOGGLE_TODO = gql`
  mutation UpdateTodo($id: ID!, $text: String!, $done: Boolean!) {
    # operation name
    updateTodo(input:{id: $id, text:$text, done: $done})
      # return fields
  {
      id
      text
      done
    }
  }
`

const TodoItem = (props) => {
  const [deleteTodo, { data }] = useMutation(DELETE_TODO, {
    refetchQueries: [{ query: ALL_TODOS }],
  })

  const [completeTodo, { data2 }] = useMutation(TOGGLE_TODO, {
    refetchQueries: [{ query: ALL_TODOS }],
    variables: {
      id: props.todo.id,
      text: props.todo.text,
      done: !props.todo.done,
    },
  })

  const names = classNames({
    todo: true,
    text: props.todo.text,
    done: props.todo.done,
  })
  return (
    <li className={names}>
      <div className="view">
        <input
          type="checkbox"
          className="toggle"
          checked={props.todo.done}
          onChange={() => completeTodo(props.todo.id, props.todo.text, props.todo.done)}
          data-cy="toggle"
        />
        <label>{`${props.todo.text}`}</label>
        <button
          className="destroy"
          data-cy="destroy"
          onClick={() => {
            deleteTodo({
              variables: {
                id: props.todo.id,
              },
            })
          }}
        />
      </div>
    </li>
  )
}
export default TodoItem
