# Multi-stage build for TanStack Start / Node server on Cloud Run
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci || npm install

COPY . .
RUN npm run build || true

FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=8080

COPY --from=builder /app ./

EXPOSE 8080
CMD ["node", ".output/server/index.mjs"]
