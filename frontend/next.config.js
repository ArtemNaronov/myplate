/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // output: 'export', // Временно отключено для разработки
  images: {
    unoptimized: true,
  },
  // Для GitHub Pages - раскомментируйте и укажите название репозитория
  // basePath: process.env.NODE_ENV === 'production' ? '/MyPlateService' : '',
  // assetPrefix: process.env.NODE_ENV === 'production' ? '/MyPlateService' : '',
}

module.exports = nextConfig


