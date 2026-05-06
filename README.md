# Rinha de Backend 2026: Detecção de fraudes bancárias

Este projeto é uma implementação de um motor de detecção de fraude ultrarrápido para a Rinha de Backend 2026. A solução utiliza técnicas avançadas
de busca vetorial e engenharia de software de baixo nível para atingir latências sub-milissegundo com alta precisão.

# Arquitetura Técnica

O sistema foi desenhado seguindo princípios de Zero-Allocation e Zero-Copy, garantindo que a CPU seja utilizada quase exclusivamente para o cálculo
matemático da distância euclidiana.

### Componentes:

- Load Balancer: HAProxy em modo Round-Robin.
- API: 2 instâncias de Go 1.26+ utilizando o framework Fiber v3.
- Engine de Busca: IVF (Inverted File Index) customizado com K-Nearest Neighbors (KNN).

### Otimizações de "Estado da Arte"

Para superar os limites de hardware e a restrição de 100MB do GitHub, aplicamos as seguintes técnicas:

#### 1. Quantização Linear de 16-bits

Foi reduzida a precisão dos vetores de float32 (4 bytes) para uint16 (2 bytes) durante a extração.

- Impacto: O índice binário (ivf.bin) foi reduzido de 180MB para 90MB, permitindo o armazenamento no Git e melhorando o uso do cache da CPU.
- Estrutura do Registro (30 bytes): [14]uint16 (vetor) + uint8 (label) + uint8 (padding).

#### 2. Zero-Copy com syscall.Mmap

Em vez de carregar o arquivo na RAM ou usar io.ReadAt, foi mapeado o arquivo diretamente no espaço de endereçamento virtual do processo.

- O acesso aos dados é feito via unsafe.Slice, tratando os bytes do arquivo como structs nativas do Go sem nenhuma cópia intermediária.

#### 3. Busca Vetorial IVF (Inverted File Index)

O dataset de 3 milhões de registros é particionado em 1024 buckets.

- A API realiza uma busca linear nos centroides para identificar os 4 buckets mais prováveis ($N=4$).
- A busca final é restrita a apenas ~0.2% da base total, garantindo a latência de 3ms.

#### 4. Gestão de Recursos Docker

- GOMAXPROCS=1: Evita o process-thrashing em ambientes com CPU limitada.
- Prefork: false: Otimiza a estabilidade das Goroutines sob carga extrema.

## Como Executar

### Pré-requisitos

- Go 1.26+
- Docker & Docker Compose

1. Extração e Geração do Índice

O script de extração processa as referências, define os centroides e gera o arquivo binário quantizado.

```bash
go run scripts/extract.go
```

2. Subindo o Ambiente (Docker)

```bash
docker compose up -d --build
```

3. Teste de Carga (k6)

```bash
k6 run test/test.js
```

## Estrutura de Pastas

```plaintext
├── cmd/server # Entrypoint da API
├── internal/vector # Engine de busca IVF e Mmap
├── internal/api # Handlers e Rotas (Fiber)
├── internal/models # Definições de structs e tipos
├── resources/ # ivf.bin e metadados JSON
├── scripts/ # Script de extração e quantização
└── test/ # Scripts de teste k6 e massa de dados
```

## API

A API possui as seguintes rotas:

- `GET /ready`: Verifica se o servidor está pronto para receber requisições, envia `204 No Content`.
- `POST /fraud-score`: Realiza o cálculo de pontuação de fraude para um cliente, envia `200 OK` com o resultado, sendo o seguinte JSON:

```JSONC
{ "approved": true, "score": 0 } // Aprovado
{ "approved": false, "score": 1 } // Reprovado
```

## Licença

Este projeto está licenciado sob a licença MIT. [Ver licença](LICENSE).

---

Autor: Shobon03

Stack: Go, Fiber, HAProxy, Docker, Vector Search.
