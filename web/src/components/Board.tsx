import { useState, useCallback } from 'react'
import {
  DndContext,
  type DragEndEvent,
  type DragOverEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import type { Board as BoardType, Card as CardType } from '../types'
import { List } from './List'
import { CardModal } from './CardModal'
import { useApi } from '../hooks/useApi'
import styles from './Board.module.css'

interface Props {
  board: BoardType
  cards: CardType[]
  onRefresh: () => void
}

export function Board({ board, cards, onRefresh }: Props) {
  const api = useApi()
  const [selectedCard, setSelectedCard] = useState<CardType | null>(null)
  const [dragOverList, setDragOverList] = useState<string | null>(null)

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  )

  const getCardsForList = useCallback(
    (listId: string) =>
      cards
        .filter((c) => c.list === listId && !c.archived)
        .sort((a, b) => a.order - b.order),
    [cards],
  )

  const handleDragOver = (event: DragOverEvent) => {
    const { over } = event
    if (over) {
      const overCard = cards.find((c) => c.id === over.id)
      setDragOverList(overCard ? overCard.list : (over.id as string))
    }
  }

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event
    setDragOverList(null)

    if (!over || active.id === over.id) return

    const activeCard = cards.find((c) => c.id === active.id)
    if (!activeCard) return

    // Determine target list
    const overCard = cards.find((c) => c.id === over.id)
    const targetList = overCard ? overCard.list : (over.id as string)

    // Determine order
    const targetCards = getCardsForList(targetList).filter(
      (c) => c.id !== activeCard.id,
    )
    let order = 0
    if (overCard) {
      const overIndex = targetCards.findIndex((c) => c.id === overCard.id)
      order = overIndex >= 0 ? overIndex : targetCards.length
    } else {
      order = targetCards.length
    }

    await api.cards.move(board.id, activeCard.id, targetList, order)
    onRefresh()
  }

  const handleAddCard = async (listId: string, title: string) => {
    await api.cards.create(board.id, { title, list: listId })
    onRefresh()
  }

  const handleSaveCard = async (updates: Partial<CardType>) => {
    if (!selectedCard) return
    await api.cards.update(board.id, selectedCard.id, updates)
    setSelectedCard(null)
    onRefresh()
  }

  const handleArchiveCard = async () => {
    if (!selectedCard) return
    await api.cards.archive(board.id, selectedCard.id, !selectedCard.archived)
    setSelectedCard(null)
    onRefresh()
  }

  const handleDeleteCard = async () => {
    if (!selectedCard) return
    await api.cards.delete(board.id, selectedCard.id)
    setSelectedCard(null)
    onRefresh()
  }

  return (
    <DndContext
      sensors={sensors}
      onDragOver={handleDragOver}
      onDragEnd={handleDragEnd}
    >
      <div className={styles.board} data-drag-over-list={dragOverList}>
        {board.lists.map((list) => (
          <List
            key={list.id}
            list={list}
            cards={getCardsForList(list.id)}
            onCardClick={setSelectedCard}
            onAddCard={(title) => handleAddCard(list.id, title)}
          />
        ))}
      </div>

      {selectedCard && (
        <CardModal
          card={selectedCard}
          onClose={() => setSelectedCard(null)}
          onSave={handleSaveCard}
          onArchive={handleArchiveCard}
          onDelete={handleDeleteCard}
        />
      )}
    </DndContext>
  )
}
