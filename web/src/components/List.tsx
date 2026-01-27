import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import { useDroppable } from '@dnd-kit/core'
import type { Card as CardType, List as ListType } from '../types'
import { Card } from './Card'
import { useState, useRef, useEffect } from 'react'
import styles from './List.module.css'

interface Props {
  list: ListType
  cards: CardType[]
  onCardClick: (card: CardType) => void
  onAddCard: (title: string) => void
  onRename: (newName: string) => void
  onDelete: () => void
}

export function List({
  list,
  cards,
  onCardClick,
  onAddCard,
  onRename,
  onDelete,
}: Props) {
  const [adding, setAdding] = useState(false)
  const [title, setTitle] = useState('')
  const [editing, setEditing] = useState(false)
  const [editName, setEditName] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  const { setNodeRef } = useDroppable({ id: list.id })

  useEffect(() => {
    if (editing && inputRef.current) {
      inputRef.current.focus()
      inputRef.current.select()
    }
  }, [editing])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (title.trim()) {
      onAddCard(title.trim())
      setTitle('')
      setAdding(false)
    }
  }

  const handleStartEdit = () => {
    setEditName(list.name)
    setEditing(true)
  }

  const handleRenameSubmit = () => {
    const trimmed = editName.trim()
    if (trimmed && trimmed !== list.name) {
      onRename(trimmed)
    }
    setEditing(false)
  }

  const handleRenameKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleRenameSubmit()
    } else if (e.key === 'Escape') {
      setEditing(false)
    }
  }

  return (
    <div className={styles.list} ref={setNodeRef}>
      <div className={styles.headerRow}>
        {editing ? (
          <input
            ref={inputRef}
            className={styles.headerInput}
            value={editName}
            onChange={(e) => setEditName(e.target.value)}
            onBlur={handleRenameSubmit}
            onKeyDown={handleRenameKeyDown}
          />
        ) : (
          <div className={styles.header} onClick={handleStartEdit}>
            {list.name}
          </div>
        )}
        <button
          className={styles.deleteListBtn}
          onClick={onDelete}
          title="Delete list"
          aria-label="Delete list"
        >
          &times;
        </button>
      </div>
      <SortableContext
        items={cards.map((c) => c.id)}
        strategy={verticalListSortingStrategy}
      >
        <div className={styles.cards}>
          {cards.map((card) => (
            <Card key={card.id} card={card} onClick={() => onCardClick(card)} />
          ))}
        </div>
      </SortableContext>
      {adding ? (
        <form onSubmit={handleSubmit} className={styles.addForm}>
          <input
            autoFocus
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Card title..."
            className={styles.input}
          />
          <div className={styles.addActions}>
            <button type="submit" className={styles.addBtn}>
              Add
            </button>
            <button
              type="button"
              onClick={() => setAdding(false)}
              className={styles.cancelBtn}
            >
              Cancel
            </button>
          </div>
        </form>
      ) : (
        <button
          onClick={() => setAdding(true)}
          className={styles.addCardTrigger}
        >
          + Add card
        </button>
      )}
    </div>
  )
}
