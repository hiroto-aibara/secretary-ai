import type { Board, Card } from '../types'

const BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error?.message || 'unknown error')
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export const api = {
  boards: {
    list: () => request<Board[]>('/boards'),
    get: (id: string) => request<Board>(`/boards/${id}`),
    create: (board: Partial<Board>) =>
      request<Board>('/boards', {
        method: 'POST',
        body: JSON.stringify(board),
      }),
    update: (id: string, board: Partial<Board>) =>
      request<Board>(`/boards/${id}`, {
        method: 'PUT',
        body: JSON.stringify(board),
      }),
    delete: (id: string) =>
      request<void>(`/boards/${id}`, { method: 'DELETE' }),
  },
  cards: {
    list: (boardId: string, archived = false) =>
      request<Card[]>(
        `/boards/${boardId}/cards${archived ? '?archived=true' : ''}`,
      ),
    get: (boardId: string, cardId: string) =>
      request<Card>(`/boards/${boardId}/cards/${cardId}`),
    create: (boardId: string, card: Partial<Card>) =>
      request<Card>(`/boards/${boardId}/cards`, {
        method: 'POST',
        body: JSON.stringify(card),
      }),
    update: (boardId: string, cardId: string, card: Partial<Card>) =>
      request<Card>(`/boards/${boardId}/cards/${cardId}`, {
        method: 'PUT',
        body: JSON.stringify(card),
      }),
    delete: (boardId: string, cardId: string) =>
      request<void>(`/boards/${boardId}/cards/${cardId}`, {
        method: 'DELETE',
      }),
    move: (boardId: string, cardId: string, list: string, order: number) =>
      request<Card>(`/boards/${boardId}/cards/${cardId}/move`, {
        method: 'PATCH',
        body: JSON.stringify({ list, order }),
      }),
    archive: (boardId: string, cardId: string, archived: boolean) =>
      request<Card>(`/boards/${boardId}/cards/${cardId}/archive`, {
        method: 'PATCH',
        body: JSON.stringify({ archived }),
      }),
  },
} as const
