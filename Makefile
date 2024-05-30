#!/usr/bin/make
include .env

MIGRATION_DIR := ./database/migrations
MIGRATION_TARGET := mysql
TZ := Asia/Taipei

.PHONY: migrate-up
migrate-up:
	migrate -path ${MIGRATION_DIR} -database "${MIGRATION_TARGET}://root:${DATABASE_ROOT_PASSWORD}@tcp(${DATABASE_URL})/${DATABASE_SCHEMA}" up

.PHONY: migrate-down
migrate-down:
	migrate -path ${MIGRATION_DIR} -database "${MIGRATION_TARGET}://root:${DATABASE_ROOT_PASSWORD}@tcp(${DATABASE_URL})/${DATABASE_SCHEMA}" down

.PHONY: create-migration
create-migration:
	migrate create -tz ${TZ} -ext sql -dir ${MIGRATION_DIR} $(MIGRATION_NAME)
