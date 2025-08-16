import timport/** @type {import('next').NextConfig} */
const config = {
  // We'll handle CORS directly with the backend
};

module.exports = config;nfig } from "next";

const config: NextConfig = {
  // We'll handle CORS directly with the backend
};

export default config;ig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [];
  },
};

export default nextConfig;NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*', // Adding /api prefix
      },
    ];
  },
};

export default nextConfig;
