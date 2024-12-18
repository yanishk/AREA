FROM node:18-alpine AS builder

WORKDIR /local

COPY package.json ./

RUN npm install

COPY . .

RUN npm run build

# -------------------------------------------- #

FROM node:18-alpine

WORKDIR /local

RUN npm install vite

COPY --from=builder /local/dist /local/dist

EXPOSE 8081

CMD ["npx", "vite", "preview", "--port", "8081", "--host"]