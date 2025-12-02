"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface Ingredient {
  name: string
  quantity: number
  unit: string
}

interface Recipe {
  id: number
  name: string
  description: string
  calories: number
  proteins: number
  fats: number
  carbs: number
  cooking_time: number
  meal_type: string
  ingredients: Ingredient[] | null
  instructions: string[] | null
}

export default function RecipeDetailPage() {
  const params = useParams()
  const [recipe, setRecipe] = useState<Recipe | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get(`/recipes/${params.id}`)
      .then((response) => {
        setRecipe(response.data)
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching recipe:", error)
        setLoading(false)
      })
  }, [params.id])

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  if (!recipe) {
    return <div className="container mx-auto px-4 py-8">Рецепт не найден</div>
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-2xl sm:text-3xl font-bold mb-4">{recipe.name}</h1>
      <p className="text-muted-foreground mb-6">{recipe.description}</p>

      <div className="grid gap-6 md:grid-cols-2 mb-8">
        <Card>
          <CardHeader>
            <CardTitle>Питание</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <p><span className="font-semibold">Калории:</span> {recipe.calories}</p>
            <p><span className="font-semibold">Белки:</span> {recipe.proteins}г</p>
            <p><span className="font-semibold">Жиры:</span> {recipe.fats}г</p>
            <p><span className="font-semibold">Углеводы:</span> {recipe.carbs}г</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Детали</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <p><span className="font-semibold">Время приготовления:</span> {recipe.cooking_time} мин</p>
            <p><span className="font-semibold">Тип блюда:</span> {recipe.meal_type === 'breakfast' ? 'Завтрак' : recipe.meal_type === 'lunch' ? 'Обед' : recipe.meal_type === 'dinner' ? 'Ужин' : recipe.meal_type}</p>
          </CardContent>
        </Card>
      </div>

      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Ингредиенты</CardTitle>
        </CardHeader>
        <CardContent>
          {recipe.ingredients && recipe.ingredients.length > 0 ? (
            <ul className="list-disc list-inside space-y-2">
              {recipe.ingredients.map((ing, idx) => (
                <li key={idx}>
                  {ing.quantity} {ing.unit} {ing.name}
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-muted-foreground">Ингредиенты не указаны</p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Инструкции</CardTitle>
        </CardHeader>
        <CardContent>
          {recipe.instructions && recipe.instructions.length > 0 ? (
            <ol className="list-decimal list-inside space-y-2">
              {recipe.instructions.map((instruction, idx) => (
                <li key={idx}>{instruction}</li>
              ))}
            </ol>
          ) : (
            <p className="text-muted-foreground">Инструкции не указаны</p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}


