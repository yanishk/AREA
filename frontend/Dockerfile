FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json ./

RUN npm install

COPY . .

RUN npm run build

# -------------------------------------------- #

FROM node:18-alpine

WORKDIR /app

RUN npm install vite

COPY --from=builder /app/dist /app/dist

EXPOSE 8081

CMD ["npx", "vite", "preview", "--port", "8081", "--host"]