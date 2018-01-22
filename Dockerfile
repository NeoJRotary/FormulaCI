FROM formulaci-dev as dev

FROM formulaci-base
WORKDIR /formulaci
COPY --from=dev /formulaci/server .
COPY server/*.sh ./
COPY web/build /formulaci/web
VOLUME [ "/formulaci/data" ]

# FormulaCI default port
EXPOSE 8099

CMD dockerd-entrypoint.sh & /formulaci/start.sh