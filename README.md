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

```json
{
  "participante": "anonymous",
  "total_liquido": 383317.21285,
  "total_bruto": 333464.3,
  "total_taxas": 16673.215,
  "descricao": "'total_liquido' é sua pontuação final. Equivale ao seu lucro. Fórmula: total_liquido + (total_liquido * p99.bonus) - (total_liquido * multa.porcentagem)",
  "p99": {
    "valor": "0.54ms",
    "bonus": "21%",
    "max_requests": "550",
    "descricao": "Fórmula para o bônus: max((11 - p99.valor) * 0.02, 0)"
  },
  "multa": {
    "porcentagem": 0,
    "total": 0,
    "composicao": {
      "num_inconsistencias": 0,
      "descricao": "Se 'num_inconsistencias' > 0, há multa de 35%."
    }
  }
}
```

## ✨ Thanks

Aprendi muito com os projetos do [Alan Silva](https://github.com/alan-venv/rinha-de-backend-2025), [Anderson Gomes](https://github.com/andersongomes001/rinha-2025/) e [Josiney Jr.](https://github.com/JosineyJr/rdb25_02). Deixe uma ⭐ pra eles!
