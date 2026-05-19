Based on Go and api best practices refactor the code in this api:
- No hardcoded strings (use .env)
- Model the architecture to facilitate tests
- Handle all errors
- Apply solid principles
- Only the user can mess with his stuff
- Apply best DI practices
- Write unit tests and integration tests
- Add salt for password hashes

End points(english of course)

POST /usuarios
GET /usuarios?usuario=<termo>
GET /usuarios/{usuarioId}
PUT /usuarios/{usuarioId}
DELETE /usuarios/{usuarioId}
POST /usuarios/{usuarioId}/seguir
POST /usuarios/{usuarioId}/parar-de-seguir
GET /usuarios/{usuarioId}/seguidores
GET /usuarios/{usuarioId}/seguindo
POST /usuarios/{usuarioId}/atualizar-senha

Publicações (todas exigem JWT)
POST /publicacoes
GET /publicacoes — feed com publicações próprias e de quem o usuário segue
GET /publicacoes/{publicacaoId}
PUT /publicacoes/{publicacaoId}
DELETE /publicacoes/{publicacaoId}
GET /usuarios/{usuarioId}/publicacoes
POST /publicacoes/{publicacaoId}/curtir
POST /publicacoes/{publicacaoId}/descurtir
