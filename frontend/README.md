# google-backup

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Customize configuration

See [Vite Configuration Reference](https://vitejs.dev/config/).

## Project Setup

```sh
docker build -t pnpnm:latest .
```

```sh
docker run -it --rm -v $(pwd):/app -w /app pnpnm:latest pnpm install
```

### Compile and Hot-Reload for Development

```sh
docker run -it --rm -v $(pwd):/app -w /app -p 8080:8080 pnpnm:latest pnpm dev --port 8080 --host
```

### Compile and Minify for Production

```sh
docker run -it --rm -v $(pwd):/app -w /app pnpnm:latest pnpm build
```
