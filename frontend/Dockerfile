# Frontend Dockerfile
FROM node:alpine AS npm-deps
WORKDIR /var/www/app
COPY package.json .
RUN npm install

FROM npm-deps AS dev-server
COPY . .
CMD ["npm", "run", "start"]

FROM dev-server AS static-build
RUN npm run build
CMD ["webpack-dev-server", "--content-base", "./dist/"]

FROM nginx:alpine
RUN mkdir -p /var/www/app
COPY --from=static-build /var/www/app/dist /var/www/app/
