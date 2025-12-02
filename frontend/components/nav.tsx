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

  const navItems = [
    { href: "/", label: "Главная" },
    { href: "/recipes", label: "Рецепты" },
    { href: "/menus", label: "Меню" },
    { href: "/menu/generate", label: "Создать меню" },
    { href: "/goals", label: "Цели" },
    { href: "/pantry", label: "Кладовая" },
  ]

  return (
    <nav className="sticky top-0 z-50 border-b bg-background">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link href="/" className="text-xl font-bold">
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
              </nav>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </nav>
  )
}


