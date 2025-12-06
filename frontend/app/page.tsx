import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { 
  BookOpen, 
  UtensilsCrossed, 
  Calendar, 
  Target, 
  Package, 
  ClipboardList, 
  ShoppingCart,
  Sparkles
} from "lucide-react"

const features = [
  {
    title: "Рецепты",
    description: "Просмотрите нашу коллекцию рецептов",
    href: "/recipes",
    icon: BookOpen,
    gradient: "from-blue-500 to-cyan-500",
  },
  {
    title: "Генератор меню",
    description: "Создайте ваше ежедневное меню",
    href: "/menu/generate",
    icon: UtensilsCrossed,
    gradient: "from-green-500 to-emerald-500",
  },
  {
    title: "Меню на неделю",
    description: "Создайте меню на 7 дней с учетом семьи",
    href: "/menu/weekly",
    icon: Calendar,
    gradient: "from-purple-500 to-pink-500",
  },
  {
    title: "Мои цели",
    description: "Установите ваши цели по питанию",
    href: "/goals",
    icon: Target,
    gradient: "from-orange-500 to-red-500",
  },
  {
    title: "Кладовая",
    description: "Управляйте вашими ингредиентами",
    href: "/pantry",
    icon: Package,
    gradient: "from-amber-500 to-yellow-500",
  },
  {
    title: "Мои меню",
    description: "Просмотрите все ваши созданные меню",
    href: "/menus",
    icon: ClipboardList,
    gradient: "from-indigo-500 to-blue-500",
  },
  {
    title: "Список покупок",
    description: "Просмотрите ваши списки покупок",
    href: "/shopping-list",
    icon: ShoppingCart,
    gradient: "from-teal-500 to-green-500",
  },
]

export default function Home() {
  return (
    <div className="min-h-screen gradient-bg">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center mb-16">
          <div className="inline-flex items-center gap-2 mb-4">
            <Sparkles className="h-8 w-8 text-primary" />
            <h1 className="text-4xl sm:text-5xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
              MyPlateService
            </h1>
          </div>
          <p className="text-lg sm:text-xl text-muted-foreground max-w-2xl mx-auto">
            Помощник по рецептам и ежедневному меню
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 max-w-6xl mx-auto">
          {features.map((feature) => {
            const Icon = feature.icon
            return (
              <Card 
                key={feature.href}
                className="flex flex-col card-hover group overflow-hidden relative"
              >
                <div className={`absolute inset-0 bg-gradient-to-br ${feature.gradient} opacity-0 group-hover:opacity-5 transition-opacity duration-300`} />
                <CardHeader className="relative">
                  <div className={`inline-flex items-center justify-center w-12 h-12 rounded-lg bg-gradient-to-br ${feature.gradient} mb-4 shadow-lg`}>
                    <Icon className="h-6 w-6 text-white" />
                  </div>
                  <CardTitle className="text-xl">{feature.title}</CardTitle>
                  <CardDescription className="text-sm">{feature.description}</CardDescription>
                </CardHeader>
                <CardContent className="flex-1 flex flex-col justify-end relative">
                  <Link href={feature.href}>
                    <Button className="w-full group-hover:shadow-lg transition-all duration-300">
                      Открыть
                    </Button>
                  </Link>
                </CardContent>
              </Card>
            )
          })}
        </div>
      </div>
    </div>
  )
}


