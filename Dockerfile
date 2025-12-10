FROM ghcr.io/gohugoio/hugo:v0.152.2 AS builder

WORKDIR /src
COPY . .

USER root

ARG BASE_URL="http://localhost:8080/"
ENV HUGO_PARAMS_LOGO_IMAGE="/images/logokin.png"
ENV HUGO_PARAMS_LOGO_LINK="/"

RUN hugo --minify -b "${BASE_URL}"

FROM nginx:alpine3.22-slim

COPY --from=builder /src/public /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
