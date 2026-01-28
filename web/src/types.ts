export interface List {
  id: string
  name: string
}

export interface Board {
  id: string
  name: string
  lists: List[]
}

export interface TodoItem {
  id: string
  text: string
  completed: boolean
}

export interface Card {
  id: string
  title: string
  list: string
  order: number
  description: string
  labels: string[]
  todos: TodoItem[]
  archived: boolean
  created_at: string
  updated_at: string
}

export interface APIError {
  error: {
    code: string
    message: string
  }
}

export interface WSEvent {
  type: 'board_updated' | 'card_updated'
  board_id: string
  timestamp: string
}
