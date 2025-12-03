"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import api from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface Menu {
  id: number
  date: string
  total_calories: number
  total_time: number
}

interface WeeklyMenu {
  id?: number
  week: Array<{
    day: number
    breakfast: { id: number; name: string; calories: number; cooking_time: number }
    lunch: { id: number; name: string; calories: number; cooking_time: number }
    dinner: { id: number; name: string; calories: number; cooking_time: number }
    totalCalories: number
    totalProteins: number
    totalFats: number
    totalCarbs: number
  }>
}

export default function MenusPage() {
  const [menus, setMenus] = useState<Menu[]>([])
  const [weeklyMenus, setWeeklyMenus] = useState<Array<{id: number, date: string, week: WeeklyMenu['week']}>>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<"daily" | "weekly">("daily")
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const loadMenus = () => {
    // Получаем список всех дневных меню пользователя
    api.get("/menus")
      .then((response) => {
        setMenus(response.data || [])
      })
      .catch((error) => {
        console.error("Error fetching menus:", error)
      })

    // Получаем все сохраненные недельные меню из базы данных
    api.get("/menus/weekly")
      .then((response) => {
        if (response.data && response.data.length > 0) {
          // Преобразуем все недельные меню
          const weeklyMenusData = response.data.map((menu: any) => {
            if (menu.meals && Array.isArray(menu.meals)) {
              return {
                id: menu.id,
                date: menu.date,
                week: menu.meals.map((day: any) => ({
                  day: day.day || 0,
                  breakfast: day.breakfast || null,
                  lunch: day.lunch || null,
                  dinner: day.dinner || null,
                  totalCalories: day.totalCalories || 0,
                  totalProteins: day.totalProteins || 0,
                  totalFats: day.totalFats || 0,
                  totalCarbs: day.totalCarbs || 0,
                  totalTime: day.totalTime || 0,
                  ingredients_used: day.ingredients_used || [],
                  missing_ingredients: day.missing_ingredients || [],
                }))
              }
            }
            return null
          }).filter(Boolean)
          setWeeklyMenus(weeklyMenusData)
        }
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching weekly menus:", error)
        setLoading(false)
      })
  }

  useEffect(() => {
    loadMenus()
  }, [])

  const handleDelete = async (menuId: number) => {
    if (!confirm("Вы уверены, что хотите удалить это меню?")) {
      return
    }

    setDeletingId(menuId)
    try {
      await api.delete(`/menus/${menuId}`)
      // Перезагружаем списки меню
      loadMenus()
    } catch (error: any) {
      console.error("Error deleting menu:", error)
      alert(error.response?.data?.error || "Не удалось удалить меню")
    } finally {
      setDeletingId(null)
    }
  }

  if (loading) {
    return <div className="container mx-auto px-4 py-8">Загрузка...</div>
  }

  const hasAnyMenus = menus.length > 0 || weeklyMenus.length > 0

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold">Мои меню</h1>
        <div className="flex gap-2">
          <Link href="/menu/generate">
            <Button variant="outline" size="sm">Создать меню на день</Button>
          </Link>
          <Link href="/menu/weekly">
            <Button variant="outline" size="sm">Создать меню на неделю</Button>
          </Link>
        </div>
      </div>

      {/* Табы */}
      <div className="flex border-b mb-6">
        <button
          onClick={() => setActiveTab("daily")}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === "daily"
              ? "border-b-2 border-primary text-primary"
              : "text-muted-foreground hover:text-foreground"
          }`}
        >
          Меню на день ({menus.length})
        </button>
        <button
          onClick={() => setActiveTab("weekly")}
          className={`px-4 py-2 font-medium transition-colors ${
            activeTab === "weekly"
              ? "border-b-2 border-primary text-primary"
              : "text-muted-foreground hover:text-foreground"
          }`}
        >
          Меню на неделю ({weeklyMenus.length})
        </button>
      </div>

      {!hasAnyMenus ? (
        <Card>
          <CardContent className="pt-6">
            <p className="text-muted-foreground text-center">
              У вас пока нет созданных меню. <Link href="/menu/generate" className="text-primary underline">Создайте меню</Link>, чтобы начать.
            </p>
          </CardContent>
        </Card>
      ) : (
        <div>
          {/* Дневные меню */}
          {activeTab === "daily" && (
            <div>
              {menus.length === 0 ? (
                <Card>
                  <CardContent className="pt-6">
                    <p className="text-muted-foreground text-center">
                      У вас пока нет дневных меню. <Link href="/menu/generate" className="text-primary underline">Создайте меню</Link>, чтобы начать.
                    </p>
                  </CardContent>
                </Card>
              ) : (
                <div className="grid gap-4">
                  {menus.map((menu) => (
                    <Card key={menu.id}>
                      <CardHeader>
                        <CardTitle>Меню на {new Date(menu.date).toLocaleDateString('ru-RU')}</CardTitle>
                        <CardDescription>
                          {menu.total_calories} ккал • {menu.total_time} мин
                        </CardDescription>
                      </CardHeader>
                      <CardContent>
                        <div className="flex flex-col sm:flex-row gap-2">
                          <Link href={`/menu/${menu.id}`} className="flex-1">
                            <Button className="w-full">Просмотреть меню</Button>
                          </Link>
                          <Link href={`/shopping-list/${menu.id}`} className="flex-1">
                            <Button className="w-full" variant="outline">Список покупок</Button>
                          </Link>
                          <Button 
                            className="w-full sm:w-auto" 
                            variant="destructive"
                            onClick={() => handleDelete(menu.id)}
                            disabled={deletingId === menu.id}
                          >
                            {deletingId === menu.id ? "Удаление..." : "Удалить"}
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Недельные меню */}
          {activeTab === "weekly" && (
            <div>
              {weeklyMenus.length === 0 ? (
                <Card>
                  <CardContent className="pt-6">
                    <p className="text-muted-foreground text-center">
                      У вас пока нет недельных меню. <Link href="/menu/weekly" className="text-primary underline">Создайте меню</Link>, чтобы начать.
                    </p>
                  </CardContent>
                </Card>
              ) : (
                <div className="grid gap-4">
                  {weeklyMenus.map((weeklyMenu) => (
                    <Card key={weeklyMenu.id}>
                      <CardHeader>
                        <CardTitle>Меню на неделю</CardTitle>
                        <CardDescription>
                          Создано {new Date(weeklyMenu.date).toLocaleDateString('ru-RU')} • {weeklyMenu.week.length} дней • {weeklyMenu.week.reduce((sum, day) => sum + (day.totalCalories || 0), 0)} ккал всего
                        </CardDescription>
                      </CardHeader>
                      <CardContent>
                        <div className="flex flex-col sm:flex-row gap-2">
                          <Button 
                            className="flex-1" 
                            onClick={() => {
                              // Сохраняем недельное меню в localStorage для отображения на странице просмотра
                              const menuToView: WeeklyMenu = {
                                id: weeklyMenu.id,
                                week: weeklyMenu.week
                              }
                              localStorage.setItem("viewWeeklyMenu", JSON.stringify(menuToView))
                              window.location.href = "/menu/weekly"
                            }}
                          >
                            Просмотреть меню
                          </Button>
                          <Button 
                            className="w-full sm:w-auto" 
                            variant="destructive"
                            onClick={() => handleDelete(weeklyMenu.id!)}
                            disabled={deletingId === weeklyMenu.id}
                          >
                            {deletingId === weeklyMenu.id ? "Удаление..." : "Удалить"}
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

