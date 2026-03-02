import fs from 'fs'

for (const zstName of fs.globSync('dist/**/*.zst')) {
    const zstStat = fs.statSync(zstName)
    const gzName = zstName.replace(/\.zst$/i, '.gz')
    if (!fs.existsSync(gzName)) continue;
    const gzStat = fs.statSync(gzName)
    if (zstStat.size >= gzStat.size) {
        console.log('rm', zstName.replaceAll('\\', '/'))
        fs.rmSync(zstName)
    }
}