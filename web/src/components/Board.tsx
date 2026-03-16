import { useState, useCallback, useRef, useMemo } from 'react'
import {
  DndContext,
  DragOverlay,
  closestCenter,
  type DragStartEvent,
  type DragEndEvent,
  type DragMoveEvent,
  PointerSensor,
  useSensor,
  useSensors,
  type CollisionDetection,
} from '@dnd-kit/core'
import {
  SortableContext,
  horizontalListSortingStrategy,
  arrayMove,
} from '@dnd-kit/sortable'
import type {
  Board as BoardType,
  Card as CardType,
  List as ListType,
} from '../types'
import { List } from './List'
import { CardModal } from './CardModal'
import { AddList } from './AddList'
import { api } from '../hooks/useApi'
import { generateUniqueListId } from '../utils/id'
import styles from './Board.module.css'
import cardStyles from './Card.module.css'

function createCollisionDetection(listIds: Set<string>): CollisionDetection {
  return (args) => {
    const dragType = args.active.data.current?.type

    if (dragType === 'list') {
      const filtered = {
        ...args,
        droppableContainers: args.droppableContainers.filter((c) =>
          listIds.has(c.id as string),
        ),
      }
      return closestCenter(filtered)
    }

    const collisions = closestCenter(args)
    const cardCollision = collisions.find(
      (c) => !listIds.has(c.id as string) && c.id !== args.active.id,
    )
    if (cardCollision) return [cardCollision]

    const listCollision = collisions.find((c) => listIds.has(c.id as string))
    return listCollision ? [listCollision] : collisions
  }
}

interface Props {
  board: BoardType
  cards: CardType[]
  onRefresh: () => void
  onBoardUpdate: () => void
}

export function Board({ board, cards, onRefresh, onBoardUpdate }: Props) {
  const [selectedCard, setSelectedCard] = useState<CardType | null>(null)
  const [activeCard, setActiveCard] = useState<CardType | null>(null)
  const [activeList, setActiveList] = useState<ListType | null>(null)
  const [dragCards, setDragCards] = useState<CardType[] | null>(null)
  const [dragLists, setDragLists] = useState<ListType[] | null>(null)
  const dragCardsRef = useRef<CardType[] | null>(null)
  const lastMoveRef = useRef<{ overId: string; insertAfter: boolean } | null>(
    null,
  )
  const justDraggedRef = useRef(false)

  const effectiveCards = dragCards ?? cards
  const effectiveLists = dragLists ?? board.lists

  const listIds = useMemo(
    () => new Set(effectiveLists.map((l) => l.id)),
    [effectiveLists],
  )
  const collisionDetection = useMemo(
    () => createCollisionDetection(listIds),
    [listIds],
  )

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  )

  const getCardsForList = useCallback(
    (listId: string) =>
      (effectiveCards ?? [])
        .filter((c) => c.list === listId && !c.archived)
        .sort((a, b) => a.order - b.order),
    [effectiveCards],
  )

  const handleCardClick = useCallback((card: CardType) => {
    if (justDraggedRef.current) return
    setSelectedCard(card)
  }, [])

  const handleDragStart = (event: DragStartEvent) => {
    const dragType = event.active.data.current?.type

    if (dragType === 'list') {
      const list = effectiveLists.find((l) => l.id === event.active.id) ?? null
      setActiveList(list)
      justDraggedRef.current = true
      return
    }

    const card = cards.find((c) => c.id === event.active.id) ?? null
    setActiveCard(card)
    const snapshot = [...cards]
    dragCardsRef.current = snapshot
    setDragCards(snapshot)
    lastMoveRef.current = null
    justDraggedRef.current = true
  }

  const handleDragMove = (event: DragMoveEvent) => {
    if (event.active.data.current?.type === 'list') return

    const { active, over } = event
    if (!over || active.id === over.id) return

    const activeId = active.id as string
    const overId = over.id as string

    const activeMidY =
      (active.rect.current.translated?.top ?? 0) +
      (active.rect.current.translated?.height ?? 0) / 2

    const isOverList = listIds.has(overId)
    const insertAfter =
      !isOverList && activeMidY > over.rect.top + over.rect.height / 2

    const last = lastMoveRef.current
    if (last && last.overId === overId && last.insertAfter === insertAfter) {
      return
    }
    lastMoveRef.current = { overId, insertAfter }

    setDragCards((prev) => {
      if (!prev) return prev
      const activeCardItem = prev.find((c) => c.id === activeId)
      if (!activeCardItem) return prev

      const overCard = prev.find((c) => c.id === overId)
      const targetList = overCard ? overCard.list : overId

      const next = prev.filter((c) => c.id !== activeId)
      const updatedCard = { ...activeCardItem, list: targetList }

      if (overCard) {
        const overIndex = next.findIndex((c) => c.id === overId)
        const insertIndex = insertAfter ? overIndex + 1 : overIndex
        next.splice(
          insertIndex >= 0 ? insertIndex : next.length,
          0,
          updatedCard,
        )
      } else {
        // 空リストへのドロップ
        next.push(updatedCard)
      }

      const affectedLists = new Set([activeCardItem.list, targetList])
      for (const listId of affectedLists) {
        let order = 0
        for (let i = 0; i < next.length; i++) {
          if (next[i].list === listId && !next[i].archived) {
            next[i] = { ...next[i], order: order++ }
          }
        }
      }

      dragCardsRef.current = next
      return next
    })
  }

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event
    const dragType = active.data.current?.type

    if (dragType === 'list') {
      setActiveList(null)
      requestAnimationFrame(() => {
        justDraggedRef.current = false
      })

      if (!over || active.id === over.id) return

      const oldIndex = board.lists.findIndex((l) => l.id === active.id)
      const newIndex = board.lists.findIndex((l) => l.id === over.id)
      if (oldIndex === newIndex) return

      const reordered = arrayMove(board.lists, oldIndex, newIndex)
      setDragLists(reordered)

      try {
        await api.boards.update(board.id, { lists: reordered })
        onBoardUpdate()
      } catch (err) {
        console.error('Failed to reorder lists:', err)
      } finally {
        setDragLists(null)
      }
      return
    }

    setActiveCard(null)
    lastMoveRef.current = null
    requestAnimationFrame(() => {
      justDraggedRef.current = false
    })

    if (!over) {
      setDragCards(null)
      return
    }

    const activeId = active.id as string
    const currentDragCards = dragCardsRef.current
    const originalCard = cards.find((c) => c.id === activeId)
    const draggedCard = currentDragCards?.find((c) => c.id === activeId)

    if (!originalCard || !draggedCard) {
      setDragCards(null)
      return
    }

    const targetList = draggedCard.list
    const targetCards = (currentDragCards ?? [])
      .filter((c) => c.list === targetList && !c.archived)
      .sort((a, b) => a.order - b.order)
    const order = targetCards.findIndex((c) => c.id === activeId)

    if (originalCard.list === targetList && originalCard.order === order) {
      setDragCards(null)
      return
    }

    try {
      await api.cards.move(
        board.id,
        activeId,
        targetList,
        order >= 0 ? order : 0,
      )
      onRefresh()
    } catch (err) {
      console.error('Failed to move card:', err)
    } finally {
      dragCardsRef.current = null
      setDragCards(null)
    }
  }

  const handleDragCancel = () => {
    setActiveCard(null)
    setActiveList(null)
    dragCardsRef.current = null
    lastMoveRef.current = null
    setDragCards(null)
    setDragLists(null)
    requestAnimationFrame(() => {
      justDraggedRef.current = false
    })
  }

  const handleAddCard = async (listId: string, title: string) => {
    try {
      await api.cards.create(board.id, { title, list: listId })
      onRefresh()
    } catch (err) {
      console.error('Failed to create card:', err)
    }
  }

  const handleSaveCard = async (updates: Partial<CardType>) => {
    if (!selectedCard) return
    try {
      await api.cards.update(board.id, selectedCard.id, updates)
      setSelectedCard(null)
      onRefresh()
    } catch (err) {
      console.error('Failed to update card:', err)
    }
  }

  const handleArchiveCard = async () => {
    if (!selectedCard) return
    try {
      await api.cards.archive(board.id, selectedCard.id, !selectedCard.archived)
      setSelectedCard(null)
      onRefresh()
    } catch (err) {
      console.error('Failed to archive card:', err)
    }
  }

  const handleDeleteCard = async () => {
    if (!selectedCard) return
    try {
      await api.cards.delete(board.id, selectedCard.id)
      setSelectedCard(null)
      onRefresh()
    } catch (err) {
      console.error('Failed to delete card:', err)
    }
  }

  const handleAddList = async (name: string) => {
    try {
      const existingIds = new Set(board.lists.map((l) => l.id))
      const newListId = generateUniqueListId(name, existingIds)
      const updatedLists = [...board.lists, { id: newListId, name }]
      await api.boards.update(board.id, { lists: updatedLists })
      onBoardUpdate()
    } catch (err) {
      console.error('Failed to add list:', err)
    }
  }

  const handleRenameList = async (listId: string, newName: string) => {
    try {
      const updatedLists = board.lists.map((l) =>
        l.id === listId ? { ...l, name: newName } : l,
      )
      await api.boards.update(board.id, { lists: updatedLists })
      onBoardUpdate()
    } catch (err) {
      console.error('Failed to rename list:', err)
    }
  }

  const handleDeleteList = async (listId: string) => {
    const listCards = getCardsForList(listId)
    const message =
      listCards.length > 0
        ? `This list contains ${listCards.length} card(s). Are you sure you want to delete it?`
        : 'Are you sure you want to delete this list?'

    if (!window.confirm(message)) return

    try {
      const updatedLists = board.lists.filter((l) => l.id !== listId)
      await api.boards.update(board.id, { lists: updatedLists })
      onBoardUpdate()
    } catch (err) {
      console.error('Failed to delete list:', err)
    }
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={collisionDetection}
      onDragStart={handleDragStart}
      onDragMove={handleDragMove}
      onDragEnd={handleDragEnd}
      onDragCancel={handleDragCancel}
    >
      <div className={styles.board}>
        <SortableContext
          items={effectiveLists.map((l) => l.id)}
          strategy={horizontalListSortingStrategy}
        >
          {effectiveLists.map((list) => (
            <List
              key={list.id}
              list={list}
              cards={getCardsForList(list.id)}
              onCardClick={handleCardClick}
              onAddCard={(title) => handleAddCard(list.id, title)}
              onRename={(newName) => handleRenameList(list.id, newName)}
              onDelete={() => handleDeleteList(list.id)}
            />
          ))}
        </SortableContext>
        <AddList onAdd={handleAddList} />
      </div>

      <DragOverlay dropAnimation={null}>
        {activeCard ? (
          <div
            className={`${cardStyles.card} ${cardStyles.overlay} ${styles.dragOverlay}`}
          >
            <div className={cardStyles.title}>{activeCard.title}</div>
            {(activeCard.labels ?? []).length > 0 && (
              <div className={cardStyles.labels}>
                {(activeCard.labels ?? []).map((label) => (
                  <span key={label} className={cardStyles.label}>
                    {label}
                  </span>
                ))}
              </div>
            )}
            {(activeCard.todos ?? []).length > 0 &&
              (() => {
                const todos = activeCard.todos ?? []
                const completed = todos.filter((t) => t.completed).length
                const total = todos.length
                const percent = (completed / total) * 100
                return (
                  <div className={cardStyles.todoProgress}>
                    <div className={cardStyles.progressBar}>
                      <div
                        className={cardStyles.progressFill}
                        style={{ width: `${percent}%` }}
                      />
                    </div>
                    <span className={cardStyles.progressText}>
                      {completed}/{total}
                    </span>
                  </div>
                )
              })()}
          </div>
        ) : activeList ? (
          <div className={styles.listOverlay}>
            <div className={styles.listOverlayHeader}>{activeList.name}</div>
            <div className={styles.listOverlayBody}>
              {getCardsForList(activeList.id).map((card) => (
                <div key={card.id} className={styles.listOverlayCard}>
                  {card.title}
                </div>
              ))}
            </div>
          </div>
        ) : null}
      </DragOverlay>

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
