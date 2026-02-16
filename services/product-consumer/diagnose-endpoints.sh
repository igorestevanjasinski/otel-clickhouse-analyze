#!/bin/bash
# Script de diagnóstico para verificar endpoints

echo "==================================="
echo "Diagnóstico de Endpoints - Product Consumer"
echo "==================================="
echo ""

# Verificar se algo está rodando na porta 8081
echo "1. Verificando se a porta 8081 está em uso:"
if lsof -i :8081 2>/dev/null | grep -q LISTEN; then
    echo "   ✓ Porta 8081 está em uso"
    lsof -i :8081 | grep LISTEN
else
    echo "   ✗ Porta 8081 NÃO está em uso"
    echo "   → A aplicação pode não ter iniciado o servidor HTTP"
fi

echo ""
echo "2. Testando endpoints:"

for endpoint in health metrics ready; do
    echo ""
    echo "   Testing /$endpoint:"
    HTTP_CODE=$(curl -s -o /tmp/response.txt -w "%{http_code}" http://localhost:8081/$endpoint 2>/dev/null)
    
    if [ "$HTTP_CODE" = "000" ]; then
        echo "      ✗ Status: Conexão recusada (servidor não está rodando)"
    elif [ "$HTTP_CODE" = "200" ]; then
        echo "      ✓ Status: $HTTP_CODE OK"
        if [ "$endpoint" = "health" ] || [ "$endpoint" = "ready" ]; then
            echo "      Response: $(cat /tmp/response.txt)"
        else
            echo "      Response: $(cat /tmp/response.txt | head -2)"
        fi
    elif [ "$HTTP_CODE" = "404" ]; then
        echo "      ✗ Status: $HTTP_CODE Not Found"
        echo "      → Endpoint não registrado no ServeMux"
    else
        echo "      ⚠ Status: $HTTP_CODE"
        echo "      Response: $(cat /tmp/response.txt)"
    fi
done

echo ""
echo "3. Verificando processos product-consumer:"
if ps aux | grep -v grep | grep -q "product-consumer"; then
    ps aux | grep -v grep | grep "product-consumer" | head -3
else
    echo "   ✗ Nenhum processo product-consumer encontrado"
fi

echo ""
echo "==================================="
rm -f /tmp/response.txt
