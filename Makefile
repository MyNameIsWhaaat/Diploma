# Makefile для перезапуска контейнера с миграциями

# Тег для контейнера миграции
MIGRATE_CONTAINER = migrate

# Команда для перезапуска контейнера миграций
db-migrate:
	@docker compose up --force-recreate --no-deps $(MIGRATE_CONTAINER)
