/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Force NextJS to read environment variables at runtime rather than build time
  experimental: {
    forceSwcTransforms: true,
  }
}

module.exports = nextConfig
