import type { Card as CardType } from '../types'
import styles from './ArchiveView.module.css'

interface Props {
  cards: CardType[]
  onRestore: (cardId: string) => void
  onClose: () => void
}

export function ArchiveView({ cards, onRestore, onClose }: Props) {
  const archivedCards = cards.filter((c) => c.archived)

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.panel} onClick={(e) => e.stopPropagation()}>
        <div className={styles.header}>
          <h3 className={styles.title}>Archived Cards</h3>
          <button className={styles.closeBtn} onClick={onClose}>
            &times;
          </button>
        </div>
        <div className={styles.list}>
          {archivedCards.length === 0 ? (
            <p className={styles.empty}>No archived cards</p>
          ) : (
            archivedCards.map((card) => (
              <div key={card.id} className={styles.item}>
                <div className={styles.cardTitle}>{card.title}</div>
                <button
                  className={styles.restoreBtn}
                  onClick={() => onRestore(card.id)}
                >
                  Restore
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
