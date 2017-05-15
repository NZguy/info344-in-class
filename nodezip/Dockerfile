FROM node
RUN mkdir /app
## WORKDIR is kinda like cd
WORKDIR  /app
## Explicityly only copy over this so that the npm install doesn't have to be rerun on subsequent builds
COPY package.json /app
RUN npm install
COPY . /app
## WE use EXPOSE in case someone runs with the -P flag (capital p)
EXPOSE 80
CMD ["node", "index.js"]
