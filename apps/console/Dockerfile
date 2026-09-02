# syntax=docker/dockerfile:1.7
FROM node:20-alpine AS build
RUN corepack enable && corepack prepare pnpm@9.15.9 --activate
WORKDIR /src
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/console/package.json apps/console/package.json
RUN pnpm install --frozen-lockfile
COPY apps/console apps/console
RUN pnpm --filter @traceforge/console build

FROM nginx:1.29-alpine
COPY apps/console/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/apps/console/dist /usr/share/nginx/html
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
