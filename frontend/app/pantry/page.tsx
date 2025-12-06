"use client"

import { useEffect, useState } from "react"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Trash2, Plus, Package, Search, X } from "lucide-react"

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
  const [searchQuery, setSearchQuery] = useState("")

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
    return (
      <div className="min-h-screen gradient-bg flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-primary border-r-transparent"></div>
          <p className="mt-4 text-muted-foreground">Загрузка кладовой...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen gradient-bg">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
          <div>
            <h1 className="text-3xl sm:text-4xl font-bold mb-2 bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
              Кладовая
            </h1>
            <p className="text-muted-foreground">Управляйте вашими ингредиентами</p>
          </div>
          <Button 
            onClick={() => setShowAddForm(!showAddForm)} 
            className="w-full sm:w-auto"
            size="sm"
          >
            {showAddForm ? (
              <>Отмена</>
            ) : (
              <>
                <Plus className="h-4 w-4 mr-2" />
                Добавить продукт
              </>
            )}
          </Button>
        </div>

      {showAddForm && (
        <Card className="mb-6">
          <CardHeader className="pb-4">
            <CardTitle className="text-lg">Добавить продукт</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleAdd} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Название</label>
                <input
                  type="text"
                  value={newItem.name}
                  onChange={(e) => setNewItem({ ...newItem, name: e.target.value })}
                  className="w-full px-3 py-2 border rounded-lg bg-background"
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
                    className="w-full px-3 py-2 border rounded-lg bg-background"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Единица измерения</label>
                  <input
                    type="text"
                    value={newItem.unit}
                    onChange={(e) => setNewItem({ ...newItem, unit: e.target.value })}
                    className="w-full px-3 py-2 border rounded-lg bg-background"
                    required
                  />
                </div>
              </div>
              <Button type="submit" size="sm">Добавить</Button>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Поиск */}
      {items.length > 0 && (
        <div className="mb-6">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <input
              type="text"
              placeholder="Поиск продуктов..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-10 py-2 border rounded-lg bg-background focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery("")}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
              >
                <X className="h-4 w-4" />
              </button>
            )}
          </div>
        </div>
      )}

      {(() => {
        const filteredItems = items.filter((item) =>
          item.name.toLowerCase().includes(searchQuery.toLowerCase())
        )

        if (filteredItems.length === 0 && items.length > 0) {
          return (
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-8">
                  <Search className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">
                    Ничего не найдено по запросу "{searchQuery}"
                  </p>
                </div>
              </CardContent>
            </Card>
          )
        }

        if (items.length === 0) {
          return (
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-8">
                  <Package className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                  <p className="text-muted-foreground">
                    В кладовой нет продуктов. Добавьте ингредиенты!
                  </p>
                </div>
              </CardContent>
            </Card>
          )
        }

        return (
          <>
            {searchQuery && (
              <p className="text-sm text-muted-foreground mb-4">
                Найдено: {filteredItems.length} {filteredItems.length === 1 ? 'продукт' : filteredItems.length < 5 ? 'продукта' : 'продуктов'}
              </p>
            )}
            <Card>
              <CardContent className="p-0">
                <div className="divide-y">
                  {filteredItems.map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between p-4 hover:bg-muted/50 transition-colors"
                >
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-base truncate">{item.name}</h3>
                    <p className="text-sm text-muted-foreground mt-1">
                      {item.quantity} {item.unit}
                    </p>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(item.id)}
                    className="ml-4 text-destructive hover:text-destructive hover:bg-destructive/10"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </>
        )
      })()}
      </div>
    </div>
  )
}


