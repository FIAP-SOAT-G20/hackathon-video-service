# Adicionar Endpoint de Métricas do Prometheus

## Descrição

Este PR adiciona suporte a métricas do Prometheus tanto para o servidor da API quanto para o worker consumer, habilitando monitoramento e observabilidade para a infraestrutura do serviço de vídeo.

## Mudanças

- **Adicionada dependência do Prometheus**: `github.com/prometheus/client_golang v1.20.5` no `go.mod`
- **Métricas do servidor**: Adicionado endpoint `/metrics` na porta 8081 para o deployment do servidor da API
- **Métricas do worker**: Adicionado endpoint `/metrics` na porta 8081 para o deployment do worker consumer
- **Integração no router**: Exposto endpoint `/metrics` através do router principal para acesso adicional

## Informações Adicionais

### Endpoints Disponíveis:
- **Servidor**: `http://localhost:8081/metrics` (servidor de métricas dedicado)
- **Worker**: `http://localhost:8081/metrics` (servidor de métricas dedicado)

### Métricas Expostas:
- Métricas de runtime do Go (uso de memória, goroutines, estatísticas do GC)
- Métricas de requisições HTTP (quando usando middleware do Prometheus)
- Métricas customizadas da aplicação (pronto para implementação futura)

### Testes:
```bash
# Testar métricas do servidor
curl http://localhost:8081/metrics

# Testar métricas do worker
curl http://localhost:8081/metrics
```

Ambos os serviços executam seus servidores de métricas em goroutines separadas para evitar bloquear a lógica principal da aplicação.

## Checklist

- [x] Testes passaram
- [x] Mudanças são cobertas por testes
- [x] Documentação atualizada
- [x] Mensagens de commit seguem [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)