SHELL := /bin/bash

.PHONY: migrate-up migrate-down migrate-create

migrate-up:
	go run ./migrate up

migrate-down:
	go run ./migrate down

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=add_users_table"; exit 1; fi
	@max=0; \
	for file in migrations/*.up.sql; do \
		[ -e "$$file" ] || continue; \
		base=$$(basename "$$file"); \
		prefix=$${base%%_*}; \
		if [[ "$$prefix" =~ ^[0-9]{1,6}$$ ]]; then \
			num=$$((10#$$prefix)); \
			if [ $$num -gt $$max ]; then max=$$num; fi; \
		fi; \
	done; \
	version=$$(printf "%05d" $$((max + 1))); \
	up_file="migrations/$${version}_$(name).up.sql"; \
	down_file="migrations/$${version}_$(name).down.sql"; \
	printf -- "-- Write your UP migration SQL here\n" > "$$up_file"; \
	printf -- "-- Write your DOWN migration SQL here\n" > "$$down_file"; \
	echo "Created $$up_file and $$down_file"
