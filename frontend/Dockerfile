# syntax=docker/dockerfile:1

FROM node:24-alpine3.24

WORKDIR /app

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev"]
