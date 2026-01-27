import { useEffect, useRef } from 'react'
import type { WSEvent } from '../types'

export function useWebSocket(onEvent: (event: WSEvent) => void) {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimer = useRef<ReturnType<typeof setTimeout>>(undefined)
  const onEventRef = useRef(onEvent)
  const connectRef = useRef<(() => void) | undefined>(undefined)

  useEffect(() => {
    onEventRef.current = onEvent
  })

  useEffect(() => {
    let isActive = true

    const connect = () => {
      if (!isActive) return

      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const ws = new WebSocket(`${protocol}//${window.location.host}/ws`)

      ws.onmessage = (msg) => {
        try {
          const event: WSEvent = JSON.parse(msg.data)
          onEventRef.current(event)
        } catch {
          // ignore malformed messages
        }
      }

      ws.onclose = () => {
        wsRef.current = null
        if (isActive) {
          reconnectTimer.current = setTimeout(
            () => connectRef.current?.(),
            3000,
          )
        }
      }

      ws.onerror = () => {
        // Only close if connection was established
        if (ws.readyState === WebSocket.OPEN) {
          ws.close()
        }
      }

      wsRef.current = ws
    }

    connectRef.current = connect
    connect()

    return () => {
      isActive = false
      clearTimeout(reconnectTimer.current)
      const ws = wsRef.current
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close()
      }
      wsRef.current = null
    }
  }, [])
}
