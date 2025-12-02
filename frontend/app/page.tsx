import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center mb-12">
          <h1 className="text-3xl sm:text-4xl font-bold mb-4">MyPlateService</h1>
          <p className="text-lg sm:text-xl text-muted-foreground">
            Помощник по рецептам и ежедневному меню
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 max-w-5xl mx-auto">
          <Card>
            <CardHeader>
              <CardTitle>Рецепты</CardTitle>
              <CardDescription>Просмотрите нашу коллекцию рецептов</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/recipes">
                <Button className="w-full">Просмотреть рецепты</Button>
              </Link>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Генератор меню</CardTitle>
              <CardDescription>Создайте ваше ежедневное меню</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/menu/generate">
                <Button className="w-full">Создать меню</Button>
              </Link>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Мои цели</CardTitle>
              <CardDescription>Установите ваши цели по питанию</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/goals">
                <Button className="w-full">Установить цели</Button>
              </Link>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Кладовая</CardTitle>
              <CardDescription>Управляйте вашими ингредиентами</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/pantry">
                <Button className="w-full">Просмотреть кладовую</Button>
              </Link>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Мои меню</CardTitle>
              <CardDescription>Просмотрите все ваши созданные меню</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/menus">
                <Button className="w-full">Просмотреть меню</Button>
              </Link>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Список покупок</CardTitle>
              <CardDescription>Просмотрите ваши списки покупок</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/shopping-list">
                <Button className="w-full" variant="outline">Просмотреть списки</Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}


