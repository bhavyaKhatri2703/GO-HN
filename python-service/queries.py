import psycopg2
import numpy as np

def connect_db():
    conn = psycopg2.connect(
        host="localhost",
        port=5555,
        database="hackernews",
        user="postgres",
        password="postgres"
    )
    return conn


def createTables(conn):
    cursor = conn.cursor()

    cursor.execute("""
        CREATE TABLE IF NOT EXISTS newStories (
            id BIGINT PRIMARY KEY,
            by TEXT,
            type TEXT,
            text TEXT,
            url TEXT,
            title TEXT,
            full_text TEXT,
            score INT,
            sem_embedding vector(384),
            bm25_embedding bm25vector
        )
    """)

    cursor.execute("""
        CREATE TABLE IF NOT EXISTS topStories (
            id BIGINT PRIMARY KEY,
            by TEXT,
            type TEXT,
            text TEXT,
            url TEXT,
            title TEXT,
            full_text TEXT,
            score INT,
            sem_embedding vector(384),
            bm25_embedding bm25vector
        )
    """)

    conn.commit()
    cursor.close()
    print("Tables created/verified")


def saveStory(conn, story, table):
    cursor = conn.cursor()

    embedding = story['embeddings']
    if isinstance(embedding, np.ndarray):
        embedding = embedding.tolist()

    cursor.execute(f"""
        INSERT INTO {table} (id, by, type, text, url, title, full_text, score, sem_embedding)
        VALUES (%s,%s,%s,%s,%s,%s,%s,%s,%s)
        ON CONFLICT (id) DO NOTHING
    """, (
        story['id'],
        story['by'],
        story['type'],
        story['text'],
        story['url'],
        story['title'],
        story['full_text'],
        story['score'],
        embedding
    ))

    print(f"saved in {table}: {story['id']}")
    conn.commit()
    cursor.close()


def deleteStory(conn, story, table):
    cursor = conn.cursor()
    cursor.execute(f"DELETE FROM {table} WHERE id = %s", (story['id'],))
    print(f"deleted from {table}: {story['id']}")
    conn.commit()
    cursor.close()
