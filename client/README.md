# NeoHome Client

Frontend MVP for NeoHome built with:
- React + TypeScript + Vite
- React Router
- Axios
- TanStack Query
- React Hook Form + Zod
- SCSS Modules

## Local run

```bash
npm install
npm run dev
```

Default API base URL: `/api/v1` (proxied to `http://localhost:3434` in dev).

## Docker run (full stack)

From repository root:

```bash
docker compose up -d --build
```

Services:
- API: `http://localhost:3434`
- Web: `http://localhost:5173`
