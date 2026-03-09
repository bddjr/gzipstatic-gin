import { defineConfig } from "vite";
import { compression } from 'vite-plugin-compression2'

export default defineConfig({
    build: {
        reportCompressedSize: false,
        assetsInlineLimit: 0,
    },
    plugins: [
        compression({
            algorithms: [
                'gz',
                'br',
                // 'zstd'
            ],
        }),
    ]
})