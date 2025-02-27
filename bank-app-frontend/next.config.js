/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // This allows passing runtime environment variables
  publicRuntimeConfig: {
    apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  },
}

module.exports = nextConfig
