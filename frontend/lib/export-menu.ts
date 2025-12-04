// Динамический импорт html2pdf только на клиенте

interface DailyMenuData {
  date: string
  total_calories: number
  total_time: number
  meals: Array<{
    meal_type: string
    recipe_name: string
    calories: number
    time: number
    proteins?: number
    fats?: number
    carbs?: number
    ingredients?: Array<{ name: string; quantity: number | string; unit: string }>
  }>
  totalProteins?: number
  totalFats?: number
  totalCarbs?: number
}

// Интерфейс для данных недельного меню
interface WeeklyMenuData {
  week: Array<{
    day: number
    breakfast: { 
      id?: number
      name: string
      calories: number
      cooking_time: number
      ingredients?: Array<{ name: string; quantity: number | string; unit: string }>
      [key: string]: any // Разрешаем дополнительные поля
    }
    lunch: { 
      id?: number
      name: string
      calories: number
      cooking_time: number
      ingredients?: Array<{ name: string; quantity: number | string; unit: string }>
      [key: string]: any // Разрешаем дополнительные поля
    }
    dinner: { 
      id?: number
      name: string
      calories: number
      cooking_time: number
      ingredients?: Array<{ name: string; quantity: number | string; unit: string }>
      [key: string]: any // Разрешаем дополнительные поля
    }
    totalCalories: number
    totalProteins: number
    totalFats: number
    totalCarbs: number
    totalTime?: number
  }>
}

const getMealTypeLabel = (type: string) => {
  switch (type) {
    case 'breakfast': return 'Завтрак'
    case 'lunch': return 'Обед'
    case 'dinner': return 'Ужин'
    default: return type
  }
}

export async function exportDailyMenuToPDF(menu: DailyMenuData) {
  // Динамический импорт только на клиенте
  const html2pdf = (await import('html2pdf.js')).default
  const dateStr = new Date(menu.date).toLocaleDateString('ru-RU', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })

  // Создаем HTML структуру
  const htmlContent = `
    <div style="font-family: Arial, sans-serif; padding: 20px; max-width: 800px; margin: 0 auto;">
      <h1 style="text-align: center; font-size: 24px; margin-bottom: 10px;">Меню на день</h1>
      <p style="text-align: center; font-size: 14px; color: #666; margin-bottom: 30px;">${dateStr}</p>
      
      <div style="margin-bottom: 30px;">
        <h2 style="font-size: 18px; margin-bottom: 15px; border-bottom: 2px solid #333; padding-bottom: 5px;">Общая информация</h2>
        <p style="margin: 5px 0;"><strong>Калории:</strong> ${menu.total_calories} ккал</p>
        <p style="margin: 5px 0;"><strong>Время приготовления:</strong> ${menu.total_time} минут</p>
        ${menu.totalProteins !== undefined && menu.totalFats !== undefined && menu.totalCarbs !== undefined ? `
          <h3 style="font-size: 16px; margin-top: 15px; margin-bottom: 10px;">Макронутриенты</h3>
          <p style="margin: 5px 0;"><strong>Белки:</strong> ${menu.totalProteins.toFixed(1)}г</p>
          <p style="margin: 5px 0;"><strong>Жиры:</strong> ${menu.totalFats.toFixed(1)}г</p>
          <p style="margin: 5px 0;"><strong>Углеводы:</strong> ${menu.totalCarbs.toFixed(1)}г</p>
        ` : ''}
      </div>
      
      <div>
        <h2 style="font-size: 18px; margin-bottom: 15px; border-bottom: 2px solid #333; padding-bottom: 5px;">Блюда</h2>
        ${menu.meals.map((meal, index) => `
          <div style="margin-bottom: 20px; padding: 15px; border: 1px solid #ddd; border-radius: 5px;">
            <div style="display: flex; justify-content: space-between; gap: 20px; align-items: flex-start;">
              <div style="flex: 1;">
                <h3 style="font-size: 16px; margin-bottom: 10px; color: #333;">${getMealTypeLabel(meal.meal_type)}</h3>
                <p style="margin: 5px 0; font-weight: 500;"><strong>Блюдо:</strong> ${meal.recipe_name}</p>
                <p style="margin: 5px 0;"><strong>Калории:</strong> ${meal.calories} ккал</p>
                <p style="margin: 5px 0;"><strong>Время:</strong> ${meal.time} мин</p>
                ${meal.proteins !== undefined && meal.fats !== undefined && meal.carbs !== undefined ? `
                  <p style="margin: 5px 0;"><strong>БЖУ:</strong> ${meal.proteins.toFixed(1)}г / ${meal.fats.toFixed(1)}г / ${meal.carbs.toFixed(1)}г</p>
                ` : ''}
              </div>
              ${(meal.ingredients && Array.isArray(meal.ingredients) && meal.ingredients.length > 0) ? `
                <div style="flex: 1; border-left: 1px solid #ddd; padding-left: 15px; min-width: 200px;">
                  <p style="margin: 0 0 8px 0; font-size: 12px; font-weight: bold; color: #555;">Ингредиенты:</p>
                  <ul style="margin: 0; padding-left: 20px; font-size: 11px; color: #666; list-style-type: disc;">
                    ${meal.ingredients.map((ing: any) => {
                      const qty = typeof ing.quantity === 'number' ? (ing.quantity % 1 === 0 ? ing.quantity.toString() : ing.quantity.toFixed(2)) : ing.quantity
                      return `<li style="margin-bottom: 3px;">${ing.name} - ${qty} ${ing.unit}</li>`
                    }).join('')}
                  </ul>
                </div>
              ` : '<div style="flex: 1;"></div>'}
            </div>
          </div>
        `).join('')}
      </div>
      
      <div style="margin-top: 30px; text-align: center; font-size: 10px; color: #999; border-top: 1px solid #ddd; padding-top: 10px;">
        MyPlateService
      </div>
    </div>
  `

  // Создаем временный элемент
  const element = document.createElement('div')
  element.innerHTML = htmlContent
  document.body.appendChild(element)

  // Настройки для PDF
  const opt = {
    margin: [10, 10, 10, 10] as [number, number, number, number],
    filename: `menu_${new Date(menu.date).toISOString().split('T')[0]}.pdf`,
    image: { type: 'jpeg' as const, quality: 0.98 },
    html2canvas: { scale: 2, useCORS: true },
    jsPDF: { unit: 'mm' as const, format: 'a4' as const, orientation: 'portrait' as const }
  }

  // Генерируем PDF
  html2pdf().set(opt).from(element).save().then(() => {
    // Удаляем временный элемент
    document.body.removeChild(element)
  })
}

export async function exportWeeklyMenuToPDF(menu: WeeklyMenuData) {
  // Динамический импорт только на клиенте
  const html2pdf = (await import('html2pdf.js')).default
  
  // Отладочный вывод для проверки данных
  if (menu.week && menu.week.length > 0) {
    const firstDay = menu.week[0]
    console.log('First day breakfast:', firstDay.breakfast)
    console.log('First day breakfast ingredients:', firstDay.breakfast?.ingredients)
    console.log('First day breakfast has ingredients?', !!firstDay.breakfast?.ingredients)
    console.log('First day breakfast ingredients length:', firstDay.breakfast?.ingredients?.length)
  }
  
  const weekTotal = menu.week.reduce((acc, day) => ({
    calories: acc.calories + (day.totalCalories || 0),
    proteins: acc.proteins + (day.totalProteins || 0),
    fats: acc.fats + (day.totalFats || 0),
    carbs: acc.carbs + (day.totalCarbs || 0),
    time: acc.time + (day.totalTime || 0)
  }), { calories: 0, proteins: 0, fats: 0, carbs: 0, time: 0 })

  const dayNames = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота', 'Воскресенье']

  // Создаем HTML структуру
  const htmlContent = `
    <div style="font-family: Arial, sans-serif; padding: 20px; max-width: 800px; margin: 0 auto;">
      <h1 style="text-align: center; font-size: 24px; margin-bottom: 30px;">Меню на неделю</h1>
      
      <div style="margin-bottom: 30px; padding: 15px; background-color: #f5f5f5; border-radius: 5px;">
        <h2 style="font-size: 18px; margin-bottom: 15px; border-bottom: 2px solid #333; padding-bottom: 5px;">Итого за неделю</h2>
        <p style="margin: 5px 0;"><strong>Калории:</strong> ${weekTotal.calories} ккал</p>
        <p style="margin: 5px 0;"><strong>Белки:</strong> ${weekTotal.proteins.toFixed(1)}г</p>
        <p style="margin: 5px 0;"><strong>Жиры:</strong> ${weekTotal.fats.toFixed(1)}г</p>
        <p style="margin: 5px 0;"><strong>Углеводы:</strong> ${weekTotal.carbs.toFixed(1)}г</p>
        ${weekTotal.time > 0 ? `<p style="margin: 5px 0;"><strong>Время приготовления:</strong> ${weekTotal.time} минут</p>` : ''}
      </div>
      
      <div>
        ${menu.week.map((day, dayIndex) => `
          <div style="margin-bottom: 15px; padding: 12px; border: 1px solid #ddd; border-radius: 5px; page-break-inside: avoid;">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
              <h2 style="font-size: 16px; margin: 0; color: #333;">
                День ${day.day} - ${dayNames[dayIndex] || `День ${day.day}`}
              </h2>
              <p style="margin: 0; font-size: 14px; color: #666; text-align: right;">
                Итого: ${day.totalCalories} ккал • Б: ${day.totalProteins.toFixed(1)}г • Ж: ${day.totalFats.toFixed(1)}г • У: ${day.totalCarbs.toFixed(1)}г
              </p>
            </div>
            
            <div style="margin-left: 10px;">
              <div style="margin-bottom: 8px;">
                <div style="display: flex; justify-content: space-between; gap: 20px; align-items: flex-start;">
                  <div style="flex: 1;">
                    <h3 style="font-size: 14px; margin-bottom: 5px; color: #555;"><strong>Завтрак:</strong></h3>
                    <p style="margin: 3px 0; margin-left: 10px; font-weight: 500;">${day.breakfast.name}</p>
                    <p style="margin: 3px 0; margin-left: 10px; font-size: 12px; color: #666;">${day.breakfast.calories} ккал • ${day.breakfast.cooking_time} мин</p>
                  </div>
                  ${(day.breakfast.ingredients && Array.isArray(day.breakfast.ingredients) && day.breakfast.ingredients.length > 0) ? `
                    <div style="flex: 1; border-left: 1px solid #ddd; padding-left: 15px; min-width: 200px;">
                      <p style="margin: 0 0 8px 0; font-size: 12px; font-weight: bold; color: #555;">Ингредиенты:</p>
                      <ul style="margin: 0; padding-left: 20px; font-size: 11px; color: #666; list-style-type: disc;">
                        ${day.breakfast.ingredients.map((ing: any) => {
                          const qty = typeof ing.quantity === 'number' ? (ing.quantity % 1 === 0 ? ing.quantity.toString() : ing.quantity.toFixed(2)) : ing.quantity
                          return `<li style="margin-bottom: 3px;">${ing.name} - ${qty} ${ing.unit}</li>`
                        }).join('')}
                      </ul>
                    </div>
                  ` : '<div style="flex: 1;"></div>'}
                </div>
              </div>
              
              <div style="margin-bottom: 8px;">
                <div style="display: flex; justify-content: space-between; gap: 20px; align-items: flex-start;">
                  <div style="flex: 1;">
                    <h3 style="font-size: 14px; margin-bottom: 5px; color: #555;"><strong>Обед:</strong></h3>
                    <p style="margin: 3px 0; margin-left: 10px; font-weight: 500;">${day.lunch.name}</p>
                    <p style="margin: 3px 0; margin-left: 10px; font-size: 12px; color: #666;">${day.lunch.calories} ккал • ${day.lunch.cooking_time} мин</p>
                  </div>
                  ${(day.lunch.ingredients && Array.isArray(day.lunch.ingredients) && day.lunch.ingredients.length > 0) ? `
                    <div style="flex: 1; border-left: 1px solid #ddd; padding-left: 15px; min-width: 200px;">
                      <p style="margin: 0 0 8px 0; font-size: 12px; font-weight: bold; color: #555;">Ингредиенты:</p>
                      <ul style="margin: 0; padding-left: 20px; font-size: 11px; color: #666; list-style-type: disc;">
                        ${day.lunch.ingredients.map((ing: any) => {
                          const qty = typeof ing.quantity === 'number' ? (ing.quantity % 1 === 0 ? ing.quantity.toString() : ing.quantity.toFixed(2)) : ing.quantity
                          return `<li style="margin-bottom: 3px;">${ing.name} - ${qty} ${ing.unit}</li>`
                        }).join('')}
                      </ul>
                    </div>
                  ` : '<div style="flex: 1;"></div>'}
                </div>
              </div>
              
              <div style="margin-bottom: 8px;">
                <div style="display: flex; justify-content: space-between; gap: 20px; align-items: flex-start;">
                  <div style="flex: 1;">
                    <h3 style="font-size: 14px; margin-bottom: 5px; color: #555;"><strong>Ужин:</strong></h3>
                    <p style="margin: 3px 0; margin-left: 10px; font-weight: 500;">${day.dinner.name}</p>
                    <p style="margin: 3px 0; margin-left: 10px; font-size: 12px; color: #666;">${day.dinner.calories} ккал • ${day.dinner.cooking_time} мин</p>
                  </div>
                  ${(day.dinner.ingredients && Array.isArray(day.dinner.ingredients) && day.dinner.ingredients.length > 0) ? `
                    <div style="flex: 1; border-left: 1px solid #ddd; padding-left: 15px; min-width: 200px;">
                      <p style="margin: 0 0 8px 0; font-size: 12px; font-weight: bold; color: #555;">Ингредиенты:</p>
                      <ul style="margin: 0; padding-left: 20px; font-size: 11px; color: #666; list-style-type: disc;">
                        ${day.dinner.ingredients.map((ing: any) => {
                          const qty = typeof ing.quantity === 'number' ? (ing.quantity % 1 === 0 ? ing.quantity.toString() : ing.quantity.toFixed(2)) : ing.quantity
                          return `<li style="margin-bottom: 3px;">${ing.name} - ${qty} ${ing.unit}</li>`
                        }).join('')}
                      </ul>
                    </div>
                  ` : '<div style="flex: 1;"></div>'}
                </div>
              </div>
            </div>
          </div>
        `).join('')}
      </div>
      
      <div style="margin-top: 30px; text-align: center; font-size: 10px; color: #999; border-top: 1px solid #ddd; padding-top: 10px;">
        MyPlateService
      </div>
    </div>
  `

  // Создаем временный элемент
  const element = document.createElement('div')
  element.innerHTML = htmlContent
  document.body.appendChild(element)

  // Настройки для PDF
  const opt = {
    margin: [10, 10, 10, 10] as [number, number, number, number],
    filename: `weekly_menu_${new Date().toISOString().split('T')[0]}.pdf`,
    image: { type: 'jpeg' as const, quality: 0.98 },
    html2canvas: { scale: 2, useCORS: true },
    jsPDF: { unit: 'mm' as const, format: 'a4' as const, orientation: 'portrait' as const }
  }

  // Генерируем PDF
  html2pdf().set(opt).from(element).save().then(() => {
    // Удаляем временный элемент
    document.body.removeChild(element)
  })
}
