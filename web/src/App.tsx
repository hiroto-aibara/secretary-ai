import { useState, useEffect, useCallback, useRef } from 'react'
import type { Board as BoardType, Card as CardType } from './types'
import { api } from './hooks/useApi'
import { useWebSocket } from './hooks/useWebSocket'
import { Board } from './components/Board'
import { ArchiveView } from './components/ArchiveView'
import { BoardModal } from './components/BoardModal'
import { generateListId } from './utils/id'
import styles from './App.module.css'

function App() {
  const [boards, setBoards] = useState<BoardType[]>([])
  const [selectedBoardId, setSelectedBoardId] = useState<string | null>(null)
  const [cards, setCards] = useState<CardType[]>([])
  const [showArchive, setShowArchive] = useState(false)
  const [allCards, setAllCards] = useState<CardType[]>([])
  const [showBoardModal, setShowBoardModal] = useState(false)

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

  const handleCreateBoard = async (
    id: string,
    name: string,
    lists: string[],
  ) => {
    try {
      await api.boards.create({
        id,
        name,
        lists: lists.map((listName) => ({ id: generateListId(listName), name: listName })),
      })
      setShowBoardModal(false)
      const data = await api.boards.list()
      setBoards(data)
      setSelectedBoardId(id)
    } catch (err) {
      console.error('Failed to create board:', err)
    }
  }

  const handleDeleteBoard = async () => {
    if (!selectedBoardId) return
    const confirmed = window.confirm(
      `Are you sure you want to delete "${selectedBoard?.name}"? This action cannot be undone.`,
    )
    if (!confirmed) return

    try {
      await api.boards.delete(selectedBoardId)
      const data = await api.boards.list()
      setBoards(data)
      setSelectedBoardId(data.length > 0 ? data[0].id : null)
      setCards([])
    } catch (err) {
      console.error('Failed to delete board:', err)
    }
  }

  const handleBoardUpdate = async () => {
    await loadBoards()
    loadCards()
  }

  return (
    <div className={styles.app}>
      <header className={styles.header}>
        <div className={styles.headerLeft}>
          <h1 className={styles.logo}>TaskMgr</h1>
          {boards.length > 0 && (
            <div className={styles.boardSelector}>
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
              <button
                className={styles.deleteBoardBtn}
                onClick={handleDeleteBoard}
                title="Delete board"
              >
                &times;
              </button>
            </div>
          )}
          <button
            className={styles.newBoardBtn}
            onClick={() => setShowBoardModal(true)}
          >
            + New Board
          </button>
        </div>
        <button className={styles.archiveBtn} onClick={handleShowArchive}>
          Archive
        </button>
      </header>

      {selectedBoard ? (
        <Board board={selectedBoard} cards={cards} onRefresh={loadCards} onBoardUpdate={handleBoardUpdate} />
      ) : (
        <div className={styles.empty}>
          <p>No boards found. Click "+ New Board" to create one.</p>
        </div>
      )}

      {showArchive && (
        <ArchiveView
          cards={allCards}
          onRestore={handleRestore}
          onClose={() => setShowArchive(false)}
        />
      )}

      {showBoardModal && (
        <BoardModal
          onClose={() => setShowBoardModal(false)}
          onCreate={handleCreateBoard}
        />
      )}
    </div>
  )
}

export default App
