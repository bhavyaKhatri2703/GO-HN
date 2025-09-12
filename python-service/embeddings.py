from sentence_transformers import SentenceTransformer
from queries import connect_db
model = SentenceTransformer("sentence-transformers/all-MiniLM-L6-v2")

def generateEmbeddings(text):
    embeddings = model.encode([text])
    return embeddings[0].astype(float).tolist()


def search(query,conn):
    embedding = generateEmbeddings(query)
    print("hello??")
    cursor = conn.cursor()

    cursor.execute("""
         SELECT  title,text, (sem_embedding <=> %s::vector) AS distance
         FROM topstories
         ORDER BY distance ASC
         LIMIT 20;
     """, (embedding,))

    results = cursor.fetchall()

    print(results)

print(generateEmbeddings("GoLang Javascript Cyber Security"))
