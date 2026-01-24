import { useState, useEffect, useCallback } from 'react'
import type { Board as BoardType, Card as CardType } from './types'
import { useApi } from './hooks/useApi'
import { useWebSocket } from './hooks/useWebSocket'
import { Board } from './components/Board'
import { ArchiveView } from './components/ArchiveView'
import styles from './App.module.css'

function App() {
  const api = useApi()
  const [boards, setBoards] = useState<BoardType[]>([])
  const [selectedBoardId, setSelectedBoardId] = useState<string | null>(null)
  const [cards, setCards] = useState<CardType[]>([])
  const [showArchive, setShowArchive] = useState(false)
  const [allCards, setAllCards] = useState<CardType[]>([])

  const selectedBoard = boards.find((b) => b.id === selectedBoardId) || null

  const loadBoards = useCallback(async () => {
    const data = await api.boards.list()
    setBoards(data)
    if (data.length > 0 && !selectedBoardId) {
      setSelectedBoardId(data[0].id)
    }
  }, [api.boards, selectedBoardId])

  const loadCards = useCallback(async () => {
    if (!selectedBoardId) return
    const data = await api.cards.list(selectedBoardId)
    setCards(data)
  }, [api.cards, selectedBoardId])

  const loadAllCards = useCallback(async () => {
    if (!selectedBoardId) return
    const data = await api.cards.list(selectedBoardId, true)
    setAllCards(data)
  }, [api.cards, selectedBoardId])

  useEffect(() => {
    loadBoards()
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (selectedBoardId) {
      loadCards()
    }
  }, [selectedBoardId]) // eslint-disable-line react-hooks/exhaustive-deps

  const handleRefresh = useCallback(() => {
    loadCards()
  }, [loadCards])

  useWebSocket((event) => {
    if (event.board_id === selectedBoardId) {
      if (event.type === 'board_updated') {
        loadBoards()
      }
      loadCards()
    }
  })

  const handleShowArchive = async () => {
    await loadAllCards()
    setShowArchive(true)
  }

  const handleRestore = async (cardId: string) => {
    if (!selectedBoardId) return
    await api.cards.archive(selectedBoardId, cardId, false)
    await loadAllCards()
    loadCards()
  }

  return (
    <div className={styles.app}>
      <header className={styles.header}>
        <div className={styles.headerLeft}>
          <h1 className={styles.logo}>TaskMgr</h1>
          {boards.length > 1 && (
            <select
              className={styles.boardSelect}
              value={selectedBoardId || ''}
              onChange={(e) => setSelectedBoardId(e.target.value)}
            >
              {boards.map((b) => (
                <option key={b.id} value={b.id}>
                  {b.name}
                </option>
              ))}
            </select>
          )}
        </div>
        <button className={styles.archiveBtn} onClick={handleShowArchive}>
          Archive
        </button>
      </header>

      {selectedBoard ? (
        <Board board={selectedBoard} cards={cards} onRefresh={handleRefresh} />
      ) : (
        <div className={styles.empty}>
          <p>No boards found. Create one via API.</p>
        </div>
      )}

      {showArchive && (
        <ArchiveView
          cards={allCards}
          onRestore={handleRestore}
          onClose={() => setShowArchive(false)}
        />
      )}
    </div>
  )
}

export default App
