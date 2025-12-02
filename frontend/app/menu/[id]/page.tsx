"use client"

import { useEffect, useState } from "react"
import { useParams, useRouter } from "next/navigation"
import Link from "next/link"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface MenuMeal {
  recipe_id: number
  meal_type: string
  calories: number
  time: number
}

interface Recipe {
  id: number
  name: string
  description: string
}

interface Menu {
  id: number
  date: string
  total_calories: number
  total_time: number
  meals: MenuMeal[]
}

export default function MenuDetailPage() {
  const params = useParams()
  const router = useRouter()
  const [menu, setMenu] = useState<Menu | null>(null)
  const [recipes, setRecipes] = useState<Record<number, Recipe>>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Получаем информацию о меню
    api.get(`/menus/${params.id}`)
      .then((response) => {
        setMenu(response.data)
        
        // Получаем информацию о рецептах
        const recipeIds = response.data.meals.map((meal: MenuMeal) => meal.recipe_id)
        const recipePromises = recipeIds.map((id: number) => 
          api.get(`/recipes/${id}`).then(res => res.data)
        )
        
        Promise.all(recipePromises).then((recipeData) => {
          const recipesMap: Record<number, Recipe> = {}
          recipeData.forEach((recipe: Recipe) => {
            recipesMap[recipe.id] = recipe
          })
          setRecipes(recipesMap)
          setLoading(false)
        })
      })
      .catch((error) => {
        console.error("Error fetching menu:", error)
        setLoading(false)
      })
  }, [params.id])

  const getMealTypeLabel = (type: string) => {
    switch (type) {
      case 'breakfast': return 'Завтрак'
      case 'lunch': return 'Обед'
      case 'dinner': return 'Ужин'
      default: return type
    }
  }

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  if (!menu) {
    return <div className="container mx-auto px-4 py-8">Меню не найдено</div>
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold mb-2">Меню на {new Date(menu.date).toLocaleDateString('ru-RU')}</h1>
          <p className="text-muted-foreground">Ваше ежедневное меню</p>
        </div>
        <Link href={`/shopping-list/${menu.id}`} className="w-full sm:w-auto">
          <Button className="w-full sm:w-auto">Список покупок</Button>
        </Link>
      </div>

      <div className="grid gap-6 md:grid-cols-2 mb-8">
        <Card>
          <CardHeader>
            <CardTitle>Калории</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold">{menu.total_calories}</p>
            <p className="text-sm text-muted-foreground">ккал</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Время</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-bold">{menu.total_time}</p>
            <p className="text-sm text-muted-foreground">минут</p>
          </CardContent>
        </Card>
      </div>

      <div className="space-y-6">
        <h2 className="text-2xl font-bold">Блюда</h2>
        
        {menu.meals.map((meal, idx) => {
          const recipe = recipes[meal.recipe_id]
          return (
            <Card key={idx}>
              <CardHeader>
                <div className="flex flex-col sm:flex-row justify-between items-start gap-3">
                  <div className="flex-1">
                    <CardTitle>{getMealTypeLabel(meal.meal_type)}</CardTitle>
                    {recipe && (
                      <CardDescription>{recipe.name}</CardDescription>
                    )}
                  </div>
                  <Link href={`/recipes/${meal.recipe_id}`} className="w-full sm:w-auto">
                    <Button variant="outline" size="sm" className="w-full sm:w-auto">Рецепт</Button>
                  </Link>
                </div>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">Калории</p>
                    <p className="font-semibold">{meal.calories} ккал</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Время</p>
                    <p className="font-semibold">{meal.time} мин</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>

      <div className="mt-8">
        <Link href="/menu/generate">
          <Button variant="outline">Создать новое меню</Button>
        </Link>
      </div>
    </div>
  )
}

