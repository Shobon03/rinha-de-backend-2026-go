# Rinha de Backend 2026: Detecção de fraudes bancárias

Este projeto é uma implementação de um motor de detecção de fraude ultrarrápido para a Rinha de Backend 2026. A solução utiliza técnicas avançadas
de busca vetorial e engenharia de software de baixo nível para atingir latências sub-milissegundo com alta precisão.

# Arquitetura Técnica

O sistema foi desenhado seguindo princípios de **Zero-Allocation** e **Zero-Copy**, garantindo que a CPU seja utilizada quase exclusivamente para o cálculo matemático da distância euclidiana.

### Componentes:

- **Load Balancer**: HAProxy configurado com 2 instâncias de API.
- **API**: 2 instâncias de Go 1.26+ utilizando a biblioteca padrão (`net/http`) para overhead mínimo.
- **Engine de Busca**: IVF (Inverted File Index) customizado com K-Nearest Neighbors (KNN).

## API

- `GET /ready`: Retorna `204 No Content` quando o índice está carregado.
- `POST /fraud-score`: Recebe o payload da transação e retorna o score de fraude.

## Estrutura de Pastas

```plaintext
├── cmd/server      # Entrypoint da API (UDS Listener)
├── internal/vector  # Engine de busca IVF, Mmap e Normalização
├── internal/api     # Handlers e Rotas (Standard Library)
├── internal/models  # Definições de structs e tipos
├── resources/       # ivf.bin e metadados JSON
├── scripts/         # Script de extração e quantização (K-Means)
└── test/            # Scripts de teste k6 e massa de dados
```

## Como Executar

### Pré-requisitos

- Go 1.26+
- Docker & Docker Compose

1. **Extração e Geração do Índice**

```bash
go run scripts/extract.go
```

2. **Subindo o Ambiente (Docker)**

```bash
docker compose up -d --build
```

3. **Teste de Carga (k6)**

```bash
k6 run test/test.js
```

### Otimizações

#### 1. Zero-Allocation JSON Parsing

Utiliza um `sync.Pool` de buffers para ler e decodificar os payloads JSON. Isso evita alocações frequentes no heap, reduzindo drasticamente a pressão sobre o Garbage Collector (GC) e evitando picos de latência (p99).

#### 2. Matemática Inteira (Integer Squared Distance)

A busca vetorial utiliza aritmética de inteiros (`int64`) para calcular a distância quadrada entre vetores quantizados. Isso elimina a necessidade de operações de ponto flutuante caras e permite um _early exit_ agressivo durante a busca.

#### 3. Zero-Copy com syscall.Mmap

O índice binário de 3 milhões de registros é mapeado diretamente no espaço de endereçamento virtual do processo. O acesso aos dados é feito via `unsafe.Slice`, tratando os bytes do arquivo como structs nativas sem nenhuma cópia intermediária.

#### 4. Quantização Linear de 16-bits

Reduzida a precisão dos vetores de `float32` (4 bytes) para `uint16` (2 bytes). O índice final (`ivf.bin`) ocupa ~90MB, permitindo que o conjunto de dados caia no cache da CPU com mais eficiência.

#### 5. Busca Vetorial IVF (Inverted File Index)

O dataset é particionado em 1024 buckets. A API identifica os 4 buckets mais prováveis ($N=4$) e realiza a busca exata apenas neles (~0.4% da base total).

## Licença

Este projeto está licenciado sob a licença MIT.

---

Autor: Shobon03  
Stack: Go (std), HAProxy, UDS, IVF, Mmap, Integer Math.
