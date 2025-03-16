import sys
import yaml

def generar_compose(nombre_archivo, cantidad_clientes):
    cantidad_clientes = int(cantidad_clientes)
    
    compose = {
        "name": "tp0",
        "services": {
            "server": {
                "container_name": "server",
                "image": "server:latest",
                "entrypoint": "python3 /main.py",
                "environment": [
                    "PYTHONUNBUFFERED=1",
                    "LOGGING_LEVEL=DEBUG"
                ],
                "networks": ["testing_net"]
            }
        },
        "networks": {
            "testing_net": {
                "ipam": {
                    "driver": "default",
                    "config": [{"subnet": "172.25.125.0/24"}]
                }
            }
        }
    }

    for i in range(1, cantidad_clientes + 1):
        compose["services"][f"client{i}"] = {
            "container_name": f"client{i}",
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={i}",
                "CLI_LOG_LEVEL=DEBUG"
            ],
            "networks": ["testing_net"],
            "depends_on": ["server"]
        }

    yaml_output = yaml.dump(
        compose,
        default_flow_style=False,
        sort_keys=False,
        indent=2
    )

    with open(nombre_archivo, "w") as f:
        f.write(yaml_output)

if __name__ == "__main__":
    archivo = sys.argv[1] if len(sys.argv) > 1 else "docker-compose-dev.yaml"
    clientes = int(sys.argv[2]) if len(sys.argv) > 2 else 1

    generar_compose(archivo, clientes)
