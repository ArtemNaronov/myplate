"use client"

import { createContext, useContext, useEffect, useState } from "react"
import api from "@/lib/api"

interface TelegramContextType {
  isTelegram: boolean
  initData: string | null
  user: any | null
}

const TelegramContext = createContext<TelegramContextType>({
  isTelegram: false,
  initData: null,
  user: null,
})

export function TelegramProvider({ children }: { children: React.ReactNode }) {
  const [isTelegram, setIsTelegram] = useState(false)
  const [initData, setInitData] = useState<string | null>(null)
  const [user, setUser] = useState<any | null>(null)

  useEffect(() => {
    // Check if running in Telegram WebApp
    if (typeof window !== "undefined" && (window as any).Telegram?.WebApp) {
      const tg = (window as any).Telegram.WebApp
      setIsTelegram(true)
      tg.ready()
      tg.expand()

      // Get initData
      const initDataStr = tg.initData
      if (initDataStr) {
        setInitData(initDataStr)
        
        // Authenticate with backend
        api.post("/auth/telegram", { init_data: initDataStr })
          .then((response) => {
            localStorage.setItem("token", response.data.token)
            setUser(response.data.user)
          })
          .catch((error) => {
            console.error("Telegram auth error:", error)
          })
      }
    } else {
      // Если не в Telegram, проверяем наличие токена
      const token = localStorage.getItem("token")
      if (token) {
        // Если токен есть, проверяем его валидность, получая информацию о пользователе
        api.get("/auth/profile")
          .then((response) => {
            setUser(response.data)
          })
          .catch((error) => {
            // Если токен невалидный или ошибка сервера, удаляем его
            if (error.response?.status === 401 || error.response?.status === 500) {
              localStorage.removeItem("token")
            }
            console.error("Ошибка при проверке токена:", error)
          })
      }
      // Если токена нет, пользователь должен войти через форму входа
    }
  }, [])

  return (
    <TelegramContext.Provider value={{ isTelegram, initData, user }}>
      {children}
    </TelegramContext.Provider>
  )
}

export function useTelegram() {
  return useContext(TelegramContext)
}


