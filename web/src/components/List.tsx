import {
  SortableContext,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { useDroppable } from '@dnd-kit/core'
import type { Card as CardType, List as ListType } from '../types'
import { Card } from './Card'
import { useState } from 'react'
import styles from './List.module.css'

interface Props {
  list: ListType
  cards: CardType[]
  onCardClick: (card: CardType) => void
  onAddCard: (title: string) => void
}

export function List({ list, cards, onCardClick, onAddCard }: Props) {
  const [adding, setAdding] = useState(false)
  const [title, setTitle] = useState('')

  const { setNodeRef } = useDroppable({ id: list.id })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (title.trim()) {
      onAddCard(title.trim())
      setTitle('')
      setAdding(false)
    }
  }

  return (
    <div className={styles.list} ref={setNodeRef}>
      <div className={styles.header}>{list.name}</div>
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
