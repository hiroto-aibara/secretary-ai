import { useState } from 'react'
import styles from './AddList.module.css'

interface Props {
  onAdd: (name: string) => void
}

export function AddList({ onAdd }: Props) {
  const [adding, setAdding] = useState(false)
  const [name, setName] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (name.trim()) {
      onAdd(name.trim())
      setName('')
      setAdding(false)
    }
  }

  const handleCancel = () => {
    setName('')
    setAdding(false)
  }

  if (!adding) {
    return (
      <button className={styles.trigger} onClick={() => setAdding(true)}>
        + Add List
      </button>
    )
  }

  return (
    <div className={styles.container}>
      <form onSubmit={handleSubmit}>
        <input
          className={styles.input}
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="List name..."
          autoFocus
        />
        <div className={styles.actions}>
          <button
            type="submit"
            className={styles.addBtn}
            disabled={!name.trim()}
          >
            Add
          </button>
          <button
            type="button"
            className={styles.cancelBtn}
            onClick={handleCancel}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  )
}
