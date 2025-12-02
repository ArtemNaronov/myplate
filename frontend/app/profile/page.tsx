"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import api from "@/lib/api"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface User {
  id: number
  email: string
  username: string
  first_name: string
  last_name: string
}

export default function ProfilePage() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  
  const [profileData, setProfileData] = useState({
    firstName: "",
    lastName: "",
  })
  
  const [passwordData, setPasswordData] = useState({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  })
  const [showPasswordForm, setShowPasswordForm] = useState(false)

  useEffect(() => {
    loadProfile()
  }, [])

  const loadProfile = async () => {
    try {
      const response = await api.get("/auth/profile")
      setUser(response.data)
      setProfileData({
        firstName: response.data.first_name || "",
        lastName: response.data.last_name || "",
      })
    } catch (err: any) {
      if (err.response?.status === 401) {
        router.push("/auth/login")
      } else {
        setError("Ошибка при загрузке профиля")
      }
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setSuccess("")

    try {
      const response = await api.put("/auth/profile", {
        first_name: profileData.firstName,
        last_name: profileData.lastName,
      })
      setUser(response.data)
      setSuccess("Профиль успешно обновлен")
    } catch (err: any) {
      setError(err.response?.data?.error || "Ошибка при обновлении профиля")
    }
  }

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setSuccess("")

    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setError("Новые пароли не совпадают")
      return
    }

    if (passwordData.newPassword.length < 6) {
      setError("Пароль должен содержать минимум 6 символов")
      return
    }

    try {
      await api.put("/auth/password", {
        old_password: passwordData.oldPassword,
        new_password: passwordData.newPassword,
      })
      setSuccess("Пароль успешно обновлен")
      setPasswordData({
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
      })
      setShowPasswordForm(false)
    } catch (err: any) {
      setError(err.response?.data?.error || "Ошибка при обновлении пароля")
    }
  }

  const handleLogout = () => {
    localStorage.removeItem("token")
    router.push("/auth/login")
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">Загрузка...</div>
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-3xl font-bold mb-6">Профиль</h1>

      {error && (
        <div className="mb-4 p-4 bg-red-50 text-red-700 rounded">
          {error}
        </div>
      )}

      {success && (
        <div className="mb-4 p-4 bg-green-50 text-green-700 rounded">
          {success}
        </div>
      )}

      <div className="space-y-6">
        {/* Информация о пользователе */}
        <Card>
          <CardHeader>
            <CardTitle>Информация о пользователе</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <p>
              <span className="font-semibold">Email:</span> {user?.email}
            </p>
            <p>
              <span className="font-semibold">Имя:</span> {user?.first_name || "Не указано"}
            </p>
            <p>
              <span className="font-semibold">Фамилия:</span> {user?.last_name || "Не указано"}
            </p>
          </CardContent>
        </Card>

        {/* Обновление профиля */}
        <Card>
          <CardHeader>
            <CardTitle>Обновить профиль</CardTitle>
            <CardDescription>Измените ваше имя и фамилию</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleUpdateProfile} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">Имя</Label>
                <Input
                  id="firstName"
                  type="text"
                  value={profileData.firstName}
                  onChange={(e) =>
                    setProfileData({ ...profileData, firstName: e.target.value })
                  }
                  placeholder="Ваше имя"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="lastName">Фамилия</Label>
                <Input
                  id="lastName"
                  type="text"
                  value={profileData.lastName}
                  onChange={(e) =>
                    setProfileData({ ...profileData, lastName: e.target.value })
                  }
                  placeholder="Ваша фамилия"
                />
              </div>

              <Button type="submit">Сохранить изменения</Button>
            </form>
          </CardContent>
        </Card>

        {/* Изменение пароля */}
        <Card>
          <CardHeader>
            <CardTitle>Изменить пароль</CardTitle>
            <CardDescription>Обновите ваш пароль</CardDescription>
          </CardHeader>
          <CardContent>
            {!showPasswordForm ? (
              <Button onClick={() => setShowPasswordForm(true)}>
                Изменить пароль
              </Button>
            ) : (
              <form onSubmit={handleUpdatePassword} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="oldPassword">Текущий пароль</Label>
                  <Input
                    id="oldPassword"
                    type="password"
                    required
                    value={passwordData.oldPassword}
                    onChange={(e) =>
                      setPasswordData({
                        ...passwordData,
                        oldPassword: e.target.value,
                      })
                    }
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="newPassword">Новый пароль</Label>
                  <Input
                    id="newPassword"
                    type="password"
                    required
                    minLength={6}
                    value={passwordData.newPassword}
                    onChange={(e) =>
                      setPasswordData({
                        ...passwordData,
                        newPassword: e.target.value,
                      })
                    }
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="confirmPassword">Подтвердите новый пароль</Label>
                  <Input
                    id="confirmPassword"
                    type="password"
                    required
                    minLength={6}
                    value={passwordData.confirmPassword}
                    onChange={(e) =>
                      setPasswordData({
                        ...passwordData,
                        confirmPassword: e.target.value,
                      })
                    }
                  />
                </div>

                <div className="flex gap-2">
                  <Button type="submit">Сохранить пароль</Button>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => {
                      setShowPasswordForm(false)
                      setPasswordData({
                        oldPassword: "",
                        newPassword: "",
                        confirmPassword: "",
                      })
                    }}
                  >
                    Отмена
                  </Button>
                </div>
              </form>
            )}
          </CardContent>
        </Card>

        {/* Выход */}
        <Card>
          <CardContent className="pt-6">
            <Button variant="outline" onClick={handleLogout} className="w-full">
              Выйти
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

