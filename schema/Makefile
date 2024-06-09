# Nombre de la imagen Docker
IMAGE_NAME = gym_sys

# Nombre del contenedor docker
CONTAINER_NAME = postgres-gym_sys_container

# Nombre de el contenedor para testing
CONTAINER_NAME_TESTING = postgres-gym_sys_container_testing

# Directorio donde se encuentra el Dockerfile
SCHEMA_DIR = .

# Construir la imagen Docker
build: 
	docker build -t $(IMAGE_NAME) $(SCHEMA_DIR)

# Construir la imagen de docker
run:
	docker run --name $(CONTAINER_NAME) -p 5432:5432 -d $(IMAGE_NAME)

# Iniciar contenedor
start: 
	docker start $(CONTAINER_NAME)

# Detener el contenedor docker 
stop:
	docker stop $(CONTAINER_NAME)

trun:
	docker run --name $(CONTAINER_NAME_TESTING) -p 5432:5432 -d $(IMAGE_NAME)

tstart:
	docker start $(CONTAINER_NAME_TESTING)

tstop:
	docker stop $(CONTAINER_NAME_TESTING)

# Reconstruir la imagen y reiniciar el contenedor
rebuild:
	stop build run

# Entrar a la base de datos del contenedor de la base de datos
gymsys:
	docker exec -it $(CONTAINER_NAME) psql -U ax01 -d gym_sys

tgymsys:
	docker exec -it $(CONTAINER_NAME_TESTING) psql -U ax01 -d gym_sys

.PHONY: build run start stop rebuild gymsys
