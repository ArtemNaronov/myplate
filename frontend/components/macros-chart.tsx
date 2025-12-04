"use client"

interface MacrosChartProps {
  proteins: number
  fats: number
  carbs: number
  size?: number
  showLabels?: boolean
}

export default function MacrosChart({ 
  proteins, 
  fats, 
  carbs, 
  size = 200,
  showLabels = true 
}: MacrosChartProps) {
  // Калории из макронутриентов: белки и углеводы = 4 ккал/г, жиры = 9 ккал/г
  const proteinCalories = proteins * 4
  const fatCalories = fats * 9
  const carbCalories = carbs * 4
  const totalCalories = proteinCalories + fatCalories + carbCalories

  // Проценты от общего количества калорий
  const proteinPercent = totalCalories > 0 ? (proteinCalories / totalCalories) * 100 : 0
  const fatPercent = totalCalories > 0 ? (fatCalories / totalCalories) * 100 : 0
  const carbPercent = totalCalories > 0 ? (carbCalories / totalCalories) * 100 : 0

  // Радиус круга
  const radius = size / 2 - 10
  const center = size / 2
  const circumference = 2 * Math.PI * radius

  // Вычисляем offset для каждого сегмента
  const proteinOffset = circumference - (proteinPercent / 100) * circumference
  const fatOffset = circumference - (fatPercent / 100) * circumference
  const carbOffset = circumference - (carbPercent / 100) * circumference

  // Начальные углы для каждого сегмента
  let currentAngle = -90 // Начинаем сверху
  const proteinAngle = (proteinPercent / 100) * 360
  const fatAngle = (fatPercent / 100) * 360
  const carbAngle = (carbPercent / 100) * 360

  // Функция для создания path для сегмента круговой диаграммы
  const createArcPath = (startAngle: number, endAngle: number) => {
    if (endAngle - startAngle <= 0) return ""
    const startRad = (startAngle * Math.PI) / 180
    const endRad = (endAngle * Math.PI) / 180
    const x1 = center + radius * Math.cos(startRad)
    const y1 = center + radius * Math.sin(startRad)
    const x2 = center + radius * Math.cos(endRad)
    const y2 = center + radius * Math.sin(endRad)
    const largeArc = endAngle - startAngle > 180 ? 1 : 0
    return `M ${center} ${center} L ${x1} ${y1} A ${radius} ${radius} 0 ${largeArc} 1 ${x2} ${y2} Z`
  }

  if (totalCalories === 0) {
    return (
      <div className="flex flex-col items-center justify-center" style={{ width: size, height: size }}>
        <div className="text-sm text-muted-foreground">Нет данных</div>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <div className="relative" style={{ width: size, height: size }}>
        <svg width={size} height={size} className="transform -rotate-90">
          {/* Белки (синий) */}
          {proteinPercent > 0 && (
            <path
              d={createArcPath(currentAngle, currentAngle + proteinAngle)}
              fill="#3b82f6"
              className="transition-all duration-300"
            />
          )}
          {/* Жиры (оранжевый) */}
          {fatPercent > 0 && (
            <path
              d={createArcPath(
                currentAngle + proteinAngle,
                currentAngle + proteinAngle + fatAngle
              )}
              fill="#f97316"
              className="transition-all duration-300"
            />
          )}
          {/* Углеводы (зеленый) */}
          {carbPercent > 0 && (
            <path
              d={createArcPath(
                currentAngle + proteinAngle + fatAngle,
                currentAngle + proteinAngle + fatAngle + carbAngle
              )}
              fill="#22c55e"
              className="transition-all duration-300"
            />
          )}
        </svg>
      </div>
      {/* Текст с калориями под диаграммой */}
      <div className="flex flex-col items-center">
        <div className="text-2xl font-bold text-foreground">{Math.round(totalCalories)}</div>
        <div className="text-xs text-muted-foreground">ккал</div>
      </div>

      {showLabels && (
        <div className="flex flex-col gap-2 w-full">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 rounded bg-blue-500"></div>
              <span className="text-sm">Белки</span>
            </div>
            <div className="text-sm font-semibold">
              {proteins.toFixed(1)}г ({proteinPercent.toFixed(1)}%)
            </div>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 rounded bg-orange-500"></div>
              <span className="text-sm">Жиры</span>
            </div>
            <div className="text-sm font-semibold">
              {fats.toFixed(1)}г ({fatPercent.toFixed(1)}%)
            </div>
          </div>
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-4 h-4 rounded bg-green-500"></div>
              <span className="text-sm">Углеводы</span>
            </div>
            <div className="text-sm font-semibold">
              {carbs.toFixed(1)}г ({carbPercent.toFixed(1)}%)
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

