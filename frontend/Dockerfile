FROM node:19

WORKDIR /app

COPY . .

RUN npm ci

# RUN npm run build

EXPOSE 5173

CMD [ "npm", "run", "dev", "--host" ]