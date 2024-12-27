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

COPY launchFrontend.sh /app/dist/
RUN chmod +x /app/dist/launchFrontend.sh

EXPOSE 8081

ENTRYPOINT ["/app/dist/launchFrontend.sh"]

CMD ["npx", "vite", "preview", "--port", "8081", "--host"]