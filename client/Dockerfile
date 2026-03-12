FROM node:22-alpine AS deps

WORKDIR /app

COPY client/package.json client/package-lock.json ./
RUN npm ci

FROM node:22-alpine AS build

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY client/package.json ./package.json
COPY client/index.html ./index.html
COPY client/vite.config.ts ./vite.config.ts
COPY client/tsconfig.json ./tsconfig.json
COPY client/tsconfig.node.json ./tsconfig.node.json
COPY client/src ./src

RUN npm run build

FROM nginx:1.27-alpine AS runtime

COPY client/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
