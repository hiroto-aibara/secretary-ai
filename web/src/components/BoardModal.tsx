import { useState } from 'react'
import styles from './BoardModal.module.css'

interface Props {
  onClose: () => void
  onCreate: (id: string, name: string, lists: string[]) => void
}

export function BoardModal({ onClose, onCreate }: Props) {
  const [name, setName] = useState('')
  const [listsText, setListsText] = useState('Todo, In Progress, Done')

  const generateId = (input: string): string => {
    return input
      .trim()
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')
  }

  const generatedId = generateId(name)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim() || !generatedId) return

    const lists = listsText
      .split(',')
      .map((l) => l.trim())
      .filter(Boolean)

    onCreate(generatedId, name.trim(), lists)
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
          <label className={styles.fieldLabel}>Board Name</label>
          <input
            className={styles.input}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g., Project Alpha"
            autoFocus
          />
          {generatedId && (
            <div className={styles.idPreview}>ID: {generatedId}</div>
          )}

          <label className={styles.fieldLabel}>
            Initial Lists (comma-separated)
          </label>
          <input
            className={styles.input}
            value={listsText}
            onChange={(e) => setListsText(e.target.value)}
            placeholder="e.g., Todo, In Progress, Done"
          />

          <div className={styles.footer}>
            <button
              type="button"
              className={styles.cancelBtn}
              onClick={onClose}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.createBtn}
              disabled={!name.trim() || !generatedId}
            >
              Create Board
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
