import requests
import json

# Test health endpoint
response = requests.post('http://localhost:8080/health')
print(f"Health: {response.json()}")

# Test counter endpoint
for i in range(3):
    response = requests.post('http://localhost:8080/api/counter/next')
    print(f"Counter {i+1}: {response.json()}")

# Test LLM endpoint with SSE
print("\nLLM Stream Test:")
response = requests.post('http://localhost:8080/api/llm/chat', 
                        json={'prompt': 'Hello AI!'},
                        stream=True)
result = ""
for line in response.iter_lines():
    if line:
        line = line.decode('utf-8')
        if line.startswith('data: '):
            token = line[6:]
            result += token
            print(f"Token: {repr(token)}")

print(f"\nFull response: {result}")
