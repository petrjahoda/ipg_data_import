#!/usr/bin/env bash
cd linux
upx ipg_data_import_service_linux
cd ..
docker rmi -f petrjahoda/ipg_data_import_service:latest
docker build -t petrjahoda/ipg_data_import_service:latest .
docker push petrjahoda/ipg_data_import_service:latest

docker rmi -f petrjahoda/ipg_data_import_service:2020.3.3
docker build -t petrjahoda/ipg_data_import_service:2020.3.3 .
docker push petrjahoda/ipg_data_import_service:2020.3.3
