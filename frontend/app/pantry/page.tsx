"use client"

import { useEffect, useState } from "react"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface PantryItem {
  id: number
  name: string
  quantity: number
  unit: string
}

export default function PantryPage() {
  const [items, setItems] = useState<PantryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [showAddForm, setShowAddForm] = useState(false)
  const [newItem, setNewItem] = useState({ name: "", quantity: 0, unit: "" })

  useEffect(() => {
    fetchItems()
  }, [])

  const fetchItems = () => {
    api.get("/pantry")
      .then((response) => {
        // Убеждаемся, что response.data это массив
        setItems(Array.isArray(response.data) ? response.data : [])
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching pantry:", error)
        setItems([]) // Устанавливаем пустой массив при ошибке
        setLoading(false)
      })
  }

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.post("/pantry", newItem)
      setNewItem({ name: "", quantity: 0, unit: "" })
      setShowAddForm(false)
      fetchItems()
    } catch (error) {
      console.error("Error adding item:", error)
      alert("Не удалось добавить продукт")
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm("Удалить этот продукт?")) return
    
    try {
      await api.delete(`/pantry/${id}`)
      fetchItems()
    } catch (error) {
      console.error("Error deleting item:", error)
      alert("Не удалось удалить продукт")
    }
  }

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold">Кладовая</h1>
        <Button onClick={() => setShowAddForm(!showAddForm)} className="w-full sm:w-auto">
          {showAddForm ? "Отмена" : "Добавить продукт"}
        </Button>
      </div>

      {showAddForm && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Добавить продукт</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleAdd} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Название</label>
                <input
                  type="text"
                  value={newItem.name}
                  onChange={(e) => setNewItem({ ...newItem, name: e.target.value })}
                  className="w-full px-3 py-2 border rounded-md"
                  required
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Количество</label>
                  <input
                    type="number"
                    step="0.01"
                    value={newItem.quantity}
                    onChange={(e) => setNewItem({ ...newItem, quantity: parseFloat(e.target.value) })}
                    className="w-full px-3 py-2 border rounded-md"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Единица измерения</label>
                  <input
                    type="text"
                    value={newItem.unit}
                    onChange={(e) => setNewItem({ ...newItem, unit: e.target.value })}
                    className="w-full px-3 py-2 border rounded-md"
                    required
                  />
                </div>
              </div>
              <Button type="submit">Добавить</Button>
            </form>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {items.map((item) => (
          <Card key={item.id}>
            <CardHeader>
              <CardTitle>{item.name}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="mb-4">
                {item.quantity} {item.unit}
              </p>
              <Button
                variant="destructive"
                onClick={() => handleDelete(item.id)}
              >
                Удалить
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>

      {items.length === 0 && (
        <p className="text-center text-muted-foreground mt-8">
          В кладовой нет продуктов. Добавьте ингредиенты!
        </p>
      )}
    </div>
  )
}


