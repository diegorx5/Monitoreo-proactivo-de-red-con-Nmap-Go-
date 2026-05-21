

# 🔎 NAC Network Monitor con Nmap + Go

Herramienta simple de monitoreo proactivo de red que permite detectar cambios en dispositivos conectados a la red local utilizando Nmap y desarrollada en Go (Golang).

---

## 🚀 ¿Qué hace?

- Descubre dispositivos activos en la red (ping sweep con Nmap)

- Escanea puertos abiertos (full TCP scan)

- Detecta el sistema operativo de los dispositivos (OS fingerprinting)

- Captura las conexiones activas locales (simulado con netstat -an)

- Registra cada escaneo con fecha/hora en una base de datos SQLite

---

## 🎯 Objetivo

Este proyecto busca simular un sistema NAC básico capaz de:

- Detectar nuevas conexiones o puertos abiertos en la red

- Registrar actividad sospechosa como reverse shells o conexiones no autorizadas

- Proporcionar histórico de escaneos para análisis forense o SOC

---

## 🧪 Prueba realizada

Durante la prueba, se simuló:

- Estado inicial sin puertos sospechosos

- Ejecución de un ataque reverse shell (evadiendo antivirus con técnicas XOR, AES, PowerShell)

- El sistema detectó automáticamente la nueva conexión y la registró

---

## 📦 Requisitos

- Tener instalado Nmap  

- Tener instalado Go (versión recomendada: 1.20+)

---

## 🔧 Instalación y ejecución

carpetas en GO

MOnitoreo-proactivo-de-red-con-nmap+go/
│
├── go.mod 
├── inventario.db
├── scan.go
├── serverWEB.go
