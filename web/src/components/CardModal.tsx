import { useState } from 'react'
import type { Card, TodoItem } from '../types'
import { generateTodoId } from '../utils/id'
import styles from './CardModal.module.css'

interface Props {
  card: Card
  onClose: () => void
  onSave: (updates: Partial<Card>) => void
  onArchive: () => void
  onDelete: () => void
}

export function CardModal({
  card,
  onClose,
  onSave,
  onArchive,
  onDelete,
}: Props) {
  const [title, setTitle] = useState(card.title)
  const [description, setDescription] = useState(card.description)
  const [labelsText, setLabelsText] = useState((card.labels ?? []).join(', '))
  const [todos, setTodos] = useState<TodoItem[]>(card.todos ?? [])
  const [newTodoText, setNewTodoText] = useState('')
  const [editingTodoId, setEditingTodoId] = useState<string | null>(null)
  const [editingTodoText, setEditingTodoText] = useState('')

  const handleSave = () => {
    const labels = [
      ...new Set(
        labelsText
          .split(',')
          .map((l) => l.trim())
          .filter(Boolean),
      ),
    ]
    onSave({ title, description, labels, todos })
  }

  const handleAddTodo = () => {
    if (!newTodoText.trim()) return
    const newTodo: TodoItem = {
      id: generateTodoId(),
      text: newTodoText.trim(),
      completed: false,
    }
    setTodos([...todos, newTodo])
    setNewTodoText('')
  }

  const handleToggleTodo = (id: string) => {
    setTodos(
      todos.map((t) => (t.id === id ? { ...t, completed: !t.completed } : t)),
    )
  }

  const handleDeleteTodo = (id: string) => {
    setTodos(todos.filter((t) => t.id !== id))
  }

  const handleStartEditTodo = (todo: TodoItem) => {
    setEditingTodoId(todo.id)
    setEditingTodoText(todo.text)
  }

  const handleSaveEditTodo = () => {
    if (!editingTodoId) return
    setTodos(
      todos.map((t) =>
        t.id === editingTodoId ? { ...t, text: editingTodoText.trim() } : t,
      ),
    )
    setEditingTodoId(null)
    setEditingTodoText('')
  }

  const handleCancelEditTodo = () => {
    setEditingTodoId(null)
    setEditingTodoText('')
  }

  const completedCount = todos.filter((t) => t.completed).length
  const totalCount = todos.length

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.header}>
          <input
            className={styles.titleInput}
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <button className={styles.closeBtn} onClick={onClose}>
            &times;
          </button>
        </div>

        <div className={styles.body}>
          <label className={styles.fieldLabel}>Description</label>
          <textarea
            className={styles.textarea}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={4}
          />

          <label className={styles.fieldLabel}>Labels (comma-separated)</label>
          <input
            className={styles.input}
            value={labelsText}
            onChange={(e) => setLabelsText(e.target.value)}
          />

          <div className={styles.checklistHeader}>
            <label className={styles.fieldLabel}>Checklist</label>
            {totalCount > 0 && (
              <span className={styles.checklistProgress}>
                {completedCount}/{totalCount} done
              </span>
            )}
          </div>

          <div className={styles.todoList}>
            {todos.map((todo) => (
              <div key={todo.id} className={styles.todoItem}>
                <input
                  type="checkbox"
                  checked={todo.completed}
                  onChange={() => handleToggleTodo(todo.id)}
                  className={styles.todoCheckbox}
                />
                {editingTodoId === todo.id ? (
                  <div className={styles.todoEditContainer}>
                    <input
                      type="text"
                      value={editingTodoText}
                      onChange={(e) => setEditingTodoText(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.nativeEvent.isComposing) return
                        if (e.key === 'Enter') handleSaveEditTodo()
                        if (e.key === 'Escape') handleCancelEditTodo()
                      }}
                      className={styles.todoEditInput}
                      autoFocus
                    />
                    <button
                      onClick={handleSaveEditTodo}
                      className={styles.todoEditSaveBtn}
                    >
                      ✓
                    </button>
                    <button
                      onClick={handleCancelEditTodo}
                      className={styles.todoEditCancelBtn}
                    >
                      ✕
                    </button>
                  </div>
                ) : (
                  <>
                    <span
                      className={`${styles.todoText} ${todo.completed ? styles.todoCompleted : ''}`}
                      onClick={() => handleStartEditTodo(todo)}
                    >
                      {todo.text}
                    </span>
                    <button
                      onClick={() => handleDeleteTodo(todo.id)}
                      className={styles.todoDeleteBtn}
                    >
                      ×
                    </button>
                  </>
                )}
              </div>
            ))}
          </div>

          <div className={styles.todoAddForm}>
            <input
              type="text"
              placeholder="Add a new todo..."
              value={newTodoText}
              onChange={(e) => setNewTodoText(e.target.value)}
              onKeyDown={(e) =>
                e.key === 'Enter' &&
                !e.nativeEvent.isComposing &&
                handleAddTodo()
              }
              className={styles.todoAddInput}
            />
            <button onClick={handleAddTodo} className={styles.todoAddBtn}>
              Add
            </button>
          </div>

          <div className={styles.meta}>
            <span>List: {card.list}</span>
            <span>
              Created: {new Date(card.created_at).toLocaleDateString()}
            </span>
          </div>
        </div>

        <div className={styles.footer}>
          <button className={styles.saveBtn} onClick={handleSave}>
            Save
          </button>
          <button className={styles.archiveBtn} onClick={onArchive}>
            {card.archived ? 'Restore' : 'Archive'}
          </button>
          <button className={styles.deleteBtn} onClick={onDelete}>
            Delete
          </button>
        </div>
      </div>
    </div>
  )
}
