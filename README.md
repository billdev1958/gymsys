# GymSystem

GymSystem es una aplicación de gestión de gimnasios desarrollada en Golang, utilizando algunos conceptos de Arquitectura Limpia (Clean Architecture) para garantizar un código modular, mantenible y escalable.

## Descripción

GymSystem permite a los administradores de gimnasios gestionar usuarios, cuentas, suscripciones y pagos de manera eficiente. La aplicación está diseñada para ser robusta y flexible, facilitando la integración con diversos servicios y permitiendo futuras ampliaciones sin comprometer la calidad del código.

## Características

- **Gestión de Usuarios**: Registro y administración de los usuarios del gimnasio.
- **Gestión de Cuentas**: Control y seguimiento de las cuentas de los usuarios.
- **Gestión de Suscripciones**: Administración de los planes de suscripción disponibles.
- **Gestión de Pagos**: Registro y seguimiento de los pagos realizados por los usuarios.

## Arquitectura

GymSystem sigue algunos principios de la Arquitectura Limpia, lo que implica una separación clara entre las capas de dominio, aplicación, infraestructura y presentación. Esto facilita la escalabilidad y el mantenimiento del código y todo esto sin usar paquetes de terceros mas que pgx como controlador de conexiones para postgres y un paquete de uuid, permitiendo a los desarrolladores agregar nuevas funcionalidades o modificar las existentes con un impacto mínimo en el resto del sistema.

### Capas del Sistema

1. **Dominio**: Contiene las entidades y lógica de negocio esencial. Es independiente de frameworks y bibliotecas externas.
2. **Aplicación**: Gestiona los casos de uso de la aplicación. Contiene la lógica que orquesta la interacción entre la capa de dominio y la infraestructura.
3. **Infraestructura**: Provee implementaciones específicas de tecnología para la interacción con bases de datos, servicios externos, etc.
4. **Presentación**: Maneja la interfaz de usuario y la comunicación con los servicios de la capa de aplicación.

## Requisitos

- **Golang**: Versión 1.22 o superior
- **PostgreSQL**: Para almacenamiento de datos
