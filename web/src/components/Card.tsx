import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import type { Card as CardType } from '../types'
import styles from './Card.module.css'

interface Props {
  card: CardType
  onClick: () => void
}

export function Card({ card, onClick }: Props) {
  const { attributes, listeners, setNodeRef, transform, transition } =
    useSortable({ id: card.id, data: { card } })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={styles.card}
      onClick={onClick}
      {...attributes}
      {...listeners}
    >
      <div className={styles.title}>{card.title}</div>
      {(card.labels ?? []).length > 0 && (
        <div className={styles.labels}>
          {(card.labels ?? []).map((label) => (
            <span key={label} className={styles.label}>
              {label}
            </span>
          ))}
        </div>
      )}
      {(card.todos ?? []).length > 0 &&
        (() => {
          const todos = card.todos ?? []
          const completed = todos.filter((t) => t.completed).length
          const total = todos.length
          const percent = (completed / total) * 100
          return (
            <div className={styles.todoProgress}>
              <div className={styles.progressBar}>
                <div
                  className={styles.progressFill}
                  style={{ width: `${percent}%` }}
                />
              </div>
              <span className={styles.progressText}>
                {completed}/{total}
              </span>
            </div>
          )
        })()}
    </div>
  )
}
