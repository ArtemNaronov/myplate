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
    return (
      <div className="min-h-screen gradient-bg flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-primary border-r-transparent"></div>
          <p className="mt-4 text-muted-foreground">Загрузка рецептов...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen gradient-bg">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2 bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
            Рецепты
          </h1>
          <p className="text-muted-foreground">Откройте для себя нашу коллекцию вкусных рецептов</p>
        </div>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {recipes.map((recipe) => (
            <Card key={recipe.id} className="flex flex-col card-hover group">
              <CardHeader>
                <CardTitle className="group-hover:text-primary transition-colors duration-300">
                  {recipe.name}
                </CardTitle>
                <CardDescription className="line-clamp-2">{recipe.description}</CardDescription>
              </CardHeader>
              <CardContent className="flex-1 flex flex-col">
                <div className="space-y-2 mb-4 flex-1">
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-semibold text-primary">Калории:</span>
                    <span>{recipe.calories} ккал</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-semibold text-primary">Время:</span>
                    <span>{recipe.cooking_time} мин</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm">
                    <span className="font-semibold text-primary">Приём пищи:</span>
                    <span>{recipe.meal_type === 'breakfast' ? 'Завтрак' : recipe.meal_type === 'lunch' ? 'Обед' : recipe.meal_type === 'dinner' ? 'Ужин' : recipe.meal_type}</span>
                  </div>
                </div>
                <Link href={`/recipes/${recipe.id}`}>
                  <Button className="w-full group-hover:shadow-lg transition-all duration-300">
                    Подробнее
                  </Button>
                </Link>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  )
}


