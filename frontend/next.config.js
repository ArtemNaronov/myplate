/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'export',
  images: {
    unoptimized: true,
  },
  // Для GitHub Pages - раскомментируйте и укажите название репозитория
  // basePath: process.env.NODE_ENV === 'production' ? '/MyPlate' : '',
  // assetPrefix: process.env.NODE_ENV === 'production' ? '/MyPlate' : '',
}

module.exports = nextConfig


