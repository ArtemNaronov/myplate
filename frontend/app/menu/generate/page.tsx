"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

export default function MenuGeneratePage() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    user_id: 1, // In production, get from auth context
    target_calories: 2000,
    diet_type: "",
    allergies: [] as string[],
    max_total_time: 0,
    consider_pantry: false,
    pantry_importance: "prefer",
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      // Подготавливаем данные для отправки - убираем нулевые значения
      const requestData = {
        ...formData,
        max_total_time: formData.max_total_time > 0 ? formData.max_total_time : undefined,
        diet_type: formData.diet_type || undefined,
      }
      
      const response = await api.post("/menus/generate", requestData)
      // Перенаправляем на страницу с детальной информацией о меню
      router.push(`/menu/${response.data.id}`)
    } catch (error: any) {
      console.error("Error generating menu:", error)
      const errorMessage = error.response?.data?.error || "Не удалось создать меню"
      alert(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-2xl sm:text-3xl font-bold mb-6">Создать ежедневное меню</h1>

      <Card>
        <CardHeader>
          <CardTitle>Настройки меню</CardTitle>
          <CardDescription>Настройте генерацию вашего ежедневного меню</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium mb-2">
                Целевые калории
              </label>
              <input
                type="number"
                value={formData.target_calories}
                onChange={(e) => setFormData({ ...formData, target_calories: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border rounded-md"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Максимальное время приготовления (минуты)
              </label>
              <input
                type="number"
                value={formData.max_total_time}
                onChange={(e) => setFormData({ ...formData, max_total_time: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Тип диеты
              </label>
              <select
                value={formData.diet_type}
                onChange={(e) => setFormData({ ...formData, diet_type: e.target.value })}
                className="w-full px-3 py-2 border rounded-md"
              >
                <option value="">Нет</option>
                <option value="vegetarian">Вегетарианская</option>
                <option value="vegan">Веганская</option>
                <option value="gluten-free">Без глютена</option>
              </select>
            </div>

            <div>
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  checked={formData.consider_pantry}
                  onChange={(e) => setFormData({ ...formData, consider_pantry: e.target.checked })}
                />
                <span>Учитывать кладовую</span>
              </label>
            </div>

            {formData.consider_pantry && (
              <div>
                <label className="block text-sm font-medium mb-2">
                  Важность кладовой
                </label>
                <select
                  value={formData.pantry_importance}
                  onChange={(e) => setFormData({ ...formData, pantry_importance: e.target.value })}
                  className="w-full px-3 py-2 border rounded-md"
                >
                  <option value="ignore">Игнорировать</option>
                  <option value="prefer">Предпочитать</option>
                  <option value="strict">Строго</option>
                </select>
              </div>
            )}

            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "Создание..." : "Создать меню"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}


