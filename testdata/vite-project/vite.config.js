import { defineConfig } from "vite";
import { compression, defineAlgorithm } from 'vite-plugin-compression2'
import zlib from 'zlib'

export default defineConfig({
    build: {
        reportCompressedSize: false,
        assetsInlineLimit: 0,
    },
    plugins: [
        compression({
            algorithms: [
                'gzip',
                'brotliCompress',
                defineAlgorithm('zstd', {
                    params: {
                        [zlib.constants.ZSTD_c_compressionLevel]: 22
                    }
                })
            ]
        })
    ]
})