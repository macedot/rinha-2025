# 💸 Sistema de Processamento de Pagamentos 💸

Essa é a minha participação na lendária [**Rinha de Backend 2025**](https://github.com/zanfranceschi/rinha-de-backend-2025)! 🥊

## 🛠️ Tecnologias Utilizadas 👨‍💻

* **Go** - [Build simple, secure, scalable systems with Go](https://go.dev/)
* **Fiber** - [An Express-inspired web framework written in Go.](https://gofiber.io/)
* **Redis** - [The Real-time Data Platform](https://redis.io/)
* **go-redis** - [Redis Go client](https://github.com/redis/go-redis)
* **HaProxy** - [The Reliable, High Perf. TCP/HTTP Load Balancer](https://www.haproxy.org/)

## 🚀 Como Rodar

**Suba tudo com Docker!** 🐳

```bash
git clone https://github.com/macedot/rinha-2025
cd rinha-2025
docker compose up -d --build
```

## ✨ Repositório no GitHub

Curtiu? Dê uma olhada no [código fonte](https://github.com/macedot/rinha-2025) e deixe uma ⭐!

## Execução local

**AMD Ryzen 9 5950X**

```json
{
  "total_liquido": 383385.838,
  "total_bruto": 333524,
  "total_taxas": 16676.2,
  "p99": {
    "valor": "0.53ms",
    "bonus": "21%",
    "max_requests": "550",
  }
}
```

**Intel(R) Core(TM) i7-2635QM CPU @ 2.00GHz**

```json
{
  "total_liquido": 363505.34,
  "total_bruto": 332728,
  "total_taxas": 16636.4,
  "p99": {
    "valor": "3.26ms",
    "bonus": "15%",
    "max_requests": "550",
  }
}
```
## ✨ Thanks

Aprendi muito com os projetos do [Alan Silva](https://github.com/alan-venv/rinha-de-backend-2025), [Anderson Gomes](https://github.com/andersongomes001/rinha-2025/) e [Josiney Jr.](https://github.com/JosineyJr/rdb25_02). Deixe uma ⭐ pra eles!
