#!/bin/bash
ARCHIVO=${1:-$"docker-compose-dev.yaml"}
CLIENTES=${2:-1}
echo "Nombre del archivo de salida: $ARCHIVO"
echo "Cantidad de clientes: $CLIENTES"
python3 mi-generador.py $1 $2