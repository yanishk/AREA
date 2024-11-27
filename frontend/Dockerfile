FROM node:14-alpine AS builder

WORKDIR /

COPY package.json ./

RUN npm install

COPY . .

RUN npm run build

# -------------------------------------------- #

FROM node:14-alpine

RUN npm install -g serve

COPY --from=builder /dist /dist

EXPOSE 8081

CMD ["serve", "-s", "/dist", "-l", "8081"]
