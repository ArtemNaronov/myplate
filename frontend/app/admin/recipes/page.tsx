"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import api from "@/lib/api"
import { useTelegram } from "@/components/telegram-provider"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"

interface RecipeImportDTO {
  title: string
  description: string
  tags: string[]
  ingredients: Array<{ name: string; amount: number; unit: string }>
  calories: number
  proteins: number
  fats: number
  carbs: number
  cooking_time?: number
  servings?: number
  instructions?: string[]
}

export default function AdminRecipesPage() {
  const router = useRouter()
  const { user } = useTelegram()
  const [loading, setLoading] = useState(false)
  const [isAdmin, setIsAdmin] = useState(false)
  const [activeTab, setActiveTab] = useState<"create" | "import" | "export">("create")
  const [createForm, setCreateForm] = useState<RecipeImportDTO>({
    title: "",
    description: "",
    tags: [],
    ingredients: [],
    calories: 0,
    proteins: 0,
    fats: 0,
    carbs: 0,
    cooking_time: 30,
    servings: 1,
    instructions: [],
  })
  const [importJson, setImportJson] = useState("")
  const [exportData, setExportData] = useState<any>(null)
  const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null)

  useEffect(() => {
    // Проверяем роль пользователя
    if (user?.role === "admin") {
      setIsAdmin(true)
    } else {
      // Пытаемся получить роль из токена
      const token = localStorage.getItem("token")
      if (token) {
        // Декодируем JWT (простая проверка)
        try {
          const payload = JSON.parse(atob(token.split(".")[1]))
          if (payload.role === "admin") {
            setIsAdmin(true)
          } else {
            router.push("/")
          }
        } catch {
          router.push("/")
        }
      } else {
        router.push("/auth/login")
      }
    }
  }, [user, router])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setMessage(null)

    try {
      await api.post("/admin/recipes", createForm)
      setMessage({ type: "success", text: "Рецепт успешно создан!" })
      setCreateForm({
        title: "",
        description: "",
        tags: [],
        ingredients: [],
        calories: 0,
        proteins: 0,
        fats: 0,
        carbs: 0,
        cooking_time: 30,
        servings: 1,
        instructions: [],
      })
    } catch (error: any) {
      setMessage({
        type: "error",
        text: error.response?.data?.error || "Ошибка при создании рецепта",
      })
    } finally {
      setLoading(false)
    }
  }

  const handleImport = async () => {
    setLoading(true)
    setMessage(null)

    try {
      const data = JSON.parse(importJson)
      const response = await api.post("/admin/recipes/import", data)
      setMessage({
        type: "success",
        text: `Импортировано: ${response.data.imported}, Ошибок: ${response.data.failed}`,
      })
      if (response.data.errors && response.data.errors.length > 0) {
        console.error("Ошибки импорта:", response.data.errors)
      }
      setImportJson("")
    } catch (error: any) {
      setMessage({
        type: "error",
        text: error.response?.data?.error || "Ошибка при импорте рецептов",
      })
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async () => {
    setLoading(true)
    setMessage(null)

    try {
      const response = await api.get("/admin/recipes/export")
      setExportData(response.data)
      setMessage({ type: "success", text: "Рецепты успешно экспортированы!" })
    } catch (error: any) {
      setMessage({
        type: "error",
        text: error.response?.data?.error || "Ошибка при экспорте рецептов",
      })
    } finally {
      setLoading(false)
    }
  }

  const downloadExport = () => {
    if (!exportData) return
    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: "application/json" })
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = "recipes_export.json"
    a.click()
    URL.revokeObjectURL(url)
  }

  if (!isAdmin) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Card>
          <CardContent className="pt-6">
            <p className="text-center text-muted-foreground">Проверка прав доступа...</p>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <h1 className="text-2xl sm:text-3xl font-bold mb-6">Управление рецептами (Админ)</h1>

      {message && (
        <div
          className={`mb-4 p-4 rounded-md ${
            message.type === "success" ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"
          }`}
        >
          {message.text}
        </div>
      )}

      <div className="flex gap-2 mb-6">
        <Button
          variant={activeTab === "create" ? "default" : "outline"}
          onClick={() => setActiveTab("create")}
        >
          Создать рецепт
        </Button>
        <Button
          variant={activeTab === "import" ? "default" : "outline"}
          onClick={() => setActiveTab("import")}
        >
          Импорт JSON
        </Button>
        <Button
          variant={activeTab === "export" ? "default" : "outline"}
          onClick={() => setActiveTab("export")}
        >
          Экспорт JSON
        </Button>
      </div>

      {activeTab === "create" && (
        <Card>
          <CardHeader>
            <CardTitle>Создать новый рецепт</CardTitle>
            <CardDescription>Заполните форму для создания рецепта</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreate} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Название *</label>
                <input
                  type="text"
                  value={createForm.title}
                  onChange={(e) => setCreateForm({ ...createForm, title: e.target.value })}
                  className="w-full px-3 py-2 border rounded-md"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Описание</label>
                <textarea
                  value={createForm.description}
                  onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
                  className="w-full px-3 py-2 border rounded-md"
                  rows={3}
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">
                  Теги (через запятую: breakfast, lunch, dinner, vegetarian, vegan, eggs, dairy, etc.)
                </label>
                <input
                  type="text"
                  value={createForm.tags.join(", ")}
                  onChange={(e) =>
                    setCreateForm({
                      ...createForm,
                      tags: e.target.value.split(",").map((t) => t.trim()).filter(Boolean),
                    })
                  }
                  placeholder="breakfast, vegetarian, eggs"
                  className="w-full px-3 py-2 border rounded-md"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Калории *</label>
                  <input
                    type="number"
                    value={createForm.calories}
                    onChange={(e) =>
                      setCreateForm({ ...createForm, calories: parseFloat(e.target.value) || 0 })
                    }
                    className="w-full px-3 py-2 border rounded-md"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Время приготовления (мин)</label>
                  <input
                    type="number"
                    value={createForm.cooking_time}
                    onChange={(e) =>
                      setCreateForm({ ...createForm, cooking_time: parseInt(e.target.value) || 30 })
                    }
                    className="w-full px-3 py-2 border rounded-md"
                  />
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-2">Белки (г)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={createForm.proteins}
                    onChange={(e) =>
                      setCreateForm({ ...createForm, proteins: parseFloat(e.target.value) || 0 })
                    }
                    className="w-full px-3 py-2 border rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Жиры (г)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={createForm.fats}
                    onChange={(e) =>
                      setCreateForm({ ...createForm, fats: parseFloat(e.target.value) || 0 })
                    }
                    className="w-full px-3 py-2 border rounded-md"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-2">Углеводы (г)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={createForm.carbs}
                    onChange={(e) =>
                      setCreateForm({ ...createForm, carbs: parseFloat(e.target.value) || 0 })
                    }
                    className="w-full px-3 py-2 border rounded-md"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Количество порций</label>
                <input
                  type="number"
                  min="1"
                  value={createForm.servings}
                  onChange={(e) =>
                    setCreateForm({ ...createForm, servings: parseInt(e.target.value) || 1 })
                  }
                  className="w-full px-3 py-2 border rounded-md"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Ингредиенты (JSON)</label>
                <textarea
                  value={JSON.stringify(createForm.ingredients, null, 2)}
                  onChange={(e) => {
                    try {
                      const ingredients = JSON.parse(e.target.value)
                      setCreateForm({ ...createForm, ingredients })
                    } catch {
                      // Игнорируем ошибки парсинга
                    }
                  }}
                  className="w-full px-3 py-2 border rounded-md font-mono text-sm"
                  rows={5}
                  placeholder='[{"name": "Яйца", "amount": 2, "unit": "шт"}]'
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Инструкции (каждая строка - шаг)</label>
                <textarea
                  value={createForm.instructions?.join("\n") || ""}
                  onChange={(e) =>
                    setCreateForm({
                      ...createForm,
                      instructions: e.target.value.split("\n").filter(Boolean),
                    })
                  }
                  className="w-full px-3 py-2 border rounded-md"
                  rows={4}
                />
              </div>

              <Button type="submit" className="w-full" disabled={loading}>
                {loading ? "Создание..." : "Создать рецепт"}
              </Button>
            </form>
          </CardContent>
        </Card>
      )}

      {activeTab === "import" && (
        <Card>
          <CardHeader>
            <CardTitle>Импорт рецептов из JSON</CardTitle>
            <CardDescription>Вставьте JSON с рецептами для импорта</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">JSON данные</label>
                <textarea
                  value={importJson}
                  onChange={(e) => setImportJson(e.target.value)}
                  className="w-full px-3 py-2 border rounded-md font-mono text-sm"
                  rows={15}
                  placeholder='{"recipes": [{"title": "...", ...}]}'
                />
              </div>
              <Button onClick={handleImport} className="w-full" disabled={loading || !importJson}>
                {loading ? "Импорт..." : "Импортировать рецепты"}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {activeTab === "export" && (
        <Card>
          <CardHeader>
            <CardTitle>Экспорт рецептов в JSON</CardTitle>
            <CardDescription>Экспортируйте все рецепты в JSON формат</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <Button onClick={handleExport} className="w-full" disabled={loading}>
                {loading ? "Экспорт..." : "Экспортировать рецепты"}
              </Button>
              {exportData && (
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <p className="text-sm font-medium">Экспортировано рецептов: {exportData.recipes?.length || 0}</p>
                    <Button onClick={downloadExport} variant="outline" size="sm">
                      Скачать JSON
                    </Button>
                  </div>
                  <textarea
                    value={JSON.stringify(exportData, null, 2)}
                    readOnly
                    className="w-full px-3 py-2 border rounded-md font-mono text-xs"
                    rows={20}
                  />
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

