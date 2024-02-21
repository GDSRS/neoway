## Como executar o projeto

```sh
docker-compose up --build
```

Para ver o dados salvos, conectar no banco de dados do serviço `db` pela porta 5436 usando as credencias definidas em config-sample.yml.


## Organização do projeto

* **config/**: Com um arquivo para carregar asconfigurações de shared/config-sample.go
* **database/**: Pacote responsável por gerencias coneção com banco de dados e executar arquivos sql
* **sql/**: Arquivos sql a serem executados durante o setup do banco e para operações de limpeza e higienização dos dados
* **utils/**: Pacote para conter funções utilitárias do projeto (validação de cpf, transformações de string, etc)

## Descrição de funcionamento
Em um loop o programa:
 * Lê um número fixo de bytes do arquivo 
 * Transforma os bytes em linhas
 * Formata os dados contidos na linha (e faz a validação de cpf e cnpj)
 * Constroi uma instrução `INSERT ...` com as linhas lidas
 * Executa a instrução de inserção
 * Executa os scripts de limpeza em SQL (definidos pelo formato do nome do arquivo) da pasta **sql/**