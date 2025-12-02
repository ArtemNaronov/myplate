"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface Menu {
  id: number
  date: string
  total_calories: number
  total_time: number
}

export default function MenusPage() {
  const [menus, setMenus] = useState<Menu[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Получаем список всех меню пользователя
    api.get("/menus")
      .then((response) => {
        setMenus(response.data || [])
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching menus:", error)
        setLoading(false)
      })
  }, [])

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold">Мои меню</h1>
        <Link href="/menu/generate" className="w-full sm:w-auto">
          <Button className="w-full sm:w-auto">Создать новое меню</Button>
        </Link>
      </div>

      {menus.length === 0 ? (
        <Card>
          <CardContent className="pt-6">
            <p className="text-muted-foreground text-center">
              У вас пока нет созданных меню. <Link href="/menu/generate" className="text-primary underline">Создайте меню</Link>, чтобы начать.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4">
          {menus.map((menu) => (
            <Card key={menu.id}>
              <CardHeader>
                <CardTitle>Меню на {new Date(menu.date).toLocaleDateString('ru-RU')}</CardTitle>
                <CardDescription>
                  {menu.total_calories} ккал • {menu.total_time} мин
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex flex-col sm:flex-row gap-2">
                  <Link href={`/menu/${menu.id}`} className="flex-1">
                    <Button className="w-full">Просмотреть меню</Button>
                  </Link>
                  <Link href={`/shopping-list/${menu.id}`} className="flex-1">
                    <Button className="w-full" variant="outline">Список покупок</Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

