import { useState } from 'react'
import styles from './BoardModal.module.css'

interface Props {
  onClose: () => void
  onCreate: (id: string, name: string, lists: string[]) => void
}

export function BoardModal({ onClose, onCreate }: Props) {
  const [id, setId] = useState('')
  const [name, setName] = useState('')
  const [listsText, setListsText] = useState('Todo, In Progress, Done')

  const sanitizeId = (input: string): string => {
    return input
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!id.trim() || !name.trim()) return

    const sanitizedId = sanitizeId(id)
    if (!sanitizedId) return

    const lists = listsText
      .split(',')
      .map((l) => l.trim())
      .filter(Boolean)

    onCreate(sanitizedId, name.trim(), lists)
  }

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.header}>
          <h2 className={styles.title}>Create New Board</h2>
          <button className={styles.closeBtn} onClick={onClose}>
            &times;
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.body}>
          <label className={styles.fieldLabel}>Board ID</label>
          <input
            className={styles.input}
            value={id}
            onChange={(e) => setId(e.target.value)}
            placeholder="e.g., project-alpha"
            autoFocus
          />

          <label className={styles.fieldLabel}>Board Name</label>
          <input
            className={styles.input}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g., Project Alpha"
          />

          <label className={styles.fieldLabel}>Initial Lists (comma-separated)</label>
          <input
            className={styles.input}
            value={listsText}
            onChange={(e) => setListsText(e.target.value)}
            placeholder="e.g., Todo, In Progress, Done"
          />

          <div className={styles.footer}>
            <button type="button" className={styles.cancelBtn} onClick={onClose}>
              Cancel
            </button>
            <button
              type="submit"
              className={styles.createBtn}
              disabled={!id.trim() || !name.trim()}
            >
              Create Board
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
