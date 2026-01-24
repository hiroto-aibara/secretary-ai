import { useState, useEffect, useCallback, useRef } from 'react'
import type { Board as BoardType, Card as CardType } from './types'
import { api } from './hooks/useApi'
import { useWebSocket } from './hooks/useWebSocket'
import { Board } from './components/Board'
import { ArchiveView } from './components/ArchiveView'
import styles from './App.module.css'

function App() {
  const [boards, setBoards] = useState<BoardType[]>([])
  const [selectedBoardId, setSelectedBoardId] = useState<string | null>(null)
  const [cards, setCards] = useState<CardType[]>([])
  const [showArchive, setShowArchive] = useState(false)
  const [allCards, setAllCards] = useState<CardType[]>([])

  const selectedBoard = boards.find((b) => b.id === selectedBoardId) || null
  const selectedBoardIdRef = useRef(selectedBoardId)

  useEffect(() => {
    selectedBoardIdRef.current = selectedBoardId
  }, [selectedBoardId])

  const loadBoards = useCallback(async () => {
    const data = await api.boards.list()
    setBoards(data)
    if (data.length > 0 && !selectedBoardIdRef.current) {
      setSelectedBoardId(data[0].id)
    }
  }, [])

  const loadCards = useCallback(async () => {
    const boardId = selectedBoardIdRef.current
    if (!boardId) return
    const data = await api.cards.list(boardId)
    setCards(data)
  }, [])

  const loadAllCards = useCallback(async () => {
    const boardId = selectedBoardIdRef.current
    if (!boardId) return
    const data = await api.cards.list(boardId, true)
    setAllCards(data)
  }, [])

  useEffect(() => {
    let active = true
    api.boards.list().then((data) => {
      if (!active) return
      setBoards(data)
      if (data.length > 0 && !selectedBoardIdRef.current) {
        setSelectedBoardId(data[0].id)
      }
    })
    return () => {
      active = false
    }
  }, [])

  useEffect(() => {
    if (!selectedBoardId) return
    let active = true
    api.cards.list(selectedBoardId).then((data) => {
      if (!active) return
      setCards(data)
    })
    return () => {
      active = false
    }
  }, [selectedBoardId])

  useWebSocket((event) => {
    if (event.board_id === selectedBoardIdRef.current) {
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
    const boardId = selectedBoardIdRef.current
    if (!boardId) return
    await api.cards.archive(boardId, cardId, false)
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
        <Board board={selectedBoard} cards={cards} onRefresh={loadCards} />
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
