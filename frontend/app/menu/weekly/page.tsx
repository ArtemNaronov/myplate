"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import api from "@/lib/api"
import { useTelegram } from "@/components/telegram-provider"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import MacrosChart from "@/components/macros-chart"

interface RecipeDTO {
  id: number
  name: string
  description?: string
  calories: number
  proteins: number
  fats: number
  carbs: number
  cooking_time: number
  servings: number
  meal_type: string
  ingredients: Array<{ name: string; quantity: number; unit: string }>
  instructions?: string[]
}

interface WeeklyDayMenu {
  day: number
  breakfast: RecipeDTO
  lunch: RecipeDTO
  dinner: RecipeDTO
  totalCalories: number
  totalProteins: number
  totalFats: number
  totalCarbs: number
  totalTime?: number
  ingredients_used?: Array<{ name: string; quantity: number; unit: string }>
  missing_ingredients?: Array<{ name: string; quantity: number; unit: string }>
}

interface WeeklyMenu {
  week: WeeklyDayMenu[]
}

export default function WeeklyMenuPage() {
  const router = useRouter()
  const { user } = useTelegram()
  const [loading, setLoading] = useState(false)
  const [weeklyMenu, setWeeklyMenu] = useState<WeeklyMenu | null>(null)
  const [viewMode, setViewMode] = useState(false) // Режим просмотра сохраненного меню
  const [formData, setFormData] = useState({
    adults: 1,
    children: 0,
    diet_type: "",
    allergies: "",
    max_total_time: 0,
    max_time_per_meal: 0,
    consider_pantry: false,
    pantry_importance: "prefer",
  })

  // Проверяем, есть ли сохраненное меню для просмотра
  useEffect(() => {
    const savedMenu = localStorage.getItem("viewWeeklyMenu")
    if (savedMenu) {
      try {
        const menu = JSON.parse(savedMenu)
        setWeeklyMenu(menu)
        setViewMode(true)
        // Удаляем из localStorage после загрузки
        localStorage.removeItem("viewWeeklyMenu")
      } catch (e) {
        console.error("Error parsing saved weekly menu:", e)
      }
    }
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!user?.id) {
      alert("Пожалуйста, войдите в систему для создания меню")
      router.push("/auth/login")
      return
    }

    if (formData.adults < 1) {
      alert("Количество взрослых должно быть не менее 1")
      return
    }

    setLoading(true)

    try {
      // Формируем query параметры
      const params = new URLSearchParams()
      params.append("adults", formData.adults.toString())
      if (formData.children > 0) {
        params.append("children", formData.children.toString())
      }
      if (formData.diet_type) {
        params.append("diet_type", formData.diet_type)
      }
      if (formData.allergies) {
        params.append("allergies", formData.allergies)
      }
      if (formData.max_total_time > 0) {
        params.append("max_total_time", formData.max_total_time.toString())
      }
      if (formData.max_time_per_meal > 0) {
        params.append("max_time_per_meal", formData.max_time_per_meal.toString())
      }
      if (formData.consider_pantry) {
        params.append("consider_pantry", "true")
        params.append("pantry_importance", formData.pantry_importance)
      }
      
      const response = await api.get(`/menu/weekly?${params.toString()}`)
      setWeeklyMenu(response.data)
      
      // Сохраняем недельное меню в базу данных
      try {
        await api.post("/menu/weekly/save", response.data)
        // Также сохраняем в localStorage для быстрого доступа
        localStorage.setItem("lastWeeklyMenu", JSON.stringify(response.data))
      } catch (saveError: any) {
        console.error("Ошибка при сохранении недельного меню:", saveError)
        // Не прерываем выполнение, просто логируем ошибку
      }
    } catch (error: any) {
      console.error("Error generating weekly menu:", error)
      const errorMessage = error.response?.data?.error || "Не удалось создать недельное меню"
      alert(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      <h1 className="text-2xl sm:text-3xl font-bold mb-6">
        {viewMode ? "Меню на неделю" : "Генерация меню на неделю"}
      </h1>

      {!viewMode && (
      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Настройки недельного меню</CardTitle>
          <CardDescription>Настройте генерацию меню на 7 дней</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-2">
                  Количество взрослых *
                </label>
                <input
                  type="number"
                  min="1"
                  value={formData.adults}
                  onChange={(e) => setFormData({ ...formData, adults: parseInt(e.target.value) || 1 })}
                  className="w-full px-3 py-2 border rounded-md"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">
                  Количество детей
                </label>
                <input
                  type="number"
                  min="0"
                  value={formData.children}
                  onChange={(e) => setFormData({ ...formData, children: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 border rounded-md"
                />
                <p className="text-xs text-muted-foreground mt-1">
                  Дети учитываются с коэффициентом 0.7
                </p>
              </div>
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
              <label className="block text-sm font-medium mb-2">
                Аллергены (через запятую)
              </label>
              <input
                type="text"
                value={formData.allergies}
                onChange={(e) => setFormData({ ...formData, allergies: e.target.value })}
                placeholder="nuts, dairy, eggs"
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-2">
                  Максимальное время приготовления (минуты)
                </label>
                <input
                  type="number"
                  min="0"
                  value={formData.max_total_time}
                  onChange={(e) => setFormData({ ...formData, max_total_time: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 border rounded-md"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-2">
                  Максимальное время на одно блюдо (минуты)
                </label>
                <input
                  type="number"
                  min="0"
                  value={formData.max_time_per_meal}
                  onChange={(e) => setFormData({ ...formData, max_time_per_meal: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 border rounded-md"
                />
              </div>
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
              {loading ? "Генерация..." : "Сгенерировать недельное меню"}
            </Button>
          </form>
        </CardContent>
      </Card>
      )}

      {viewMode && weeklyMenu && (
        <div className="mb-4">
          <Button 
            variant="outline" 
            onClick={() => {
              setViewMode(false)
              setWeeklyMenu(null)
            }}
          >
            Создать новое меню
          </Button>
        </div>
      )}

      {weeklyMenu && weeklyMenu.week && weeklyMenu.week.length > 0 && (
        <div className="space-y-6">
          <h2 className="text-xl font-bold">Меню на неделю</h2>
          
          {/* Общая статистика за неделю */}
          {(() => {
            const weekTotal = weeklyMenu.week.reduce((acc, day) => ({
              calories: acc.calories + (day.totalCalories || 0),
              proteins: acc.proteins + (day.totalProteins || 0),
              fats: acc.fats + (day.totalFats || 0),
              carbs: acc.carbs + (day.totalCarbs || 0),
              time: acc.time + (day.totalTime || 0)
            }), { calories: 0, proteins: 0, fats: 0, carbs: 0, time: 0 })
            
            // Калории из макронутриентов
            const proteinCalories = weekTotal.proteins * 4
            const fatCalories = weekTotal.fats * 9
            const carbCalories = weekTotal.carbs * 4
            const totalMacroCalories = proteinCalories + fatCalories + carbCalories
            
            // Проценты
            const proteinPercent = totalMacroCalories > 0 ? (proteinCalories / totalMacroCalories) * 100 : 0
            const fatPercent = totalMacroCalories > 0 ? (fatCalories / totalMacroCalories) * 100 : 0
            const carbPercent = totalMacroCalories > 0 ? (carbCalories / totalMacroCalories) * 100 : 0
            
            return (
              <Card className="mb-6">
                <CardHeader>
                  <CardTitle>Итого за неделю</CardTitle>
                  <CardDescription>
                    {weekTotal.calories} ккал • {weekTotal.proteins.toFixed(1)}г Б • {weekTotal.fats.toFixed(1)}г Ж • {weekTotal.carbs.toFixed(1)}г У
                    {weekTotal.time > 0 && ` • ${weekTotal.time} мин`}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
                    {/* Диаграмма слева */}
                    <div className="flex-shrink-0">
                      <MacrosChart 
                        proteins={weekTotal.proteins} 
                        fats={weekTotal.fats} 
                        carbs={weekTotal.carbs}
                        size={250}
                        showLabels={false}
                      />
                    </div>
                    
                    {/* Информация справа */}
                    <div className="flex-1 w-full md:w-auto">
                      <div className="space-y-4">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="w-5 h-5 rounded bg-blue-500"></div>
                            <span className="font-medium">Белки</span>
                          </div>
                          <div className="text-right">
                            <div className="font-semibold">{weekTotal.proteins.toFixed(1)}г</div>
                            <div className="text-sm text-muted-foreground">{proteinPercent.toFixed(1)}%</div>
                          </div>
                        </div>
                        
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="w-5 h-5 rounded bg-orange-500"></div>
                            <span className="font-medium">Жиры</span>
                          </div>
                          <div className="text-right">
                            <div className="font-semibold">{weekTotal.fats.toFixed(1)}г</div>
                            <div className="text-sm text-muted-foreground">{fatPercent.toFixed(1)}%</div>
                          </div>
                        </div>
                        
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <div className="w-5 h-5 rounded bg-green-500"></div>
                            <span className="font-medium">Углеводы</span>
                          </div>
                          <div className="text-right">
                            <div className="font-semibold">{weekTotal.carbs.toFixed(1)}г</div>
                            <div className="text-sm text-muted-foreground">{carbPercent.toFixed(1)}%</div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )
          })()}
          
          {weeklyMenu.week.map((dayMenu) => (
            <Card key={dayMenu.day}>
              <CardHeader>
                <CardTitle>День {dayMenu.day}</CardTitle>
                <CardDescription>
                  {dayMenu.totalCalories} ккал • {(dayMenu.totalProteins || 0).toFixed(1)}г Б • {(dayMenu.totalFats || 0).toFixed(1)}г Ж • {(dayMenu.totalCarbs || 0).toFixed(1)}г У
                  {dayMenu.totalTime && ` • ${dayMenu.totalTime} мин`}
                </CardDescription>
              </CardHeader>
              <CardContent>
                {(() => {
                  const dayProteins = dayMenu.totalProteins || 0
                  const dayFats = dayMenu.totalFats || 0
                  const dayCarbs = dayMenu.totalCarbs || 0
                  
                  // Калории из макронутриентов
                  const dayProteinCalories = dayProteins * 4
                  const dayFatCalories = dayFats * 9
                  const dayCarbCalories = dayCarbs * 4
                  const dayTotalMacroCalories = dayProteinCalories + dayFatCalories + dayCarbCalories
                  
                  // Проценты
                  const dayProteinPercent = dayTotalMacroCalories > 0 ? (dayProteinCalories / dayTotalMacroCalories) * 100 : 0
                  const dayFatPercent = dayTotalMacroCalories > 0 ? (dayFatCalories / dayTotalMacroCalories) * 100 : 0
                  const dayCarbPercent = dayTotalMacroCalories > 0 ? (dayCarbCalories / dayTotalMacroCalories) * 100 : 0
                  
                  return (
                    <div className="mb-6">
                      <div className="flex flex-col md:flex-row items-center md:items-start gap-6">
                        {/* Диаграмма слева */}
                        <div className="flex-shrink-0">
                          <MacrosChart 
                            proteins={dayProteins} 
                            fats={dayFats} 
                            carbs={dayCarbs}
                            size={200}
                            showLabels={false}
                          />
                        </div>
                        
                        {/* Информация справа */}
                        <div className="flex-1 w-full md:w-auto">
                          <div className="space-y-4">
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3">
                                <div className="w-5 h-5 rounded bg-blue-500"></div>
                                <span className="font-medium">Белки</span>
                              </div>
                              <div className="text-right">
                                <div className="font-semibold">{dayProteins.toFixed(1)}г</div>
                                <div className="text-sm text-muted-foreground">{dayProteinPercent.toFixed(1)}%</div>
                              </div>
                            </div>
                            
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3">
                                <div className="w-5 h-5 rounded bg-orange-500"></div>
                                <span className="font-medium">Жиры</span>
                              </div>
                              <div className="text-right">
                                <div className="font-semibold">{dayFats.toFixed(1)}г</div>
                                <div className="text-sm text-muted-foreground">{dayFatPercent.toFixed(1)}%</div>
                              </div>
                            </div>
                            
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-3">
                                <div className="w-5 h-5 rounded bg-green-500"></div>
                                <span className="font-medium">Углеводы</span>
                              </div>
                              <div className="text-right">
                                <div className="font-semibold">{dayCarbs.toFixed(1)}г</div>
                                <div className="text-sm text-muted-foreground">{dayCarbPercent.toFixed(1)}%</div>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  )
                })()}
                <div className="grid gap-4 md:grid-cols-3">
                  <div>
                    <h3 className="font-semibold mb-2">Завтрак</h3>
                    {dayMenu.breakfast && (
                      <Link href={`/recipes/${dayMenu.breakfast.id}`} className="block">
                        <div className="text-sm p-3 rounded-lg border hover:bg-accent transition-colors cursor-pointer">
                          <p className="font-medium">{dayMenu.breakfast.name}</p>
                          <p className="text-muted-foreground">
                            {dayMenu.breakfast.calories} ккал • {dayMenu.breakfast.cooking_time} мин
                          </p>
                        </div>
                      </Link>
                    )}
                  </div>
                  <div>
                    <h3 className="font-semibold mb-2">Обед</h3>
                    {dayMenu.lunch && (
                      <Link href={`/recipes/${dayMenu.lunch.id}`} className="block">
                        <div className="text-sm p-3 rounded-lg border hover:bg-accent transition-colors cursor-pointer">
                          <p className="font-medium">{dayMenu.lunch.name}</p>
                          <p className="text-muted-foreground">
                            {dayMenu.lunch.calories} ккал • {dayMenu.lunch.cooking_time} мин
                          </p>
                        </div>
                      </Link>
                    )}
                  </div>
                  <div>
                    <h3 className="font-semibold mb-2">Ужин</h3>
                    {dayMenu.dinner && (
                      <Link href={`/recipes/${dayMenu.dinner.id}`} className="block">
                        <div className="text-sm p-3 rounded-lg border hover:bg-accent transition-colors cursor-pointer">
                          <p className="font-medium">{dayMenu.dinner.name}</p>
                          <p className="text-muted-foreground">
                            {dayMenu.dinner.calories} ккал • {dayMenu.dinner.cooking_time} мин
                          </p>
                        </div>
                      </Link>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}

