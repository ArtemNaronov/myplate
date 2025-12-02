import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { TelegramProvider } from "@/components/telegram-provider"
import { Nav } from "@/components/nav"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "MyPlate - Помощник по рецептам и меню",
  description: "Генерация ежедневного меню на основе калорий, времени приготовления и наличия ингредиентов",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ru">
      <body className={inter.className}>
        <TelegramProvider>
          <Nav />
          {children}
        </TelegramProvider>
      </body>
    </html>
  )
}

