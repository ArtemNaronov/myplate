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
      // Если токена нет, авторизуемся через тестовый endpoint
      const token = localStorage.getItem("token")
      if (!token) {
        api.get("/auth/test")
          .then((response) => {
            if (response.data && response.data.token) {
              localStorage.setItem("token", response.data.token)
              setUser(response.data.user)
              console.log("Авторизация успешна (тестовый пользователь)")
            }
          })
          .catch((error) => {
            console.error("Ошибка тестовой авторизации:", error)
            // Показываем предупреждение пользователю
            if (error.response) {
              console.error("Ответ сервера:", error.response.data)
            }
          })
      } else {
        // Если токен есть, проверяем его валидность, получая информацию о пользователе
        // (можно добавить endpoint для проверки токена)
        console.log("Токен найден в localStorage")
      }
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


