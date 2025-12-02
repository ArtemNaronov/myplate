"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface ShoppingItem {
  name: string
  quantity: number
  unit: string
  reason: string[]
}

export default function ShoppingListPage() {
  const params = useParams()
  const [items, setItems] = useState<ShoppingItem[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get(`/shopping-list/${params.menu_id}`)
      .then((response) => {
        setItems(response.data.items || [])
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching shopping list:", error)
        setLoading(false)
      })
  }, [params.menu_id])

  const handleCopy = () => {
    const text = items.map(item => 
      `${item.name}: ${item.quantity} ${item.unit}`
    ).join("\n")
    navigator.clipboard.writeText(text)
    alert("Список покупок скопирован в буфер обмена!")
  }

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold">Список покупок</h1>
        <Button onClick={handleCopy} className="w-full sm:w-auto">
          Копировать список
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Что купить</CardTitle>
          <CardDescription>Продукты, необходимые для вашего меню</CardDescription>
        </CardHeader>
        <CardContent>
          {items.length === 0 ? (
            <p className="text-muted-foreground">Ничего не нужно покупать. У вас всё есть в кладовой!</p>
          ) : (
            <ul className="space-y-4">
              {items.map((item, idx) => (
                <li key={idx} className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 border-b pb-2">
                  <div className="flex-1">
                    <p className="font-semibold">{item.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {item.quantity} {item.unit} • Нужно для: {item.reason && Array.isArray(item.reason) && item.reason.length > 0 ? item.reason.map(r => r === 'breakfast' ? 'Завтрак' : r === 'lunch' ? 'Обед' : r === 'dinner' ? 'Ужин' : r).join(", ") : 'Не указано'}
                    </p>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  )
}


