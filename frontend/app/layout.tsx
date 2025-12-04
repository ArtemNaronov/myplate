import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { TelegramProvider } from "@/components/telegram-provider"
import { ThemeProvider } from "@/components/theme-provider"
import { Nav } from "@/components/nav"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "MyPlateService - Помощник по рецептам и меню",
  description: "Генерация ежедневного меню на основе калорий, времени приготовления и наличия ингредиентов",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ru" suppressHydrationWarning>
      <body className={inter.className}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <TelegramProvider>
            <Nav />
            {children}
          </TelegramProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}

