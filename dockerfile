FROM alpine:latest as build
RUN apk add tzdata

FROM scratch as final
ADD /linux /
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
CMD ["/ipg_data_import_service_linux"]