"use client"

import { useEffect, useState } from "react"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface Goals {
  daily_calories: number
  protein_ratio: number
  fat_ratio: number
  carb_ratio: number
}

export default function GoalsPage() {
  const [goals, setGoals] = useState<Goals | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    api.get("/users/goals")
      .then((response) => {
        setGoals(response.data)
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching goals:", error)
        setLoading(false)
      })
  }, [])

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setSaving(true)

    try {
      await api.post("/users/goals", goals)
      alert("Цели успешно сохранены!")
    } catch (error) {
      console.error("Error saving goals:", error)
      alert("Не удалось сохранить цели")
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  const defaultGoals: Goals = {
    daily_calories: 2000,
    protein_ratio: 30,
    fat_ratio: 30,
    carb_ratio: 40,
  }

  const currentGoals = goals || defaultGoals

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-2xl sm:text-3xl font-bold mb-6">Мои цели</h1>

      <Card>
        <CardHeader>
          <CardTitle>Цели по питанию</CardTitle>
          <CardDescription>Установите ваши ежедневные цели по калориям и макроэлементам</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label className="block text-sm font-medium mb-2">
                Ежедневные калории
              </label>
              <input
                type="number"
                value={currentGoals.daily_calories}
                onChange={(e) => setGoals({ ...currentGoals, daily_calories: parseInt(e.target.value) })}
                className="w-full px-3 py-2 border rounded-md"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Доля белков (%)
              </label>
              <input
                type="number"
                step="0.1"
                value={currentGoals.protein_ratio}
                onChange={(e) => setGoals({ ...currentGoals, protein_ratio: parseFloat(e.target.value) })}
                className="w-full px-3 py-2 border rounded-md"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Доля жиров (%)
              </label>
              <input
                type="number"
                step="0.1"
                value={currentGoals.fat_ratio}
                onChange={(e) => setGoals({ ...currentGoals, fat_ratio: parseFloat(e.target.value) })}
                className="w-full px-3 py-2 border rounded-md"
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">
                Доля углеводов (%)
              </label>
              <input
                type="number"
                step="0.1"
                value={currentGoals.carb_ratio}
                onChange={(e) => setGoals({ ...currentGoals, carb_ratio: parseFloat(e.target.value) })}
                className="w-full px-3 py-2 border rounded-md"
                required
              />
            </div>

            <Button type="submit" className="w-full" disabled={saving}>
              {saving ? "Сохранение..." : "Сохранить цели"}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}


