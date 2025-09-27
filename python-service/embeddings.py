from sentence_transformers import SentenceTransformer
from queries import connect_db
model = SentenceTransformer("sentence-transformers/all-MiniLM-L6-v2")

def generateEmbeddings(text):
    embeddings = model.encode([text])
    return embeddings[0].astype(float).tolist()
