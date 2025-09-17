# OrderFlow Pro

Um ecossistema de microsserviços em Go, construído com arquitetura orientada a eventos para simular um sistema de processamento de pedidos de e-commerce. Este projeto serve como um estudo de caso prático de design de sistemas distribuídos e práticas de DevOps.

[![Go CI/CD Pipeline](https://github.com/mlucas4330/orderflow-pro/actions/workflows/ci.yml/badge.svg)](https://github.com/mlucas4330/orderflow-pro/actions/workflows/ci.yml)

---

## Sobre o Projeto

O OrderFlow Pro foi desenvolvido para explorar os desafios e as soluções na construção de sistemas de backend modernos. A premissa de negócio é um sistema de pedidos para a "Café Ponto Com", uma loja de cafés especiais, mas o foco principal é a implementação de padrões de arquitetura que garantam escalabilidade, resiliência e manutenibilidade.

## Arquitetura do Sistema

A arquitetura é composta por múltiplos microsserviços desacoplados que se comunicam através de protocolos síncronos (REST, gRPC) e assíncronos (Kafka, RabbitMQ). O objetivo é demonstrar uma separação clara de responsabilidades, onde cada serviço é um componente independente com um propósito bem definido.

---

## Tecnologias e Padrões Implementados

-   **Linguagem:** Go
-   **Comunicação Síncrona:** API RESTful (Gin), gRPC
-   **Comunicação Assíncrona:** Apache Kafka (Event Streaming), RabbitMQ (Task Queues)
-   **Persistência:** PostgreSQL (com `pgxpool`), Migrations com `goose`
-   **Cache:** Redis (padrão Cache-Aside)
-   **Containerização e Orquestração:** Docker, Docker Compose, Kubernetes (Manifestos)
-   **CI/CD:** Pipeline de Integração Contínua com GitHub Actions
-   **Padrões de Design:**
    -   Clean Architecture (com separação de `handler`, `repository`, `service`)
    -   Injeção de Dependência
    -   Idempotência para Operações de Escrita Críticas
    -   Resiliência com Lógicas de `Retry` e `Dead Letter Queues`
    -   Segurança de API com `JWT`

---

## Estrutura do Monorepo

O projeto está organizado como um monorepo, com uma estrutura que promove a partilha de código e a separação de responsabilidades:

-   **`cmd/`**: Contém os pontos de entrada (`main.go`) e `Dockerfiles` para cada microsserviço.
-   **`internal/`**: Contém a lógica de implementação privada dos serviços (handlers, repositórios, consumidores, etc.).
-   **`pkg/`**: Contém "bibliotecas" partilhadas entre os serviços, como os modelos de domínio (`model`) e os contratos de mensagem (`messaging`).
-   **`k8s/`**: Contém os manifestos de Infraestrutura como Código (`IaC`) para o deploy no Kubernetes.
-   **`db/`**: Armazena as migrações SQL versionadas.
-   **`Makefile`**: Serve como a interface de automação para o desenvolvedor, contendo os comandos para executar, testar e construir o projeto.

---

## Decisões Chave de Arquitetura

O design deste sistema foi guiado por várias decisões arquiteturais importantes:

-   **Kafka vs. RabbitMQ:** A escolha de utilizar dois brokers de mensageria foi deliberada. O Kafka foi usado para um log de eventos de alta vazão, ideal para desacoplar serviços que reagem a fatos ocorridos (ex: `OrderCreated`). O RabbitMQ foi usado para filas de tarefas específicas, garantindo a entrega de "ordens de serviço" (ex: "enviar e-mail de confirmação").

-   **Idempotência:** A operação de criação de pedidos (`POST /orders`), por ser crítica, foi tornada idempotente através do padrão `Idempotency-Key` no cabeçalho HTTP. Isto garante que falhas de rede e retentativas do cliente não resultem em pedidos duplicados.

-   **Configuração e Segredos:** A configuração segue os princípios da metodologia 12-Factor App. A aplicação é agnóstica ao ambiente e lê a sua configuração de variáveis de ambiente. A gestão de segredos é feita de forma explícita, com os valores sensíveis a serem injetados em tempo de execução e nunca versionados no Git.

-   **Estratégia de Testes:** A qualidade é garantida por duas camadas principais de testes: testes de integração que validam a comunicação com a infraestrutura real (Postgres, Redis) e testes de unidade rápidos que usam `mocks` para validar a lógica de negócio dos handlers de forma isolada.

---

## Artigos da Série

A jornada de construção e as decisões de arquitetura deste projeto foram documentadas numa série de artigos no Medium:

1.  **Artigo #1:** https://medium.com/@mlucas4330/anatomia-de-uma-api-restful-de-alta-performance-em-go-crud-cache-com-redis-e-ci-dad147f038ec
2.  **Artigo #2:** https://medium.com/@mlucas4330/construindo-a-funda%C3%A7%C3%A3o-de-um-sistema-de-microsservi%C3%A7os-em-go-do-zero-ao-banco-de-dados-com-docker-6279080b6c7c
3.  **Artigo #3:** Em andamento
