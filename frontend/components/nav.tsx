"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet"
import { Menu } from "lucide-react"
import { ThemeToggle } from "@/components/theme-toggle"

export function Nav() {
  const pathname = usePathname()
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const [isAuthenticated, setIsAuthenticated] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem("token")
    setIsAuthenticated(!!token)
  }, [pathname])

  const handleLogout = () => {
    localStorage.removeItem("token")
    setIsAuthenticated(false)
    router.push("/auth/login")
  }

  const [userRole, setUserRole] = useState<string | null>(null)

  useEffect(() => {
    // Получаем роль из токена
    const token = localStorage.getItem("token")
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split(".")[1]))
        setUserRole(payload.role || "user")
      } catch {
        setUserRole("user")
      }
    }
  }, [pathname])

  const navItems = [
    { href: "/", label: "Главная" },
    { href: "/recipes", label: "Рецепты" },
    { href: "/menus", label: "Меню" },
    { href: "/menu/generate", label: "Создать меню" },
    { href: "/menu/weekly", label: "Меню на неделю" },
    { href: "/goals", label: "Цели" },
    { href: "/pantry", label: "Кладовая" },
  ]

  return (
    <nav className="sticky top-0 z-50 border-b glass backdrop-blur-xl bg-background/80">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="text-xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent hover:opacity-80 transition-opacity">
            MyPlateService
          </Link>
          
          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-2">
            {navItems.map((item) => (
              <Link key={item.href} href={item.href}>
                <Button
                  variant={pathname === item.href ? "default" : "ghost"}
                  size="sm"
                >
                  {item.label}
                </Button>
              </Link>
            ))}
            {isAuthenticated ? (
              <>
                {userRole === "admin" && (
                  <Link href="/admin/recipes">
                    <Button
                      variant={pathname === "/admin/recipes" ? "default" : "ghost"}
                      size="sm"
                    >
                      Админ
                    </Button>
                  </Link>
                )}
                <Link href="/profile">
                  <Button
                    variant={pathname === "/profile" ? "default" : "ghost"}
                    size="sm"
                  >
                    Профиль
                  </Button>
                </Link>
                <Button variant="ghost" size="sm" onClick={handleLogout}>
                  Выйти
                </Button>
              </>
            ) : (
              <Link href="/auth/login">
                <Button variant="outline" size="sm">
                  Войти
                </Button>
              </Link>
            )}
            <ThemeToggle />
          </div>

          {/* Mobile Navigation */}
          <Sheet open={open} onOpenChange={setOpen}>
            <SheetTrigger asChild className="md:hidden">
              <Button variant="ghost" size="icon">
                <Menu className="h-6 w-6" />
                <span className="sr-only">Открыть меню</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="right" className="w-[300px] sm:w-[400px]">
              <SheetHeader>
                <SheetTitle>Навигация</SheetTitle>
              </SheetHeader>
              <nav className="flex flex-col space-y-2 mt-6">
                {navItems.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    onClick={() => setOpen(false)}
                  >
                    <Button
                      variant={pathname === item.href ? "default" : "ghost"}
                      className="w-full justify-start"
                    >
                      {item.label}
                    </Button>
                  </Link>
                ))}
                {isAuthenticated ? (
                  <>
                    {userRole === "admin" && (
                      <Link href="/admin/recipes" onClick={() => setOpen(false)}>
                        <Button
                          variant={pathname === "/admin/recipes" ? "default" : "ghost"}
                          className="w-full justify-start"
                        >
                          Админ
                        </Button>
                      </Link>
                    )}
                    <Link href="/profile" onClick={() => setOpen(false)}>
                      <Button
                        variant={pathname === "/profile" ? "default" : "ghost"}
                        className="w-full justify-start"
                      >
                        Профиль
                      </Button>
                    </Link>
                    <Button
                      variant="ghost"
                      className="w-full justify-start"
                      onClick={() => {
                        handleLogout()
                        setOpen(false)
                      }}
                    >
                      Выйти
                    </Button>
                  </>
                ) : (
                  <Link href="/auth/login" onClick={() => setOpen(false)}>
                    <Button variant="outline" className="w-full justify-start">
                      Войти
                    </Button>
                  </Link>
                )}
                <div className="pt-4 border-t">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">Тема</span>
                    <ThemeToggle />
                  </div>
                </div>
              </nav>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </nav>
  )
}


