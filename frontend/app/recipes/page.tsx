"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface Recipe {
  id: number
  name: string
  description: string
  calories: number
  cooking_time: number
  meal_type: string
}

export default function RecipesPage() {
  const [recipes, setRecipes] = useState<Recipe[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get("/recipes")
      .then((response) => {
        setRecipes(response.data)
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching recipes:", error)
        setLoading(false)
      })
  }, [])

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Рецепты</h1>
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {recipes.map((recipe) => (
          <Card key={recipe.id}>
            <CardHeader>
              <CardTitle>{recipe.name}</CardTitle>
              <CardDescription>{recipe.description}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2 mb-4">
                <p className="text-sm">
                  <span className="font-semibold">Калории:</span> {recipe.calories}
                </p>
                <p className="text-sm">
                  <span className="font-semibold">Время:</span> {recipe.cooking_time} мин
                </p>
                <p className="text-sm">
                  <span className="font-semibold">Приём пищи:</span> {recipe.meal_type === 'breakfast' ? 'Завтрак' : recipe.meal_type === 'lunch' ? 'Обед' : recipe.meal_type === 'dinner' ? 'Ужин' : recipe.meal_type}
                </p>
              </div>
              <Link href={`/recipes/${recipe.id}`}>
                <Button className="w-full">Подробнее</Button>
              </Link>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}


