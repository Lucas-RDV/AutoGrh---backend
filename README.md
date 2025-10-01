# 📌 AutoGRH – Mapeamento de Rotas (Backend)

Este documento lista as rotas expostas pelo backend do **AutoGRH**, conforme definido no `router.go` e nos controllers.

Cada seção contém:

* **Método e Endpoint**
* **Descrição**
* **JSON esperado (request)**
* **JSON retornado (response)**

---

## 🔑 Autenticação

### `POST /auth/login`

* **Descrição**: Realiza login de usuário ativo.
* **Request JSON**:

```json
{
  "login": "string",
  "senha": "string"
}
```

* **Response JSON**:

```json
{
  "token": "string",
  "expiresAt": "2025-01-01T12:00:00Z",
  "usuario": {
    "id": 1,
    "username": "admin",
    "isAdmin": true
  }
}
```

---

## 👤 Usuários (Admin)

### `GET /admin/usuarios`

* Lista todos os usuários.
* **Response JSON**:

```json
[
  {
    "id": 1,
    "username": "admin",
    "isAdmin": true,
    "ativo": true
  }
]
```

### `POST /admin/usuarios`

* Cria novo usuário.
* **Request JSON**:

```json
{
  "username": "string",
  "senha": "string",
  "isAdmin": true
}
```

### `PUT /admin/usuarios/{id}`

* Atualiza usuário existente.
* **Request JSON**:

```json
{
  "username": "string",
  "senha": "novaSenhaOpcional",
  "isAdmin": false,
  "ativo": true
}
```

### `DELETE /admin/usuarios/{id}`

* Desativa usuário (soft delete).

---

## 🧑 Pessoas

### `GET /pessoas`

* Lista todas as pessoas.

### `POST /pessoas`

* Cria nova pessoa.
* **Request JSON**:

```json
{
  "nome": "string",
  "cpf": "string",
  "dataNascimento": "2025-01-01",
  "endereco": "string"
}
```

### `PUT /pessoas/{id}`

* Atualiza pessoa.

### `DELETE /pessoas/{id}`

* Remove pessoa (soft delete).

---

## 👔 Funcionários

### `GET /funcionarios`

* Lista todos os funcionários.

### `POST /funcionarios`

* Cria funcionário a partir de pessoa.
* **Request JSON**:

```json
{
  "pessoaID": 1,
  "cargo": "string",
  "setor": "string",
  "dataAdmissao": "2025-01-01"
}
```

### `PUT /funcionarios/{id}`

* Atualiza funcionário.

### `DELETE /funcionarios/{id}`

* Remove funcionário (soft delete).

---

## 📄 Documentos

### `GET /documentos`

* Lista documentos.

### `POST /documentos`

* Insere documento associado ao funcionário.
* **Request JSON (multipart/form-data)**:

```
file: (arquivo)
funcionarioID: number
descricao: string
```

### `GET /documentos/{id}`

* Download de documento.

### `DELETE /documentos/{id}`

* Remove documento.

---

## ⛔ Faltas

### `GET /faltas`

* Lista faltas.

### `POST /faltas`

* Cria nova falta.
* **Request JSON**:

```json
{
  "funcionarioID": 1,
  "data": "2025-01-01",
  "justificada": false
}
```

### `DELETE /faltas/{id}`

* Remove falta.

---

## 🏖️ Férias

### `GET /ferias`

* Lista férias.

### `POST /ferias`

* Cria férias manualmente.
* **Request JSON**:

```json
{
  "funcionarioID": 1,
  "dias": 30,
  "inicio": "2025-01-01",
  "fim": "2025-01-30"
}
```

### `PUT /ferias/{id}/vencer`

* Marca férias como vencidas.

### `PUT /ferias/{id}/terco`

* Marca 1/3 como pago.

---

## 💤 Descansos

### `GET /descansos`

* Lista descansos.

### `POST /descansos`

* Cria novo descanso.
* **Request JSON**:

```json
{
  "feriasID": 1,
  "inicio": "2025-01-01",
  "fim": "2025-01-10"
}
```

### `PUT /descansos/{id}/aprovar`

* Admin aprova descanso.

### `PUT /descansos/{id}/pagar`

* Admin marca descanso como pago.

### `GET /descansos/aprovados`

* Lista descansos aprovados.

### `GET /descansos/pendentes`

* Lista descansos pendentes.

---

## 💰 Salários

### `GET /salarios`

* Lista salários.

### `POST /salarios`

* Cria novo salário (encerra anterior).
* **Request JSON**:

```json
{
  "funcionarioID": 1,
  "valor": 3000.00,
  "inicio": "2025-01-01"
}
```

---

## 📑 Folhas de Pagamento

### `GET /folhas`

* Lista folhas de pagamento.

### `POST /folhas`

* Cria nova folha de pagamento.

### `PUT /folhas/{id}/recalcular`

* Recalcula folha de pagamento.

### `PUT /folhas/{id}/fechar`

* Admin fecha/paga folha.

---

## 💵 Pagamentos

### `GET /pagamentos`

* Lista pagamentos.

### `PUT /pagamentos/{id}`

* Atualiza manualmente um pagamento.

---

## 💳 Vales

### `GET /vales`

* Lista vales.

### `POST /vales`

* Cria novo vale.
* **Request JSON**:

```json
{
  "funcionarioID": 1,
  "valor": 500.00,
  "descricao": "Adiantamento"
}
```

### `PUT /vales/{id}`

* Atualiza vale (antes de aprovado).

### `DELETE /vales/{id}`

* Admin exclui vale.

### `POST /vales/folha`

* Cria folha de vales.

### `PUT /vales/folha/{id}/aprovar`

* Admin aprova folha de vales.

### `PUT /vales/folha/{id}/pagar`

* Admin paga folha de vales.

---

## 📢 Avisos (Sprint 5)

### `GET /avisos`

* Lista avisos ativos.
* **Response JSON**:

```json
[
  {
    "id": 1,
    "tipo": "ferias",
    "descricao": "Férias de João vencem em 2025-10-10",
    "dataEvento": "2025-10-10"
  }
]
```

---

# ✅ Observações

* Todas as rotas protegidas por `AuthMiddleware` exigem **JWT válido**.
* Perfis de acesso definidos:

    * **RH**: criar/listar/recalcular.
    * **Admin**: aprovar/pagar/excluir.
* JSONs podem variar levemente dependendo do estado do banco, mas seguem esse padrão geral.
