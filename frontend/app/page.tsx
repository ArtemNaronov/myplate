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
          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Рецепты</CardTitle>
              <CardDescription>Просмотрите нашу коллекцию рецептов</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
              <Link href="/recipes">
                <Button className="w-full">Просмотреть рецепты</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Генератор меню</CardTitle>
              <CardDescription>Создайте ваше ежедневное меню</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
              <Link href="/menu/generate">
                <Button className="w-full">Создать меню</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Меню на неделю</CardTitle>
              <CardDescription>Создайте меню на 7 дней с учетом семьи</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
              <Link href="/menu/weekly">
                <Button className="w-full">Создать недельное меню</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Мои цели</CardTitle>
              <CardDescription>Установите ваши цели по питанию</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
              <Link href="/goals">
                <Button className="w-full">Установить цели</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Кладовая</CardTitle>
              <CardDescription>Управляйте вашими ингредиентами</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
              <Link href="/pantry">
                <Button className="w-full">Просмотреть кладовую</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Мои меню</CardTitle>
              <CardDescription>Просмотрите все ваши созданные меню</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
              <Link href="/menus">
                <Button className="w-full">Просмотреть меню</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="flex flex-col">
            <CardHeader>
              <CardTitle>Список покупок</CardTitle>
              <CardDescription>Просмотрите ваши списки покупок</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col justify-end">
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


