# Runtime
FROM nginx:1.27-alpine AS runtime

## Copy Nginx config
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf
COPY ./nginx/confd-default.conf /etc/nginx/conf.d/default.conf

## Copy Godot export html data
COPY ./web_export/ /www/data/
RUN chmod a+r -R /www/data/

RUN mv /www/data/ArcaneReveal.html /www/data/index.html
