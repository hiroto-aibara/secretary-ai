import { useState } from 'react'
import type { Card } from '../types'
import styles from './CardModal.module.css'

interface Props {
  card: Card
  onClose: () => void
  onSave: (updates: Partial<Card>) => void
  onArchive: () => void
  onDelete: () => void
}

export function CardModal({ card, onClose, onSave, onArchive, onDelete }: Props) {
  const [title, setTitle] = useState(card.title)
  const [description, setDescription] = useState(card.description)
  const [labelsText, setLabelsText] = useState((card.labels ?? []).join(', '))

  const handleSave = () => {
    const labels = labelsText
      .split(',')
      .map((l) => l.trim())
      .filter(Boolean)
    onSave({ title, description, labels })
  }

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

          <div className={styles.meta}>
            <span>List: {card.list}</span>
            <span>Created: {new Date(card.created_at).toLocaleDateString()}</span>
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
